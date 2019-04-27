package azure

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/Azure/go-autorest/autorest/azure"
	"gopkg.in/gcfg.v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/klog"
)

const (
	vmTypeVMSS			= "vmss"
	vmTypeStandard			= "standard"
	vmTypeACS			= "acs"
	vmTypeAKS			= "aks"
	scaleToZeroSupportedStandard	= false
	scaleToZeroSupportedVMSS	= true
	refreshInterval			= 1 * time.Minute
	deploymentParametersPath	= "/var/lib/azure/azuredeploy.parameters.json"
)

type AzureManager struct {
	config			*Config
	azClient		*azClient
	env			azure.Environment
	asgCache		*asgCache
	lastRefresh		time.Time
	asgAutoDiscoverySpecs	[]cloudprovider.LabelAutoDiscoveryConfig
	explicitlyConfigured	map[string]bool
}
type Config struct {
	Cloud				string			`json:"cloud" yaml:"cloud"`
	TenantID			string			`json:"tenantId" yaml:"tenantId"`
	SubscriptionID			string			`json:"subscriptionId" yaml:"subscriptionId"`
	ResourceGroup			string			`json:"resourceGroup" yaml:"resourceGroup"`
	VMType				string			`json:"vmType" yaml:"vmType"`
	AADClientID			string			`json:"aadClientId" yaml:"aadClientId"`
	AADClientSecret			string			`json:"aadClientSecret" yaml:"aadClientSecret"`
	AADClientCertPath		string			`json:"aadClientCertPath" yaml:"aadClientCertPath"`
	AADClientCertPassword		string			`json:"aadClientCertPassword" yaml:"aadClientCertPassword"`
	UseManagedIdentityExtension	bool			`json:"useManagedIdentityExtension" yaml:"useManagedIdentityExtension"`
	Deployment			string			`json:"deployment" yaml:"deployment"`
	DeploymentParameters		map[string]interface{}	`json:"deploymentParameters" yaml:"deploymentParameters"`
	ClusterName			string			`json:"clusterName" yaml:"clusterName"`
	NodeResourceGroup		string			`json:"nodeResourceGroup" yaml:"nodeResourceGroup"`
}

func (c *Config) TrimSpace() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Cloud = strings.TrimSpace(c.Cloud)
	c.TenantID = strings.TrimSpace(c.TenantID)
	c.SubscriptionID = strings.TrimSpace(c.SubscriptionID)
	c.ResourceGroup = strings.TrimSpace(c.ResourceGroup)
	c.VMType = strings.TrimSpace(c.VMType)
	c.AADClientID = strings.TrimSpace(c.AADClientID)
	c.AADClientSecret = strings.TrimSpace(c.AADClientSecret)
	c.AADClientCertPath = strings.TrimSpace(c.AADClientCertPath)
	c.AADClientCertPassword = strings.TrimSpace(c.AADClientCertPassword)
	c.Deployment = strings.TrimSpace(c.Deployment)
	c.ClusterName = strings.TrimSpace(c.ClusterName)
	c.NodeResourceGroup = strings.TrimSpace(c.NodeResourceGroup)
}
func CreateAzureManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*AzureManager, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	var cfg Config
	if configReader != nil {
		if err := gcfg.ReadInto(&cfg, configReader); err != nil {
			klog.Errorf("Couldn't read config: %v", err)
			return nil, err
		}
	} else {
		cfg.Cloud = os.Getenv("ARM_CLOUD")
		cfg.SubscriptionID = os.Getenv("ARM_SUBSCRIPTION_ID")
		cfg.ResourceGroup = os.Getenv("ARM_RESOURCE_GROUP")
		cfg.TenantID = os.Getenv("ARM_TENANT_ID")
		cfg.AADClientID = os.Getenv("ARM_CLIENT_ID")
		cfg.AADClientSecret = os.Getenv("ARM_CLIENT_SECRET")
		cfg.VMType = strings.ToLower(os.Getenv("ARM_VM_TYPE"))
		cfg.AADClientCertPath = os.Getenv("ARM_CLIENT_CERT_PATH")
		cfg.AADClientCertPassword = os.Getenv("ARM_CLIENT_CERT_PASSWORD")
		cfg.Deployment = os.Getenv("ARM_DEPLOYMENT")
		cfg.ClusterName = os.Getenv("AZURE_CLUSTER_NAME")
		cfg.NodeResourceGroup = os.Getenv("AZURE_NODE_RESOURCE_GROUP")
		useManagedIdentityExtensionFromEnv := os.Getenv("ARM_USE_MANAGED_IDENTITY_EXTENSION")
		if len(useManagedIdentityExtensionFromEnv) > 0 {
			cfg.UseManagedIdentityExtension, err = strconv.ParseBool(useManagedIdentityExtensionFromEnv)
			if err != nil {
				return nil, err
			}
		}
	}
	cfg.TrimSpace()
	if cfg.VMType == "" {
		cfg.VMType = vmTypeVMSS
	}
	if cfg.VMType == vmTypeStandard && len(cfg.DeploymentParameters) == 0 {
		parameters, err := readDeploymentParameters(deploymentParametersPath)
		if err != nil {
			klog.Errorf("readDeploymentParameters failed with error: %v", err)
			return nil, err
		}
		cfg.DeploymentParameters = parameters
	}
	env := azure.PublicCloud
	if cfg.Cloud != "" {
		env, err = azure.EnvironmentFromName(cfg.Cloud)
		if err != nil {
			return nil, err
		}
	}
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	klog.Infof("Starting azure manager with subscription ID %q", cfg.SubscriptionID)
	azClient, err := newAzClient(&cfg, &env)
	if err != nil {
		return nil, err
	}
	manager := &AzureManager{config: &cfg, env: env, azClient: azClient, explicitlyConfigured: make(map[string]bool)}
	cache, err := newAsgCache()
	if err != nil {
		return nil, err
	}
	manager.asgCache = cache
	specs, err := discoveryOpts.ParseLabelAutoDiscoverySpecs()
	if err != nil {
		return nil, err
	}
	manager.asgAutoDiscoverySpecs = specs
	if err := manager.fetchExplicitAsgs(discoveryOpts.NodeGroupSpecs); err != nil {
		return nil, err
	}
	if err := manager.forceRefresh(); err != nil {
		return nil, err
	}
	return manager, nil
}
func (m *AzureManager) fetchExplicitAsgs(specs []string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	changed := false
	for _, spec := range specs {
		asg, err := m.buildAsgFromSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse node group spec: %v", err)
		}
		if m.RegisterAsg(asg) {
			changed = true
		}
		m.explicitlyConfigured[asg.Id()] = true
	}
	if changed {
		if err := m.regenerateCache(); err != nil {
			return err
		}
	}
	return nil
}
func (m *AzureManager) buildAsgFromSpec(spec string) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scaleToZeroSupported := scaleToZeroSupportedStandard
	if strings.EqualFold(m.config.VMType, vmTypeVMSS) {
		scaleToZeroSupported = scaleToZeroSupportedVMSS
	}
	s, err := dynamic.SpecFromString(spec, scaleToZeroSupported)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node group spec: %v", err)
	}
	switch m.config.VMType {
	case vmTypeStandard:
		return NewAgentPool(s, m)
	case vmTypeVMSS:
		return NewScaleSet(s, m)
	case vmTypeACS:
		fallthrough
	case vmTypeAKS:
		return NewContainerServiceAgentPool(s, m)
	default:
		return nil, fmt.Errorf("vmtype %s not supported", m.config.VMType)
	}
}
func (m *AzureManager) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if m.lastRefresh.Add(refreshInterval).After(time.Now()) {
		return nil
	}
	return m.forceRefresh()
}
func (m *AzureManager) forceRefresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := m.fetchAutoAsgs(); err != nil {
		klog.Errorf("Failed to fetch ASGs: %v", err)
		return err
	}
	m.lastRefresh = time.Now()
	klog.V(2).Infof("Refreshed ASG list, next refresh after %v", m.lastRefresh.Add(refreshInterval))
	return nil
}
func (m *AzureManager) fetchAutoAsgs() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	groups, err := m.getFilteredAutoscalingGroups(m.asgAutoDiscoverySpecs)
	if err != nil {
		return fmt.Errorf("cannot autodiscover ASGs: %s", err)
	}
	changed := false
	exists := make(map[string]bool)
	for _, asg := range groups {
		asgID := asg.Id()
		exists[asgID] = true
		if m.explicitlyConfigured[asgID] {
			klog.V(3).Infof("Ignoring explicitly configured ASG %s for autodiscovery.", asg.Id())
			continue
		}
		if m.RegisterAsg(asg) {
			klog.V(3).Infof("Autodiscovered ASG %s using tags %v", asg.Id(), m.asgAutoDiscoverySpecs)
			changed = true
		}
	}
	for _, asg := range m.getAsgs() {
		asgID := asg.Id()
		if !exists[asgID] && !m.explicitlyConfigured[asgID] {
			m.UnregisterAsg(asg)
			changed = true
		}
	}
	if changed {
		if err := m.regenerateCache(); err != nil {
			return err
		}
	}
	return nil
}
func (m *AzureManager) getAsgs() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.get()
}
func (m *AzureManager) RegisterAsg(asg cloudprovider.NodeGroup) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.Register(asg)
}
func (m *AzureManager) UnregisterAsg(asg cloudprovider.NodeGroup) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.Unregister(asg)
}
func (m *AzureManager) GetAsgForInstance(instance *azureRef) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return m.asgCache.FindForInstance(instance, m.config.VMType)
}
func (m *AzureManager) regenerateCache() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.asgCache.mutex.Lock()
	defer m.asgCache.mutex.Unlock()
	return m.asgCache.regenerate()
}
func (m *AzureManager) Cleanup() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.asgCache.Cleanup()
}
func (m *AzureManager) getFilteredAutoscalingGroups(filter []cloudprovider.LabelAutoDiscoveryConfig) (asgs []cloudprovider.NodeGroup, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(filter) == 0 {
		return nil, nil
	}
	switch m.config.VMType {
	case vmTypeVMSS:
		asgs, err = m.listScaleSets(filter)
	case vmTypeStandard:
		asgs, err = m.listAgentPools(filter)
	case vmTypeACS:
	case vmTypeAKS:
		return nil, nil
	default:
		err = fmt.Errorf("vmType %q not supported", m.config.VMType)
	}
	if err != nil {
		return nil, err
	}
	return asgs, nil
}
func (m *AzureManager) listScaleSets(filter []cloudprovider.LabelAutoDiscoveryConfig) (asgs []cloudprovider.NodeGroup, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := getContextWithCancel()
	defer cancel()
	result, err := m.azClient.virtualMachineScaleSetsClient.List(ctx, m.config.ResourceGroup)
	if err != nil {
		klog.Errorf("VirtualMachineScaleSetsClient.List for %v failed: %v", m.config.ResourceGroup, err)
		return nil, err
	}
	for _, scaleSet := range result {
		if len(filter) > 0 {
			if scaleSet.Tags == nil || len(scaleSet.Tags) == 0 {
				continue
			}
			if !matchDiscoveryConfig(scaleSet.Tags, filter) {
				continue
			}
		}
		spec := &dynamic.NodeGroupSpec{Name: *scaleSet.Name, MinSize: 1, MaxSize: -1, SupportScaleToZero: scaleToZeroSupportedVMSS}
		asg, _ := NewScaleSet(spec, m)
		asgs = append(asgs, asg)
	}
	return asgs, nil
}
func (m *AzureManager) listAgentPools(filter []cloudprovider.LabelAutoDiscoveryConfig) (asgs []cloudprovider.NodeGroup, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := getContextWithCancel()
	defer cancel()
	deploy, err := m.azClient.deploymentsClient.Get(ctx, m.config.ResourceGroup, m.config.Deployment)
	if err != nil {
		klog.Errorf("deploymentsClient.Get(%s, %s) failed: %v", m.config.ResourceGroup, m.config.Deployment, err)
		return nil, err
	}
	parameters := deploy.Properties.Parameters.(map[string]interface{})
	for k := range parameters {
		if k == "masterVMSize" || !strings.HasSuffix(k, "VMSize") {
			continue
		}
		poolName := strings.TrimRight(k, "VMSize")
		spec := &dynamic.NodeGroupSpec{Name: poolName, MinSize: 1, MaxSize: -1, SupportScaleToZero: scaleToZeroSupportedStandard}
		asg, _ := NewAgentPool(spec, m)
		asgs = append(asgs, asg)
	}
	return asgs, nil
}
