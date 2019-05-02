package ess

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) DescribeScalingGroups(request *DescribeScalingGroupsRequest) (response *DescribeScalingGroupsResponse, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = CreateDescribeScalingGroupsResponse()
 err = client.DoAction(request, response)
 return
}
func (client *Client) DescribeScalingGroupsWithChan(request *DescribeScalingGroupsRequest) (<-chan *DescribeScalingGroupsResponse, <-chan error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 responseChan := make(chan *DescribeScalingGroupsResponse, 1)
 errChan := make(chan error, 1)
 err := client.AddAsyncTask(func() {
  defer close(responseChan)
  defer close(errChan)
  response, err := client.DescribeScalingGroups(request)
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
func (client *Client) DescribeScalingGroupsWithCallback(request *DescribeScalingGroupsRequest, callback func(response *DescribeScalingGroupsResponse, err error)) <-chan int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := make(chan int, 1)
 err := client.AddAsyncTask(func() {
  var response *DescribeScalingGroupsResponse
  var err error
  defer close(result)
  response, err = client.DescribeScalingGroups(request)
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

type DescribeScalingGroupsRequest struct {
 *requests.RpcRequest
 ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
 ScalingGroupId10     string           `position:"Query" name:"ScalingGroupId.10"`
 ScalingGroupId12     string           `position:"Query" name:"ScalingGroupId.12"`
 ScalingGroupId13     string           `position:"Query" name:"ScalingGroupId.13"`
 ScalingGroupId14     string           `position:"Query" name:"ScalingGroupId.14"`
 ScalingGroupId15     string           `position:"Query" name:"ScalingGroupId.15"`
 OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
 PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
 PageSize             requests.Integer `position:"Query" name:"PageSize"`
 ScalingGroupName20   string           `position:"Query" name:"ScalingGroupName.20"`
 ScalingGroupName19   string           `position:"Query" name:"ScalingGroupName.19"`
 ScalingGroupId20     string           `position:"Query" name:"ScalingGroupId.20"`
 ScalingGroupName18   string           `position:"Query" name:"ScalingGroupName.18"`
 ScalingGroupName17   string           `position:"Query" name:"ScalingGroupName.17"`
 ScalingGroupName16   string           `position:"Query" name:"ScalingGroupName.16"`
 ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
 ScalingGroupName     string           `position:"Query" name:"ScalingGroupName"`
 OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
 ScalingGroupName1    string           `position:"Query" name:"ScalingGroupName.1"`
 ScalingGroupName2    string           `position:"Query" name:"ScalingGroupName.2"`
 ScalingGroupId2      string           `position:"Query" name:"ScalingGroupId.2"`
 ScalingGroupId1      string           `position:"Query" name:"ScalingGroupId.1"`
 ScalingGroupId6      string           `position:"Query" name:"ScalingGroupId.6"`
 ScalingGroupId16     string           `position:"Query" name:"ScalingGroupId.16"`
 ScalingGroupName7    string           `position:"Query" name:"ScalingGroupName.7"`
 ScalingGroupName11   string           `position:"Query" name:"ScalingGroupName.11"`
 ScalingGroupId5      string           `position:"Query" name:"ScalingGroupId.5"`
 ScalingGroupId17     string           `position:"Query" name:"ScalingGroupId.17"`
 ScalingGroupName8    string           `position:"Query" name:"ScalingGroupName.8"`
 ScalingGroupName10   string           `position:"Query" name:"ScalingGroupName.10"`
 ScalingGroupId4      string           `position:"Query" name:"ScalingGroupId.4"`
 ScalingGroupId18     string           `position:"Query" name:"ScalingGroupId.18"`
 ScalingGroupName9    string           `position:"Query" name:"ScalingGroupName.9"`
 ScalingGroupId3      string           `position:"Query" name:"ScalingGroupId.3"`
 ScalingGroupId19     string           `position:"Query" name:"ScalingGroupId.19"`
 ScalingGroupName3    string           `position:"Query" name:"ScalingGroupName.3"`
 ScalingGroupName15   string           `position:"Query" name:"ScalingGroupName.15"`
 ScalingGroupId9      string           `position:"Query" name:"ScalingGroupId.9"`
 ScalingGroupName4    string           `position:"Query" name:"ScalingGroupName.4"`
 ScalingGroupName14   string           `position:"Query" name:"ScalingGroupName.14"`
 ScalingGroupId8      string           `position:"Query" name:"ScalingGroupId.8"`
 ScalingGroupName5    string           `position:"Query" name:"ScalingGroupName.5"`
 ScalingGroupName13   string           `position:"Query" name:"ScalingGroupName.13"`
 ScalingGroupId7      string           `position:"Query" name:"ScalingGroupId.7"`
 ScalingGroupName6    string           `position:"Query" name:"ScalingGroupName.6"`
 ScalingGroupName12   string           `position:"Query" name:"ScalingGroupName.12"`
}
type DescribeScalingGroupsResponse struct {
 *responses.BaseResponse
 TotalCount    int           `json:"TotalCount" xml:"TotalCount"`
 PageNumber    int           `json:"PageNumber" xml:"PageNumber"`
 PageSize      int           `json:"PageSize" xml:"PageSize"`
 RequestId     string        `json:"RequestId" xml:"RequestId"`
 ScalingGroups ScalingGroups `json:"ScalingGroups" xml:"ScalingGroups"`
}

func CreateDescribeScalingGroupsRequest() (request *DescribeScalingGroupsRequest) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 request = &DescribeScalingGroupsRequest{RpcRequest: &requests.RpcRequest{}}
 request.InitWithApiInfo("Ess", "2014-08-28", "DescribeScalingGroups", "ess", "openAPI")
 return
}
func CreateDescribeScalingGroupsResponse() (response *DescribeScalingGroupsResponse) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = &DescribeScalingGroupsResponse{BaseResponse: &responses.BaseResponse{}}
 return
}

type ScalingGroups struct {
 ScalingGroup []ScalingGroup `json:"ScalingGroup" xml:"ScalingGroup"`
}
type ScalingGroup struct {
 DefaultCooldown              int             `json:"DefaultCooldown" xml:"DefaultCooldown"`
 MaxSize                      int             `json:"MaxSize" xml:"MaxSize"`
 PendingWaitCapacity          int             `json:"PendingWaitCapacity" xml:"PendingWaitCapacity"`
 RemovingWaitCapacity         int             `json:"RemovingWaitCapacity" xml:"RemovingWaitCapacity"`
 PendingCapacity              int             `json:"PendingCapacity" xml:"PendingCapacity"`
 RemovingCapacity             int             `json:"RemovingCapacity" xml:"RemovingCapacity"`
 ScalingGroupName             string          `json:"ScalingGroupName" xml:"ScalingGroupName"`
 ActiveCapacity               int             `json:"ActiveCapacity" xml:"ActiveCapacity"`
 StandbyCapacity              int             `json:"StandbyCapacity" xml:"StandbyCapacity"`
 ProtectedCapacity            int             `json:"ProtectedCapacity" xml:"ProtectedCapacity"`
 ActiveScalingConfigurationId string          `json:"ActiveScalingConfigurationId" xml:"ActiveScalingConfigurationId"`
 LaunchTemplateId             string          `json:"LaunchTemplateId" xml:"LaunchTemplateId"`
 LaunchTemplateVersion        string          `json:"LaunchTemplateVersion" xml:"LaunchTemplateVersion"`
 ScalingGroupId               string          `json:"ScalingGroupId" xml:"ScalingGroupId"`
 RegionId                     string          `json:"RegionId" xml:"RegionId"`
 TotalCapacity                int             `json:"TotalCapacity" xml:"TotalCapacity"`
 MinSize                      int             `json:"MinSize" xml:"MinSize"`
 LifecycleState               string          `json:"LifecycleState" xml:"LifecycleState"`
 CreationTime                 string          `json:"CreationTime" xml:"CreationTime"`
 ModificationTime             string          `json:"ModificationTime" xml:"ModificationTime"`
 VpcId                        string          `json:"VpcId" xml:"VpcId"`
 VSwitchId                    string          `json:"VSwitchId" xml:"VSwitchId"`
 MultiAZPolicy                string          `json:"MultiAZPolicy" xml:"MultiAZPolicy"`
 HealthCheckType              string          `json:"HealthCheckType" xml:"HealthCheckType"`
 VSwitchIds                   VSwitchIds      `json:"VSwitchIds" xml:"VSwitchIds"`
 RemovalPolicies              RemovalPolicies `json:"RemovalPolicies" xml:"RemovalPolicies"`
 DBInstanceIds                DBInstanceIds   `json:"DBInstanceIds" xml:"DBInstanceIds"`
 LoadBalancerIds              LoadBalancerIds `json:"LoadBalancerIds" xml:"LoadBalancerIds"`
}
type VSwitchIds struct {
 VSwitchId []string `json:"VSwitchId" xml:"VSwitchId"`
}
type RemovalPolicies struct {
 RemovalPolicy []string `json:"RemovalPolicy" xml:"RemovalPolicy"`
}
type DBInstanceIds struct {
 DBInstanceId []string `json:"DBInstanceId" xml:"DBInstanceId"`
}
type LoadBalancerIds struct {
 LoadBalancerId []string `json:"LoadBalancerId" xml:"LoadBalancerId"`
}
