package ess

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) ModifyScalingRule(request *ModifyScalingRuleRequest) (response *ModifyScalingRuleResponse, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = CreateModifyScalingRuleResponse()
 err = client.DoAction(request, response)
 return
}
func (client *Client) ModifyScalingRuleWithChan(request *ModifyScalingRuleRequest) (<-chan *ModifyScalingRuleResponse, <-chan error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 responseChan := make(chan *ModifyScalingRuleResponse, 1)
 errChan := make(chan error, 1)
 err := client.AddAsyncTask(func() {
  defer close(responseChan)
  defer close(errChan)
  response, err := client.ModifyScalingRule(request)
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
func (client *Client) ModifyScalingRuleWithCallback(request *ModifyScalingRuleRequest, callback func(response *ModifyScalingRuleResponse, err error)) <-chan int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := make(chan int, 1)
 err := client.AddAsyncTask(func() {
  var response *ModifyScalingRuleResponse
  var err error
  defer close(result)
  response, err = client.ModifyScalingRule(request)
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

type ModifyScalingRuleRequest struct {
 *requests.RpcRequest
 ScalingRuleName      string           `position:"Query" name:"ScalingRuleName"`
 ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
 ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
 AdjustmentValue      requests.Integer `position:"Query" name:"AdjustmentValue"`
 OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
 Cooldown             requests.Integer `position:"Query" name:"Cooldown"`
 AdjustmentType       string           `position:"Query" name:"AdjustmentType"`
 OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
 ScalingRuleId        string           `position:"Query" name:"ScalingRuleId"`
}
type ModifyScalingRuleResponse struct {
 *responses.BaseResponse
 RequestId string `json:"RequestId" xml:"RequestId"`
}

func CreateModifyScalingRuleRequest() (request *ModifyScalingRuleRequest) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 request = &ModifyScalingRuleRequest{RpcRequest: &requests.RpcRequest{}}
 request.InitWithApiInfo("Ess", "2014-08-28", "ModifyScalingRule", "ess", "openAPI")
 return
}
func CreateModifyScalingRuleResponse() (response *ModifyScalingRuleResponse) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = &ModifyScalingRuleResponse{BaseResponse: &responses.BaseResponse{}}
 return
}
