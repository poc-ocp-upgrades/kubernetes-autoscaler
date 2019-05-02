package ess

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) ExecuteScalingRule(request *ExecuteScalingRuleRequest) (response *ExecuteScalingRuleResponse, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = CreateExecuteScalingRuleResponse()
 err = client.DoAction(request, response)
 return
}
func (client *Client) ExecuteScalingRuleWithChan(request *ExecuteScalingRuleRequest) (<-chan *ExecuteScalingRuleResponse, <-chan error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 responseChan := make(chan *ExecuteScalingRuleResponse, 1)
 errChan := make(chan error, 1)
 err := client.AddAsyncTask(func() {
  defer close(responseChan)
  defer close(errChan)
  response, err := client.ExecuteScalingRule(request)
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
func (client *Client) ExecuteScalingRuleWithCallback(request *ExecuteScalingRuleRequest, callback func(response *ExecuteScalingRuleResponse, err error)) <-chan int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := make(chan int, 1)
 err := client.AddAsyncTask(func() {
  var response *ExecuteScalingRuleResponse
  var err error
  defer close(result)
  response, err = client.ExecuteScalingRule(request)
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

type ExecuteScalingRuleRequest struct {
 *requests.RpcRequest
 ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
 ScalingRuleAri       string           `position:"Query" name:"ScalingRuleAri"`
 ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
 ClientToken          string           `position:"Query" name:"ClientToken"`
 OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
 OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}
type ExecuteScalingRuleResponse struct {
 *responses.BaseResponse
 ScalingActivityId string `json:"ScalingActivityId" xml:"ScalingActivityId"`
 RequestId         string `json:"RequestId" xml:"RequestId"`
}

func CreateExecuteScalingRuleRequest() (request *ExecuteScalingRuleRequest) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 request = &ExecuteScalingRuleRequest{RpcRequest: &requests.RpcRequest{}}
 request.InitWithApiInfo("Ess", "2014-08-28", "ExecuteScalingRule", "ess", "openAPI")
 return
}
func CreateExecuteScalingRuleResponse() (response *ExecuteScalingRuleResponse) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = &ExecuteScalingRuleResponse{BaseResponse: &responses.BaseResponse{}}
 return
}
