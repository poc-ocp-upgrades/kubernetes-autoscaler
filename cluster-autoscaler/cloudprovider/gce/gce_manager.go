package gce

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	provider_gce "k8s.io/kubernetes/pkg/cloudprovider/providers/gce"
	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gce "google.golang.org/api/compute/v1"
	gcfg "gopkg.in/gcfg.v1"
	"k8s.io/klog"
)

const (
	refreshInterval		= 1 * time.Minute
	machinesRefreshInterval	= 1 * time.Hour
	httpTimeout		= 30 * time.Second
	scaleToZeroSupported	= true
)

var (
	defaultOAuthScopes []string = []string{"https://www.googleapis.com/auth/compute", "https://www.googleapis.com/auth/devstorage.read_only", "https://www.googleapis.com/auth/service.management.readonly", "https://www.googleapis.com/auth/servicecontrol"}
)

type GceManager interface {
	Refresh() error
	Cleanup() error
	GetMigs() []*MigInformation
	GetMigNodes(mig Mig) ([]string, error)
	GetMigForInstance(instance *GceRef) (Mig, error)
	GetMigTemplateNode(mig Mig) (*apiv1.Node, error)
	GetResourceLimiter() (*cloudprovider.ResourceLimiter, error)
	GetMigSize(mig Mig) (int64, error)
	SetMigSize(mig Mig, size int64) error
	DeleteInstances(instances []*GceRef) error
}
type gceManagerImpl struct {
	cache				GceCache
	lastRefresh			time.Time
	machinesCacheLastRefresh	time.Time
	GceService			AutoscalingGceClient
	location			string
	projectId			string
	templates			*GceTemplateBuilder
	interrupt			chan struct{}
	regional			bool
	explicitlyConfigured		map[GceRef]bool
	migAutoDiscoverySpecs		[]cloudprovider.MIGAutoDiscoveryConfig
}

func CreateGceManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, regional bool) (GceManager, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	tokenSource := google.ComputeTokenSource("")
	if len(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")) > 0 {
		tokenSource, err = google.DefaultTokenSource(oauth2.NoContext, gce.ComputeScope)
		if err != nil {
			return nil, err
		}
	}
	var projectId, location string
	if configReader != nil {
		var cfg provider_gce.ConfigFile
		if err := gcfg.ReadInto(&cfg, configReader); err != nil {
			klog.Errorf("Couldn't read config: %v", err)
			return nil, err
		}
		if cfg.Global.TokenURL == "" {
			klog.Warning("Empty tokenUrl in cloud config")
		} else {
			tokenSource = provider_gce.NewAltTokenSource(cfg.Global.TokenURL, cfg.Global.TokenBody)
			klog.V(1).Infof("Using TokenSource from config %#v", tokenSource)
		}
		projectId = cfg.Global.ProjectID
		location = cfg.Global.LocalZone
	} else {
		klog.V(1).Infof("Using default TokenSource %#v", tokenSource)
	}
	if len(projectId) == 0 || len(location) == 0 {
		discoveredProjectId, discoveredLocation, err := getProjectAndLocation(regional)
		if err != nil {
			return nil, err
		}
		if len(projectId) == 0 {
			projectId = discoveredProjectId
		}
		if len(location) == 0 {
			location = discoveredLocation
		}
	}
	klog.V(1).Infof("GCE projectId=%s location=%s", projectId, location)
	client := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client.Timeout = httpTimeout
	gceService, err := NewAutoscalingGceClientV1(client, projectId)
	if err != nil {
		return nil, err
	}
	manager := &gceManagerImpl{cache: NewGceCache(gceService), GceService: gceService, location: location, regional: regional, projectId: projectId, templates: &GceTemplateBuilder{}, interrupt: make(chan struct{}), explicitlyConfigured: make(map[GceRef]bool)}
	if err := manager.fetchExplicitMigs(discoveryOpts.NodeGroupSpecs); err != nil {
		return nil, fmt.Errorf("failed to fetch MIGs: %v", err)
	}
	if manager.migAutoDiscoverySpecs, err = discoveryOpts.ParseMIGAutoDiscoverySpecs(); err != nil {
		return nil, err
	}
	if err := manager.forceRefresh(); err != nil {
		return nil, err
	}
	go wait.Until(func() {
		if err := manager.cache.RegenerateInstancesCache(); err != nil {
			klog.Errorf("Error while regenerating Mig cache: %v", err)
		}
	}, time.Hour, manager.interrupt)
	return manager, nil
}
func (m *gceManagerImpl) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	close(m.interrupt)
	return nil
}
func (m *gceManagerImpl) registerMig(mig Mig) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	changed := m.cache.RegisterMig(mig)
	if changed {
		if _, err := m.GetMigTemplateNode(mig); err != nil {
			klog.Errorf("Can't build node from template for %s, won't be able to scale from 0: %v", mig.GceRef().String(), err)
		}
	}
	return changed
}
func (m *gceManagerImpl) GetMigSize(mig Mig) (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	targetSize, err := m.GceService.FetchMigTargetSize(mig.GceRef())
	if err != nil {
		return -1, err
	}
	return targetSize, nil
}
func (m *gceManagerImpl) SetMigSize(mig Mig, size int64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.V(0).Infof("Setting mig size %s to %d", mig.Id(), size)
	return m.GceService.ResizeMig(mig.GceRef(), size)
}
func (m *gceManagerImpl) DeleteInstances(instances []*GceRef) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(instances) == 0 {
		return nil
	}
	commonMig, err := m.GetMigForInstance(instances[0])
	if err != nil {
		return err
	}
	for _, instance := range instances {
		mig, err := m.GetMigForInstance(instance)
		if err != nil {
			return err
		}
		if mig != commonMig {
			return fmt.Errorf("Cannot delete instances which don't belong to the same MIG.")
		}
	}
	return m.GceService.DeleteInstances(commonMig.GceRef(), instances)
}
func (m *gceManagerImpl) GetMigs() []*MigInformation {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.cache.GetMigs()
}
func (m *gceManagerImpl) GetMigForInstance(instance *GceRef) (Mig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.cache.GetMigForInstance(instance)
}
func (m *gceManagerImpl) GetMigNodes(mig Mig) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	instances, err := m.GceService.FetchMigInstances(mig.GceRef())
	if err != nil {
		return []string{}, err
	}
	result := make([]string, 0)
	for _, ref := range instances {
		result = append(result, fmt.Sprintf("gce://%s/%s/%s", ref.Project, ref.Zone, ref.Name))
	}
	return result, nil
}
func (m *gceManagerImpl) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.lastRefresh.Add(refreshInterval).After(time.Now()) {
		return nil
	}
	return m.forceRefresh()
}
func (m *gceManagerImpl) forceRefresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.clearMachinesCache()
	if err := m.fetchAutoMigs(); err != nil {
		klog.Errorf("Failed to fetch MIGs: %v", err)
		return err
	}
	m.lastRefresh = time.Now()
	klog.V(2).Infof("Refreshed GCE resources, next refresh after %v", m.lastRefresh.Add(refreshInterval))
	return nil
}
func (m *gceManagerImpl) fetchExplicitMigs(specs []string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	changed := false
	for _, spec := range specs {
		mig, err := m.buildMigFromFlag(spec)
		if err != nil {
			return err
		}
		if m.registerMig(mig) {
			changed = true
		}
		m.explicitlyConfigured[mig.GceRef()] = true
	}
	if changed {
		return m.cache.RegenerateInstancesCache()
	}
	return nil
}
func (m *gceManagerImpl) buildMigFromFlag(flag string) (Mig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s, err := dynamic.SpecFromString(flag, scaleToZeroSupported)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node group spec: %v", err)
	}
	return m.buildMigFromSpec(s)
}
func (m *gceManagerImpl) buildMigFromAutoCfg(link string, cfg cloudprovider.MIGAutoDiscoveryConfig) (Mig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := &dynamic.NodeGroupSpec{Name: link, MinSize: cfg.MinSize, MaxSize: cfg.MaxSize, SupportScaleToZero: scaleToZeroSupported}
	return m.buildMigFromSpec(s)
}
func (m *gceManagerImpl) buildMigFromSpec(s *dynamic.NodeGroupSpec) (Mig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid node group spec: %v", err)
	}
	project, zone, name, err := ParseMigUrl(s.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mig url: %s got error: %v", s.Name, err)
	}
	mig := &gceMig{gceRef: GceRef{Project: project, Name: name, Zone: zone}, gceManager: m, minSize: s.MinSize, maxSize: s.MaxSize}
	return mig, nil
}
func (m *gceManagerImpl) fetchAutoMigs() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	exists := make(map[GceRef]bool)
	changed := false
	for _, cfg := range m.migAutoDiscoverySpecs {
		links, err := m.findMigsNamed(cfg.Re)
		if err != nil {
			return fmt.Errorf("cannot autodiscover managed instance groups: %v", err)
		}
		for _, link := range links {
			mig, err := m.buildMigFromAutoCfg(link, cfg)
			if err != nil {
				return err
			}
			exists[mig.GceRef()] = true
			if m.explicitlyConfigured[mig.GceRef()] {
				klog.V(3).Infof("Ignoring explicitly configured MIG %s in autodiscovery.", mig.GceRef().String())
				continue
			}
			if m.registerMig(mig) {
				klog.V(3).Infof("Autodiscovered MIG %s using regexp %s", mig.GceRef().String(), cfg.Re.String())
				changed = true
			}
		}
	}
	for _, mig := range m.GetMigs() {
		if !exists[mig.Config.GceRef()] && !m.explicitlyConfigured[mig.Config.GceRef()] {
			m.cache.UnregisterMig(mig.Config)
			changed = true
		}
	}
	if changed {
		return m.cache.RegenerateInstancesCache()
	}
	return nil
}
func (m *gceManagerImpl) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.cache.GetResourceLimiter()
}
func (m *gceManagerImpl) clearMachinesCache() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.machinesCacheLastRefresh.Add(machinesRefreshInterval).After(time.Now()) {
		return
	}
	machinesCache := make(map[MachineTypeKey]*gce.MachineType)
	m.cache.SetMachinesCache(machinesCache)
	nextRefresh := time.Now()
	m.machinesCacheLastRefresh = nextRefresh
	klog.V(2).Infof("Cleared machine types cache, next clear after %v", nextRefresh)
}
func getProjectAndLocation(regional bool) (string, string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result, err := metadata.Get("instance/zone")
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(result, "/")
	if len(parts) != 4 {
		return "", "", fmt.Errorf("unexpected response: %s", result)
	}
	location := parts[3]
	if regional {
		location, err = provider_gce.GetGCERegion(location)
		if err != nil {
			return "", "", err
		}
	}
	projectID, err := metadata.ProjectID()
	if err != nil {
		return "", "", err
	}
	return projectID, location, nil
}
func (m *gceManagerImpl) findMigsNamed(name *regexp.Regexp) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.regional {
		return m.findMigsInRegion(m.location, name)
	}
	return m.GceService.FetchMigsWithName(m.location, name)
}
func (m *gceManagerImpl) getZones(region string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	zones, err := m.GceService.FetchZones(region)
	if err != nil {
		return nil, fmt.Errorf("cannot get zones for GCE region %s: %v", region, err)
	}
	return zones, nil
}
func (m *gceManagerImpl) findMigsInRegion(region string, name *regexp.Regexp) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	links := make([]string, 0)
	zones, err := m.getZones(region)
	if err != nil {
		return nil, err
	}
	for _, z := range zones {
		zl, err := m.GceService.FetchMigsWithName(z, name)
		if err != nil {
			return nil, err
		}
		for _, link := range zl {
			links = append(links, link)
		}
	}
	return links, nil
}
func (m *gceManagerImpl) GetMigTemplateNode(mig Mig) (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	template, err := m.GceService.FetchMigTemplate(mig.GceRef())
	if err != nil {
		return nil, err
	}
	cpu, mem, err := m.getCpuAndMemoryForMachineType(template.Properties.MachineType, mig.GceRef().Zone)
	if err != nil {
		return nil, err
	}
	return m.templates.BuildNodeFromTemplate(mig, template, cpu, mem)
}
func (m *gceManagerImpl) getCpuAndMemoryForMachineType(machineType string, zone string) (cpu int64, mem int64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.HasPrefix(machineType, "custom-") {
		return parseCustomMachineType(machineType)
	}
	machine := m.cache.GetMachineFromCache(machineType, zone)
	if machine == nil {
		machine, err = m.GceService.FetchMachineType(zone, machineType)
		if err != nil {
			return 0, 0, err
		}
		m.cache.AddMachineToCache(machineType, zone, machine)
	}
	return machine.GuestCpus, machine.MemoryMb * units.MiB, nil
}
func parseCustomMachineType(machineType string) (cpu, mem int64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var count int
	count, err = fmt.Sscanf(machineType, "custom-%d-%d", &cpu, &mem)
	if err != nil {
		return
	}
	if count != 2 {
		return 0, 0, fmt.Errorf("failed to parse all params in %s", machineType)
	}
	mem = mem * units.MiB
	return
}
