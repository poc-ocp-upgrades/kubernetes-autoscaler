package azure

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	"k8s.io/klog"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
)

type ScaleSet struct {
	azureRef
	manager		*AzureManager
	minSize		int
	maxSize		int
	mutex		sync.Mutex
	lastRefresh	time.Time
	curSize		int64
}

func NewScaleSet(spec *dynamic.NodeGroupSpec, az *AzureManager) (*ScaleSet, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scaleSet := &ScaleSet{azureRef: azureRef{Name: spec.Name}, minSize: spec.MinSize, maxSize: spec.MaxSize, manager: az, curSize: -1}
	return scaleSet, nil
}
func (scaleSet *ScaleSet) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return scaleSet.minSize
}
func (scaleSet *ScaleSet) Exist() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func (scaleSet *ScaleSet) Create() (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrAlreadyExist
}
func (scaleSet *ScaleSet) Delete() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cloudprovider.ErrNotImplemented
}
func (scaleSet *ScaleSet) Autoprovisioned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (scaleSet *ScaleSet) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return scaleSet.maxSize
}
func (scaleSet *ScaleSet) getVMSSInfo() (compute.VirtualMachineScaleSet, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := getContextWithCancel()
	defer cancel()
	resourceGroup := scaleSet.manager.config.ResourceGroup
	setInfo, err := scaleSet.manager.azClient.virtualMachineScaleSetsClient.Get(ctx, resourceGroup, scaleSet.Name)
	if err != nil {
		return compute.VirtualMachineScaleSet{}, err
	}
	return setInfo, nil
}
func (scaleSet *ScaleSet) getCurSize() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scaleSet.mutex.Lock()
	defer scaleSet.mutex.Unlock()
	if scaleSet.lastRefresh.Add(15 * time.Second).After(time.Now()) {
		return scaleSet.curSize, nil
	}
	klog.V(5).Infof("Get scale set size for %q", scaleSet.Name)
	set, err := scaleSet.getVMSSInfo()
	if err != nil {
		return -1, err
	}
	klog.V(5).Infof("Getting scale set (%q) capacity: %d\n", scaleSet.Name, *set.Sku.Capacity)
	scaleSet.curSize = *set.Sku.Capacity
	scaleSet.lastRefresh = time.Now()
	return scaleSet.curSize, nil
}
func (scaleSet *ScaleSet) GetScaleSetSize() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return scaleSet.getCurSize()
}
func (scaleSet *ScaleSet) SetScaleSetSize(size int64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scaleSet.mutex.Lock()
	defer scaleSet.mutex.Unlock()
	resourceGroup := scaleSet.manager.config.ResourceGroup
	op, err := scaleSet.getVMSSInfo()
	if err != nil {
		return err
	}
	op.Sku.Capacity = &size
	op.VirtualMachineScaleSetProperties.ProvisioningState = nil
	updateCtx, updateCancel := getContextWithCancel()
	defer updateCancel()
	klog.V(3).Infof("Waiting for virtualMachineScaleSetsClient.CreateOrUpdate(%s)", scaleSet.Name)
	resp, err := scaleSet.manager.azClient.virtualMachineScaleSetsClient.CreateOrUpdate(updateCtx, resourceGroup, scaleSet.Name, op)
	isSuccess, realError := isSuccessHTTPResponse(resp, err)
	if isSuccess {
		klog.V(3).Infof("virtualMachineScaleSetsClient.CreateOrUpdate(%s) success", scaleSet.Name)
		scaleSet.curSize = size
		scaleSet.lastRefresh = time.Now()
		return nil
	}
	klog.Errorf("virtualMachineScaleSetsClient.CreateOrUpdate for scale set %q failed: %v", scaleSet.Name, realError)
	return realError
}
func (scaleSet *ScaleSet) TargetSize() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	size, err := scaleSet.GetScaleSetSize()
	return int(size), err
}
func (scaleSet *ScaleSet) IncreaseSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta <= 0 {
		return fmt.Errorf("size increase must be positive")
	}
	size, err := scaleSet.GetScaleSetSize()
	if err != nil {
		return err
	}
	if int(size)+delta > scaleSet.MaxSize() {
		return fmt.Errorf("size increase too large - desired:%d max:%d", int(size)+delta, scaleSet.MaxSize())
	}
	return scaleSet.SetScaleSetSize(size + int64(delta))
}
func (scaleSet *ScaleSet) GetScaleSetVms() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := getContextWithCancel()
	defer cancel()
	resourceGroup := scaleSet.manager.config.ResourceGroup
	vmList, err := scaleSet.manager.azClient.virtualMachineScaleSetVMsClient.List(ctx, resourceGroup, scaleSet.Name, "", "", "")
	if err != nil {
		klog.Errorf("VirtualMachineScaleSetVMsClient.List failed for %s: %v", scaleSet.Name, err)
		return nil, err
	}
	allVMs := make([]string, 0)
	for _, vm := range vmList {
		if len(*vm.ID) == 0 {
			continue
		}
		allVMs = append(allVMs, *vm.ID)
	}
	return allVMs, nil
}
func (scaleSet *ScaleSet) DecreaseTargetSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if delta >= 0 {
		return fmt.Errorf("size decrease size must be negative")
	}
	size, err := scaleSet.GetScaleSetSize()
	if err != nil {
		return err
	}
	nodes, err := scaleSet.Nodes()
	if err != nil {
		return err
	}
	if int(size)+delta < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d", size, delta, len(nodes))
	}
	return scaleSet.SetScaleSetSize(size + int64(delta))
}
func (scaleSet *ScaleSet) Belongs(node *apiv1.Node) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.V(6).Infof("Check if node belongs to this scale set: scaleset:%v, node:%v\n", scaleSet, node)
	ref := &azureRef{Name: node.Spec.ProviderID}
	targetAsg, err := scaleSet.manager.GetAsgForInstance(ref)
	if err != nil {
		return false, err
	}
	if targetAsg == nil {
		return false, fmt.Errorf("%s doesn't belong to a known scale set", node.Name)
	}
	if !strings.EqualFold(targetAsg.Id(), scaleSet.Id()) {
		return false, nil
	}
	return true, nil
}
func (scaleSet *ScaleSet) DeleteInstances(instances []*azureRef) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(instances) == 0 {
		return nil
	}
	klog.V(3).Infof("Deleting vmss instances %q", instances)
	commonAsg, err := scaleSet.manager.GetAsgForInstance(instances[0])
	if err != nil {
		return err
	}
	instanceIDs := []string{}
	for _, instance := range instances {
		asg, err := scaleSet.manager.GetAsgForInstance(instance)
		if err != nil {
			return err
		}
		if !strings.EqualFold(asg.Id(), commonAsg.Id()) {
			return fmt.Errorf("cannot delete instance (%s) which don't belong to the same Scale Set (%q)", instance.Name, commonAsg)
		}
		instanceID, err := getLastSegment(instance.Name)
		if err != nil {
			klog.Errorf("getLastSegment failed with error: %v", err)
			return err
		}
		instanceIDs = append(instanceIDs, instanceID)
	}
	requiredIds := &compute.VirtualMachineScaleSetVMInstanceRequiredIDs{InstanceIds: &instanceIDs}
	ctx, cancel := getContextWithCancel()
	defer cancel()
	resourceGroup := scaleSet.manager.config.ResourceGroup
	_, err = scaleSet.manager.azClient.virtualMachineScaleSetsClient.DeleteInstances(ctx, resourceGroup, commonAsg.Id(), *requiredIds)
	return err
}
func (scaleSet *ScaleSet) DeleteNodes(nodes []*apiv1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	klog.V(8).Infof("Delete nodes requested: %q\n", nodes)
	size, err := scaleSet.GetScaleSetSize()
	if err != nil {
		return err
	}
	if int(size) <= scaleSet.MinSize() {
		return fmt.Errorf("min size reached, nodes will not be deleted")
	}
	refs := make([]*azureRef, 0, len(nodes))
	for _, node := range nodes {
		belongs, err := scaleSet.Belongs(node)
		if err != nil {
			return err
		}
		if belongs != true {
			return fmt.Errorf("%s belongs to a different asg than %s", node.Name, scaleSet.Id())
		}
		ref := &azureRef{Name: node.Spec.ProviderID}
		refs = append(refs, ref)
	}
	return scaleSet.DeleteInstances(refs)
}
func (scaleSet *ScaleSet) Id() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return scaleSet.Name
}
func (scaleSet *ScaleSet) Debug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s (%d:%d)", scaleSet.Id(), scaleSet.MinSize(), scaleSet.MaxSize())
}
func buildInstanceOS(template compute.VirtualMachineScaleSet) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	instanceOS := cloudprovider.DefaultOS
	if template.VirtualMachineProfile != nil && template.VirtualMachineProfile.OsProfile != nil && template.VirtualMachineProfile.OsProfile.WindowsConfiguration != nil {
		instanceOS = "windows"
	}
	return instanceOS
}
func buildGenericLabels(template compute.VirtualMachineScaleSet, nodeName string) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(map[string]string)
	result[kubeletapis.LabelArch] = cloudprovider.DefaultArch
	result[kubeletapis.LabelOS] = buildInstanceOS(template)
	result[kubeletapis.LabelInstanceType] = *template.Sku.Name
	result[kubeletapis.LabelZoneRegion] = *template.Location
	result[kubeletapis.LabelZoneFailureDomain] = "0"
	result[kubeletapis.LabelHostname] = nodeName
	return result
}
func (scaleSet *ScaleSet) buildNodeFromTemplate(template compute.VirtualMachineScaleSet) (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := apiv1.Node{}
	nodeName := fmt.Sprintf("%s-asg-%d", scaleSet.Name, rand.Int63())
	node.ObjectMeta = metav1.ObjectMeta{Name: nodeName, SelfLink: fmt.Sprintf("/api/v1/nodes/%s", nodeName), Labels: map[string]string{}}
	node.Status = apiv1.NodeStatus{Capacity: apiv1.ResourceList{}}
	vmssType := InstanceTypes[*template.Sku.Name]
	if vmssType == nil {
		return nil, fmt.Errorf("instance type %q not supported", *template.Sku.Name)
	}
	node.Status.Capacity[apiv1.ResourcePods] = *resource.NewQuantity(110, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceCPU] = *resource.NewQuantity(vmssType.VCPU, resource.DecimalSI)
	node.Status.Capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(vmssType.GPU, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceMemory] = *resource.NewQuantity(vmssType.MemoryMb*1024*1024, resource.DecimalSI)
	node.Status.Allocatable = node.Status.Capacity
	if template.Tags != nil {
		for k, v := range template.Tags {
			if v != nil {
				node.Labels[k] = *v
			} else {
				node.Labels[k] = ""
			}
		}
	}
	node.Labels = cloudprovider.JoinStringMaps(node.Labels, buildGenericLabels(template, nodeName))
	node.Status.Conditions = cloudprovider.BuildReadyConditions()
	return &node, nil
}
func (scaleSet *ScaleSet) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	template, err := scaleSet.getVMSSInfo()
	if err != nil {
		return nil, err
	}
	node, err := scaleSet.buildNodeFromTemplate(template)
	if err != nil {
		return nil, err
	}
	nodeInfo := schedulercache.NewNodeInfo(cloudprovider.BuildKubeProxy(scaleSet.Name))
	nodeInfo.SetNode(node)
	return nodeInfo, nil
}
func (scaleSet *ScaleSet) Nodes() ([]cloudprovider.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scaleSet.mutex.Lock()
	defer scaleSet.mutex.Unlock()
	vms, err := scaleSet.GetScaleSetVms()
	if err != nil {
		return nil, err
	}
	instances := make([]cloudprovider.Instance, 0, len(vms))
	for i := range vms {
		name := "azure://" + strings.ToLower(vms[i])
		instances = append(instances, cloudprovider.Instance{Id: name})
	}
	return instances, nil
}
