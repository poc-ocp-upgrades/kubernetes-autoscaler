package ess

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) RemoveInstances(request *RemoveInstancesRequest) (response *RemoveInstancesResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateRemoveInstancesResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) RemoveInstancesWithChan(request *RemoveInstancesRequest) (<-chan *RemoveInstancesResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *RemoveInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveInstances(request)
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
func (client *Client) RemoveInstancesWithCallback(request *RemoveInstancesRequest, callback func(response *RemoveInstancesResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveInstancesResponse
		var err error
		defer close(result)
		response, err = client.RemoveInstances(request)
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

type RemoveInstancesRequest struct {
	*requests.RpcRequest
	InstanceId10		string			`position:"Query" name:"InstanceId.10"`
	ResourceOwnerId		requests.Integer	`position:"Query" name:"ResourceOwnerId"`
	InstanceId12		string			`position:"Query" name:"InstanceId.12"`
	InstanceId11		string			`position:"Query" name:"InstanceId.11"`
	ScalingGroupId		string			`position:"Query" name:"ScalingGroupId"`
	InstanceId20		string			`position:"Query" name:"InstanceId.20"`
	InstanceId1		string			`position:"Query" name:"InstanceId.1"`
	InstanceId3		string			`position:"Query" name:"InstanceId.3"`
	ResourceOwnerAccount	string			`position:"Query" name:"ResourceOwnerAccount"`
	InstanceId2		string			`position:"Query" name:"InstanceId.2"`
	InstanceId5		string			`position:"Query" name:"InstanceId.5"`
	InstanceId4		string			`position:"Query" name:"InstanceId.4"`
	OwnerAccount		string			`position:"Query" name:"OwnerAccount"`
	InstanceId7		string			`position:"Query" name:"InstanceId.7"`
	InstanceId6		string			`position:"Query" name:"InstanceId.6"`
	InstanceId9		string			`position:"Query" name:"InstanceId.9"`
	InstanceId8		string			`position:"Query" name:"InstanceId.8"`
	OwnerId			requests.Integer	`position:"Query" name:"OwnerId"`
	InstanceId18		string			`position:"Query" name:"InstanceId.18"`
	InstanceId17		string			`position:"Query" name:"InstanceId.17"`
	InstanceId19		string			`position:"Query" name:"InstanceId.19"`
	InstanceId14		string			`position:"Query" name:"InstanceId.14"`
	InstanceId13		string			`position:"Query" name:"InstanceId.13"`
	InstanceId16		string			`position:"Query" name:"InstanceId.16"`
	InstanceId15		string			`position:"Query" name:"InstanceId.15"`
}
type RemoveInstancesResponse struct {
	*responses.BaseResponse
	ScalingActivityId	string	`json:"ScalingActivityId" xml:"ScalingActivityId"`
	RequestId		string	`json:"RequestId" xml:"RequestId"`
}

func CreateRemoveInstancesRequest() (request *RemoveInstancesRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &RemoveInstancesRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ess", "2014-08-28", "RemoveInstances", "ess", "openAPI")
	return
}
func CreateRemoveInstancesResponse() (response *RemoveInstancesResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &RemoveInstancesResponse{BaseResponse: &responses.BaseResponse{}}
	return
}
