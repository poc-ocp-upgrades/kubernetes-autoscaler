package azure

import (
	"testing"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

func newTestAzureManager(t *testing.T) *AzureManager {
	_logClusterCodePath()
	defer _logClusterCodePath()
	manager := &AzureManager{env: azure.PublicCloud, explicitlyConfigured: make(map[string]bool), config: &Config{ResourceGroup: "test", VMType: vmTypeVMSS}, azClient: &azClient{disksClient: &DisksClientMock{}, interfacesClient: &InterfacesClientMock{}, storageAccountsClient: &AccountsClientMock{}, deploymentsClient: &DeploymentsClientMock{FakeStore: make(map[string]resources.DeploymentExtended)}, virtualMachinesClient: &VirtualMachinesClientMock{FakeStore: make(map[string]map[string]compute.VirtualMachine)}, virtualMachineScaleSetsClient: &VirtualMachineScaleSetsClientMock{FakeStore: make(map[string]map[string]compute.VirtualMachineScaleSet)}, virtualMachineScaleSetVMsClient: &VirtualMachineScaleSetVMsClientMock{}}}
	cache, error := newAsgCache()
	assert.NoError(t, error)
	manager.asgCache = cache
	return manager
}
func newTestProvider(t *testing.T) *AzureCloudProvider {
	_logClusterCodePath()
	defer _logClusterCodePath()
	manager := newTestAzureManager(t)
	resourceLimiter := cloudprovider.NewResourceLimiter(map[string]int64{cloudprovider.ResourceNameCores: 1, cloudprovider.ResourceNameMemory: 10000000}, map[string]int64{cloudprovider.ResourceNameCores: 10, cloudprovider.ResourceNameMemory: 100000000})
	return &AzureCloudProvider{azureManager: manager, resourceLimiter: resourceLimiter}
}
func TestBuildAzureCloudProvider(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceLimiter := cloudprovider.NewResourceLimiter(map[string]int64{cloudprovider.ResourceNameCores: 1, cloudprovider.ResourceNameMemory: 10000000}, map[string]int64{cloudprovider.ResourceNameCores: 10, cloudprovider.ResourceNameMemory: 100000000})
	m := newTestAzureManager(t)
	_, err := BuildAzureCloudProvider(m, resourceLimiter)
	assert.NoError(t, err)
}
func TestName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := newTestProvider(t)
	assert.Equal(t, provider.Name(), "azure")
}
func TestNodeGroups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := newTestProvider(t)
	assert.Equal(t, len(provider.NodeGroups()), 0)
	registered := provider.azureManager.RegisterAsg(newTestScaleSet(provider.azureManager, "test-asg"))
	assert.True(t, registered)
	assert.Equal(t, len(provider.NodeGroups()), 1)
}
func TestNodeGroupForNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := newTestProvider(t)
	registered := provider.azureManager.RegisterAsg(newTestScaleSet(provider.azureManager, "test-asg"))
	assert.True(t, registered)
	assert.Equal(t, len(provider.NodeGroups()), 1)
	node := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "azure://" + fakeVirtualMachineScaleSetVMID}}
	group, err := provider.NodeGroupForNode(node)
	assert.NoError(t, err)
	assert.NotNil(t, group, "Group should not be nil")
	assert.Equal(t, group.Id(), "test-asg")
	assert.Equal(t, group.MinSize(), 1)
	assert.Equal(t, group.MaxSize(), 5)
	nodeNotInGroup := &apiv1.Node{Spec: apiv1.NodeSpec{ProviderID: "azure:///subscriptions/subscripion/resourceGroups/test-resource-group/providers/Microsoft.Compute/virtualMachines/test-instance-id-not-in-group"}}
	group, err = provider.NodeGroupForNode(nodeNotInGroup)
	assert.NoError(t, err)
	assert.Nil(t, group)
}
