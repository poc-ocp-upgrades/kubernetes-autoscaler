package ess

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) DescribeScalingInstances(request *DescribeScalingInstancesRequest) (response *DescribeScalingInstancesResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateDescribeScalingInstancesResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) DescribeScalingInstancesWithChan(request *DescribeScalingInstancesRequest) (<-chan *DescribeScalingInstancesResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *DescribeScalingInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeScalingInstances(request)
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
func (client *Client) DescribeScalingInstancesWithCallback(request *DescribeScalingInstancesRequest, callback func(response *DescribeScalingInstancesResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeScalingInstancesResponse
		var err error
		defer close(result)
		response, err = client.DescribeScalingInstances(request)
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

type DescribeScalingInstancesRequest struct {
	*requests.RpcRequest
	InstanceId10			string				`position:"Query" name:"InstanceId.10"`
	ResourceOwnerId			requests.Integer	`position:"Query" name:"ResourceOwnerId"`
	InstanceId12			string				`position:"Query" name:"InstanceId.12"`
	InstanceId11			string				`position:"Query" name:"InstanceId.11"`
	ScalingGroupId			string				`position:"Query" name:"ScalingGroupId"`
	LifecycleState			string				`position:"Query" name:"LifecycleState"`
	CreationType			string				`position:"Query" name:"CreationType"`
	PageNumber				requests.Integer	`position:"Query" name:"PageNumber"`
	PageSize				requests.Integer	`position:"Query" name:"PageSize"`
	InstanceId20			string				`position:"Query" name:"InstanceId.20"`
	InstanceId1				string				`position:"Query" name:"InstanceId.1"`
	InstanceId3				string				`position:"Query" name:"InstanceId.3"`
	ResourceOwnerAccount	string				`position:"Query" name:"ResourceOwnerAccount"`
	InstanceId2				string				`position:"Query" name:"InstanceId.2"`
	InstanceId5				string				`position:"Query" name:"InstanceId.5"`
	InstanceId4				string				`position:"Query" name:"InstanceId.4"`
	OwnerAccount			string				`position:"Query" name:"OwnerAccount"`
	InstanceId7				string				`position:"Query" name:"InstanceId.7"`
	InstanceId6				string				`position:"Query" name:"InstanceId.6"`
	InstanceId9				string				`position:"Query" name:"InstanceId.9"`
	InstanceId8				string				`position:"Query" name:"InstanceId.8"`
	OwnerId					requests.Integer	`position:"Query" name:"OwnerId"`
	ScalingConfigurationId	string				`position:"Query" name:"ScalingConfigurationId"`
	HealthStatus			string				`position:"Query" name:"HealthStatus"`
	InstanceId18			string				`position:"Query" name:"InstanceId.18"`
	InstanceId17			string				`position:"Query" name:"InstanceId.17"`
	InstanceId19			string				`position:"Query" name:"InstanceId.19"`
	InstanceId14			string				`position:"Query" name:"InstanceId.14"`
	InstanceId13			string				`position:"Query" name:"InstanceId.13"`
	InstanceId16			string				`position:"Query" name:"InstanceId.16"`
	InstanceId15			string				`position:"Query" name:"InstanceId.15"`
}
type DescribeScalingInstancesResponse struct {
	*responses.BaseResponse
	TotalCount			int					`json:"TotalCount" xml:"TotalCount"`
	PageNumber			int					`json:"PageNumber" xml:"PageNumber"`
	PageSize			int					`json:"PageSize" xml:"PageSize"`
	RequestId			string				`json:"RequestId" xml:"RequestId"`
	ScalingInstances	ScalingInstances	`json:"ScalingInstances" xml:"ScalingInstances"`
}

func CreateDescribeScalingInstancesRequest() (request *DescribeScalingInstancesRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &DescribeScalingInstancesRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ess", "2014-08-28", "DescribeScalingInstances", "ess", "openAPI")
	return
}
func CreateDescribeScalingInstancesResponse() (response *DescribeScalingInstancesResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &DescribeScalingInstancesResponse{BaseResponse: &responses.BaseResponse{}}
	return
}

type ScalingInstances struct {
	ScalingInstance []ScalingInstance `json:"ScalingInstance" xml:"ScalingInstance"`
}
type ScalingInstance struct {
	InstanceId				string	`json:"InstanceId" xml:"InstanceId"`
	ScalingConfigurationId	string	`json:"ScalingConfigurationId" xml:"ScalingConfigurationId"`
	ScalingGroupId			string	`json:"ScalingGroupId" xml:"ScalingGroupId"`
	HealthStatus			string	`json:"HealthStatus" xml:"HealthStatus"`
	LoadBalancerWeight		int		`json:"LoadBalancerWeight" xml:"LoadBalancerWeight"`
	LifecycleState			string	`json:"LifecycleState" xml:"LifecycleState"`
	CreationTime			string	`json:"CreationTime" xml:"CreationTime"`
	CreationType			string	`json:"CreationType" xml:"CreationType"`
}
