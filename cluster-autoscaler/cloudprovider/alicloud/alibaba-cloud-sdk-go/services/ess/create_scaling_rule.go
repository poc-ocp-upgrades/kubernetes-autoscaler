package ess

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) CreateScalingRule(request *CreateScalingRuleRequest) (response *CreateScalingRuleResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateCreateScalingRuleResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) CreateScalingRuleWithChan(request *CreateScalingRuleRequest) (<-chan *CreateScalingRuleResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *CreateScalingRuleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateScalingRule(request)
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
func (client *Client) CreateScalingRuleWithCallback(request *CreateScalingRuleRequest, callback func(response *CreateScalingRuleResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateScalingRuleResponse
		var err error
		defer close(result)
		response, err = client.CreateScalingRule(request)
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

type CreateScalingRuleRequest struct {
	*requests.RpcRequest
	ScalingRuleName		string			`position:"Query" name:"ScalingRuleName"`
	ResourceOwnerAccount	string			`position:"Query" name:"ResourceOwnerAccount"`
	AdjustmentValue		requests.Integer	`position:"Query" name:"AdjustmentValue"`
	ScalingGroupId		string			`position:"Query" name:"ScalingGroupId"`
	OwnerAccount		string			`position:"Query" name:"OwnerAccount"`
	Cooldown		requests.Integer	`position:"Query" name:"Cooldown"`
	AdjustmentType		string			`position:"Query" name:"AdjustmentType"`
	OwnerId			requests.Integer	`position:"Query" name:"OwnerId"`
}
type CreateScalingRuleResponse struct {
	*responses.BaseResponse
	ScalingRuleId	string	`json:"ScalingRuleId" xml:"ScalingRuleId"`
	ScalingRuleAri	string	`json:"ScalingRuleAri" xml:"ScalingRuleAri"`
	RequestId	string	`json:"RequestId" xml:"RequestId"`
}

func CreateCreateScalingRuleRequest() (request *CreateScalingRuleRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &CreateScalingRuleRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ess", "2014-08-28", "CreateScalingRule", "ess", "openAPI")
	return
}
func CreateCreateScalingRuleResponse() (response *CreateScalingRuleResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &CreateScalingRuleResponse{BaseResponse: &responses.BaseResponse{}}
	return
}
