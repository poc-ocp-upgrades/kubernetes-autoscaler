package ess

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) ModifyScalingGroup(request *ModifyScalingGroupRequest) (response *ModifyScalingGroupResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateModifyScalingGroupResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) ModifyScalingGroupWithChan(request *ModifyScalingGroupRequest) (<-chan *ModifyScalingGroupResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *ModifyScalingGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyScalingGroup(request)
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
func (client *Client) ModifyScalingGroupWithCallback(request *ModifyScalingGroupRequest, callback func(response *ModifyScalingGroupResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyScalingGroupResponse
		var err error
		defer close(result)
		response, err = client.ModifyScalingGroup(request)
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

type ModifyScalingGroupRequest struct {
	*requests.RpcRequest
	ResourceOwnerId					requests.Integer	`position:"Query" name:"ResourceOwnerId"`
	HealthCheckType					string				`position:"Query" name:"HealthCheckType"`
	LaunchTemplateId				string				`position:"Query" name:"LaunchTemplateId"`
	ResourceOwnerAccount			string				`position:"Query" name:"ResourceOwnerAccount"`
	ScalingGroupName				string				`position:"Query" name:"ScalingGroupName"`
	ScalingGroupId					string				`position:"Query" name:"ScalingGroupId"`
	OwnerAccount					string				`position:"Query" name:"OwnerAccount"`
	ActiveScalingConfigurationId	string				`position:"Query" name:"ActiveScalingConfigurationId"`
	MinSize							requests.Integer	`position:"Query" name:"MinSize"`
	OwnerId							requests.Integer	`position:"Query" name:"OwnerId"`
	LaunchTemplateVersion			string				`position:"Query" name:"LaunchTemplateVersion"`
	MaxSize							requests.Integer	`position:"Query" name:"MaxSize"`
	DefaultCooldown					requests.Integer	`position:"Query" name:"DefaultCooldown"`
	RemovalPolicy1					string				`position:"Query" name:"RemovalPolicy.1"`
	RemovalPolicy2					string				`position:"Query" name:"RemovalPolicy.2"`
}
type ModifyScalingGroupResponse struct {
	*responses.BaseResponse
	RequestId	string	`json:"RequestId" xml:"RequestId"`
}

func CreateModifyScalingGroupRequest() (request *ModifyScalingGroupRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &ModifyScalingGroupRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ess", "2014-08-28", "ModifyScalingGroup", "ess", "openAPI")
	return
}
func CreateModifyScalingGroupResponse() (response *ModifyScalingGroupResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &ModifyScalingGroupResponse{BaseResponse: &responses.BaseResponse{}}
	return
}
