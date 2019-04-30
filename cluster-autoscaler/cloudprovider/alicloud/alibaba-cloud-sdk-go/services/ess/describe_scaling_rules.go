package ess

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) DescribeScalingRules(request *DescribeScalingRulesRequest) (response *DescribeScalingRulesResponse, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = CreateDescribeScalingRulesResponse()
	err = client.DoAction(request, response)
	return
}
func (client *Client) DescribeScalingRulesWithChan(request *DescribeScalingRulesRequest) (<-chan *DescribeScalingRulesResponse, <-chan error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	responseChan := make(chan *DescribeScalingRulesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeScalingRules(request)
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
func (client *Client) DescribeScalingRulesWithCallback(request *DescribeScalingRulesRequest, callback func(response *DescribeScalingRulesResponse, err error)) <-chan int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeScalingRulesResponse
		var err error
		defer close(result)
		response, err = client.DescribeScalingRules(request)
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

type DescribeScalingRulesRequest struct {
	*requests.RpcRequest
	ScalingRuleName1	string			`position:"Query" name:"ScalingRuleName.1"`
	ResourceOwnerId		requests.Integer	`position:"Query" name:"ResourceOwnerId"`
	ScalingRuleName2	string			`position:"Query" name:"ScalingRuleName.2"`
	ScalingRuleName3	string			`position:"Query" name:"ScalingRuleName.3"`
	ScalingRuleName4	string			`position:"Query" name:"ScalingRuleName.4"`
	ScalingRuleName5	string			`position:"Query" name:"ScalingRuleName.5"`
	ScalingGroupId		string			`position:"Query" name:"ScalingGroupId"`
	ScalingRuleName6	string			`position:"Query" name:"ScalingRuleName.6"`
	ScalingRuleName7	string			`position:"Query" name:"ScalingRuleName.7"`
	ScalingRuleName8	string			`position:"Query" name:"ScalingRuleName.8"`
	ScalingRuleAri9		string			`position:"Query" name:"ScalingRuleAri.9"`
	ScalingRuleName9	string			`position:"Query" name:"ScalingRuleName.9"`
	PageNumber		requests.Integer	`position:"Query" name:"PageNumber"`
	PageSize		requests.Integer	`position:"Query" name:"PageSize"`
	ScalingRuleId10		string			`position:"Query" name:"ScalingRuleId.10"`
	ResourceOwnerAccount	string			`position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount		string			`position:"Query" name:"OwnerAccount"`
	OwnerId			requests.Integer	`position:"Query" name:"OwnerId"`
	ScalingRuleAri1		string			`position:"Query" name:"ScalingRuleAri.1"`
	ScalingRuleAri2		string			`position:"Query" name:"ScalingRuleAri.2"`
	ScalingRuleName10	string			`position:"Query" name:"ScalingRuleName.10"`
	ScalingRuleAri3		string			`position:"Query" name:"ScalingRuleAri.3"`
	ScalingRuleAri4		string			`position:"Query" name:"ScalingRuleAri.4"`
	ScalingRuleId8		string			`position:"Query" name:"ScalingRuleId.8"`
	ScalingRuleAri5		string			`position:"Query" name:"ScalingRuleAri.5"`
	ScalingRuleId9		string			`position:"Query" name:"ScalingRuleId.9"`
	ScalingRuleAri6		string			`position:"Query" name:"ScalingRuleAri.6"`
	ScalingRuleAri7		string			`position:"Query" name:"ScalingRuleAri.7"`
	ScalingRuleAri10	string			`position:"Query" name:"ScalingRuleAri.10"`
	ScalingRuleAri8		string			`position:"Query" name:"ScalingRuleAri.8"`
	ScalingRuleId4		string			`position:"Query" name:"ScalingRuleId.4"`
	ScalingRuleId5		string			`position:"Query" name:"ScalingRuleId.5"`
	ScalingRuleId6		string			`position:"Query" name:"ScalingRuleId.6"`
	ScalingRuleId7		string			`position:"Query" name:"ScalingRuleId.7"`
	ScalingRuleId1		string			`position:"Query" name:"ScalingRuleId.1"`
	ScalingRuleId2		string			`position:"Query" name:"ScalingRuleId.2"`
	ScalingRuleId3		string			`position:"Query" name:"ScalingRuleId.3"`
}
type DescribeScalingRulesResponse struct {
	*responses.BaseResponse
	TotalCount	int		`json:"TotalCount" xml:"TotalCount"`
	PageNumber	int		`json:"PageNumber" xml:"PageNumber"`
	PageSize	int		`json:"PageSize" xml:"PageSize"`
	RequestId	string		`json:"RequestId" xml:"RequestId"`
	ScalingRules	ScalingRules	`json:"ScalingRules" xml:"ScalingRules"`
}

func CreateDescribeScalingRulesRequest() (request *DescribeScalingRulesRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	request = &DescribeScalingRulesRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Ess", "2014-08-28", "DescribeScalingRules", "ess", "openAPI")
	return
}
func CreateDescribeScalingRulesResponse() (response *DescribeScalingRulesResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	response = &DescribeScalingRulesResponse{BaseResponse: &responses.BaseResponse{}}
	return
}

type ScalingRules struct {
	ScalingRule []ScalingRule `json:"ScalingRule" xml:"ScalingRule"`
}
type ScalingRule struct {
	ScalingRuleId	string	`json:"ScalingRuleId" xml:"ScalingRuleId"`
	ScalingGroupId	string	`json:"ScalingGroupId" xml:"ScalingGroupId"`
	ScalingRuleName	string	`json:"ScalingRuleName" xml:"ScalingRuleName"`
	Cooldown	int	`json:"Cooldown" xml:"Cooldown"`
	AdjustmentType	string	`json:"AdjustmentType" xml:"AdjustmentType"`
	AdjustmentValue	int	`json:"AdjustmentValue" xml:"AdjustmentValue"`
	MinSize		int	`json:"MinSize" xml:"MinSize"`
	MaxSize		int	`json:"MaxSize" xml:"MaxSize"`
	ScalingRuleAri	string	`json:"ScalingRuleAri" xml:"ScalingRuleAri"`
}
