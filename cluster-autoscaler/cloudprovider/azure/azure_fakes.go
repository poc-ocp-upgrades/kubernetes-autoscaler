package azure

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/stretchr/testify/mock"
)

const (
	fakeVirtualMachineScaleSetVMID = "/subscriptions/test-subscription-id/resourcegroups/test-asg/providers/microsoft.compute/virtualmachinescalesets/agents/virtualmachines/0"
)

type VirtualMachineScaleSetsClientMock struct {
	mock.Mock
	mutex		sync.Mutex
	FakeStore	map[string]map[string]compute.VirtualMachineScaleSet
}

func (client *VirtualMachineScaleSetsClientMock) Get(ctx context.Context, resourceGroupName string, vmScaleSetName string) (result compute.VirtualMachineScaleSet, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	capacity := int64(2)
	properties := compute.VirtualMachineScaleSetProperties{}
	return compute.VirtualMachineScaleSet{Name: &vmScaleSetName, Sku: &compute.Sku{Capacity: &capacity}, VirtualMachineScaleSetProperties: &properties}, nil
}
func (client *VirtualMachineScaleSetsClientMock) CreateOrUpdate(ctx context.Context, resourceGroupName string, VMScaleSetName string, parameters compute.VirtualMachineScaleSet) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if _, ok := client.FakeStore[resourceGroupName]; !ok {
		client.FakeStore[resourceGroupName] = make(map[string]compute.VirtualMachineScaleSet)
	}
	client.FakeStore[resourceGroupName][VMScaleSetName] = parameters
	return &http.Response{StatusCode: http.StatusOK}, nil
}
func (client *VirtualMachineScaleSetsClientMock) DeleteInstances(ctx context.Context, resourceGroupName string, vmScaleSetName string, vmInstanceIDs compute.VirtualMachineScaleSetVMInstanceRequiredIDs) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := client.Called(resourceGroupName, vmScaleSetName, vmInstanceIDs)
	return nil, args.Error(1)
}
func (client *VirtualMachineScaleSetsClientMock) List(ctx context.Context, resourceGroupName string) (result []compute.VirtualMachineScaleSet, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.mutex.Lock()
	defer client.mutex.Unlock()
	result = []compute.VirtualMachineScaleSet{}
	if _, ok := client.FakeStore[resourceGroupName]; ok {
		for _, v := range client.FakeStore[resourceGroupName] {
			result = append(result, v)
		}
	}
	return result, nil
}

type VirtualMachineScaleSetVMsClientMock struct{ mock.Mock }

func (m *VirtualMachineScaleSetVMsClientMock) Get(ctx context.Context, resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVM, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ID := fakeVirtualMachineScaleSetVMID
	vmID := "123E4567-E89B-12D3-A456-426655440000"
	properties := compute.VirtualMachineScaleSetVMProperties{VMID: &vmID}
	return compute.VirtualMachineScaleSetVM{ID: &ID, InstanceID: &instanceID, VirtualMachineScaleSetVMProperties: &properties}, nil
}
func (m *VirtualMachineScaleSetVMsClientMock) List(ctx context.Context, resourceGroupName string, virtualMachineScaleSetName string, filter string, selectParameter string, expand string) (result []compute.VirtualMachineScaleSetVM, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ID := fakeVirtualMachineScaleSetVMID
	instanceID := "0"
	vmID := "123E4567-E89B-12D3-A456-426655440000"
	properties := compute.VirtualMachineScaleSetVMProperties{VMID: &vmID}
	result = append(result, compute.VirtualMachineScaleSetVM{ID: &ID, InstanceID: &instanceID, VirtualMachineScaleSetVMProperties: &properties})
	return result, nil
}

type VirtualMachinesClientMock struct {
	mock.Mock
	mutex		sync.Mutex
	FakeStore	map[string]map[string]compute.VirtualMachine
}

func (m *VirtualMachinesClientMock) Get(ctx context.Context, resourceGroupName string, VMName string, expand compute.InstanceViewTypes) (result compute.VirtualMachine, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.FakeStore[resourceGroupName]; ok {
		if entity, ok := m.FakeStore[resourceGroupName][VMName]; ok {
			return entity, nil
		}
	}
	return result, autorest.DetailedError{StatusCode: http.StatusNotFound, Message: "Not such VM"}
}
func (m *VirtualMachinesClientMock) List(ctx context.Context, resourceGroupName string) (result []compute.VirtualMachine, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.FakeStore[resourceGroupName]; ok {
		for _, v := range m.FakeStore[resourceGroupName] {
			result = append(result, v)
		}
	}
	return result, nil
}
func (m *VirtualMachinesClientMock) Delete(ctx context.Context, resourceGroupName string, VMName string) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(resourceGroupName, VMName)
	return nil, args.Error(1)
}

type InterfacesClientMock struct{ mock.Mock }

func (m *InterfacesClientMock) Delete(ctx context.Context, resourceGroupName string, networkInterfaceName string) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(resourceGroupName, networkInterfaceName)
	return nil, args.Error(1)
}

type DisksClientMock struct{ mock.Mock }

func (m *DisksClientMock) Delete(ctx context.Context, resourceGroupName string, diskName string) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(resourceGroupName, diskName)
	return nil, args.Error(1)
}

type AccountsClientMock struct{ mock.Mock }

func (m *AccountsClientMock) ListKeys(ctx context.Context, resourceGroupName string, accountName string) (result storage.AccountListKeysResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(resourceGroupName, accountName)
	return storage.AccountListKeysResult{}, args.Error(1)
}

type DeploymentsClientMock struct {
	mock.Mock
	mutex		sync.Mutex
	FakeStore	map[string]resources.DeploymentExtended
}

func (m *DeploymentsClientMock) Get(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExtended, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	deploy, ok := m.FakeStore[deploymentName]
	if !ok {
		return result, fmt.Errorf("deployment not found")
	}
	return deploy, nil
}
func (m *DeploymentsClientMock) ExportTemplate(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExportResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	deploy, ok := m.FakeStore[deploymentName]
	if !ok {
		return result, fmt.Errorf("deployment not found")
	}
	return resources.DeploymentExportResult{Template: deploy.Properties.Template}, nil
}
func (m *DeploymentsClientMock) CreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.Deployment) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	deploy, ok := m.FakeStore[deploymentName]
	if !ok {
		deploy = resources.DeploymentExtended{Properties: &resources.DeploymentPropertiesExtended{}}
		m.FakeStore[deploymentName] = deploy
	}
	deploy.Properties.Parameters = parameters.Properties.Parameters
	deploy.Properties.Template = parameters.Properties.Template
	return nil, nil
}
