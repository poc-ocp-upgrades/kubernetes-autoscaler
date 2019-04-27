package ecs

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) DescribeInstanceTypes(request *DescribeInstanceTypesRequest) (response *DescribeInstanceTypesResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateDescribeInstanceTypesResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) DescribeInstanceTypesWithChan(request *DescribeInstanceTypesRequest) (<-chan *DescribeInstanceTypesResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *DescribeInstanceTypesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeInstanceTypes(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}
func (client *Client) DescribeInstanceTypesWithCallback(request *DescribeInstanceTypesRequest, callback func(response *DescribeInstanceTypesResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeInstanceTypesResponse
		var err error
		defer close(result)
		response, err = client.DescribeInstanceTypes(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

type DescribeInstanceTypesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId		requests.Integer	`position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount	string			`position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount		string			`position:"Query" name:"OwnerAccount"`
	InstanceTypeFamily	string			`position:"Query" name:"InstanceTypeFamily"`
	OwnerId			requests.Integer	`position:"Query" name:"OwnerId"`
}
type DescribeInstanceTypesResponse struct {
	*responses.BaseResponse
	RequestId	string					`json:"RequestId" xml:"RequestId"`
	InstanceTypes	InstanceTypesInDescribeInstanceTypes	`json:"InstanceTypes" xml:"InstanceTypes"`
}
type InstanceTypesInDescribeInstanceTypes struct {
	InstanceType []InstanceType `json:"InstanceType" xml:"InstanceType"`
}
type InstanceType struct {
	MemorySize		float64	`json:"MemorySize" xml:"MemorySize"`
	InstancePpsRx		int	`json:"InstancePpsRx" xml:"InstancePpsRx"`
	CpuCoreCount		int	`json:"CpuCoreCount" xml:"CpuCoreCount"`
	Cores			int	`json:"Cores" xml:"Cores"`
	Memory			int	`json:"Memory" xml:"Memory"`
	InstanceTypeId		string	`json:"InstanceTypeId" xml:"InstanceTypeId"`
	InstanceBandwidthRx	int	`json:"InstanceBandwidthRx" xml:"InstanceBandwidthRx"`
	InstanceType		string	`json:"InstanceType" xml:"InstanceType"`
	BaselineCredit		int	`json:"BaselineCredit" xml:"BaselineCredit"`
	EniQuantity		int	`json:"EniQuantity" xml:"EniQuantity"`
	Generation		string	`json:"Generation" xml:"Generation"`
	GPUAmount		int	`json:"GPUAmount" xml:"GPUAmount"`
	SupportIoOptimized	string	`json:"SupportIoOptimized" xml:"SupportIoOptimized"`
	InstanceTypeFamily	string	`json:"InstanceTypeFamily" xml:"InstanceTypeFamily"`
	InitialCredit		int	`json:"InitialCredit" xml:"InitialCredit"`
	InstancePpsTx		int	`json:"InstancePpsTx" xml:"InstancePpsTx"`
	LocalStorageAmount	int	`json:"LocalStorageAmount" xml:"LocalStorageAmount"`
	InstanceFamilyLevel	string	`json:"InstanceFamilyLevel" xml:"InstanceFamilyLevel"`
	LocalStorageCapacity	int	`json:"LocalStorageCapacity" xml:"LocalStorageCapacity"`
	GPUSpec			string	`json:"GPUSpec" xml:"GPUSpec"`
	LocalStorageCategory	string	`json:"LocalStorageCategory" xml:"LocalStorageCategory"`
	InstanceBandwidthTx	int	`json:"InstanceBandwidthTx" xml:"InstanceBandwidthTx"`
}

func CreateDescribeInstanceTypesRequest() (request *DescribeInstanceTypesRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &DescribeInstanceTypesRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeInstanceTypes", "ecs", "openAPI")
	return
}
func CreateDescribeInstanceTypesResponse() (response *DescribeInstanceTypesResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &DescribeInstanceTypesResponse{BaseResponse: &responses.BaseResponse{}}
	return
}
