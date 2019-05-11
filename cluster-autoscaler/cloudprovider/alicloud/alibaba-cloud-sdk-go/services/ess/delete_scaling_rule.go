package ess

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) DeleteScalingRule(request *DeleteScalingRuleRequest) (response *DeleteScalingRuleResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateDeleteScalingRuleResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) DeleteScalingRuleWithChan(request *DeleteScalingRuleRequest) (<-chan *DeleteScalingRuleResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *DeleteScalingRuleResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteScalingRule(request)
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
func (client *Client) DeleteScalingRuleWithCallback(request *DeleteScalingRuleRequest, callback func(response *DeleteScalingRuleResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteScalingRuleResponse
		var err error
		defer close(result)
		response, err = client.DeleteScalingRule(request)
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

type DeleteScalingRuleRequest struct {
	*requests.RpcRequest
	ResourceOwnerAccount	string				`position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount			string				`position:"Query" name:"OwnerAccount"`
	OwnerId					requests.Integer	`position:"Query" name:"OwnerId"`
	ScalingRuleId			string				`position:"Query" name:"ScalingRuleId"`
}
type DeleteScalingRuleResponse struct {
	*responses.BaseResponse
	RequestId	string	`json:"RequestId" xml:"RequestId"`
}

func CreateDeleteScalingRuleRequest() (request *DeleteScalingRuleRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &DeleteScalingRuleRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ess", "2014-08-28", "DeleteScalingRule", "ess", "openAPI")
	return
}
func CreateDeleteScalingRuleResponse() (response *DeleteScalingRuleResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &DeleteScalingRuleResponse{BaseResponse: &responses.BaseResponse{}}
	return
}
