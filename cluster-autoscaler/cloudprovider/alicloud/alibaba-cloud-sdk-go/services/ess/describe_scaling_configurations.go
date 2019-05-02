package ess

import (
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
)

func (client *Client) DescribeScalingConfigurations(request *DescribeScalingConfigurationsRequest) (response *DescribeScalingConfigurationsResponse, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = CreateDescribeScalingConfigurationsResponse()
 err = client.DoAction(request, response)
 return
}
func (client *Client) DescribeScalingConfigurationsWithChan(request *DescribeScalingConfigurationsRequest) (<-chan *DescribeScalingConfigurationsResponse, <-chan error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 responseChan := make(chan *DescribeScalingConfigurationsResponse, 1)
 errChan := make(chan error, 1)
 err := client.AddAsyncTask(func() {
  defer close(responseChan)
  defer close(errChan)
  response, err := client.DescribeScalingConfigurations(request)
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
func (client *Client) DescribeScalingConfigurationsWithCallback(request *DescribeScalingConfigurationsRequest, callback func(response *DescribeScalingConfigurationsResponse, err error)) <-chan int {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := make(chan int, 1)
 err := client.AddAsyncTask(func() {
  var response *DescribeScalingConfigurationsResponse
  var err error
  defer close(result)
  response, err = client.DescribeScalingConfigurations(request)
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

type DescribeScalingConfigurationsRequest struct {
 *requests.RpcRequest
 ScalingConfigurationId6    string           `position:"Query" name:"ScalingConfigurationId.6"`
 ScalingConfigurationId7    string           `position:"Query" name:"ScalingConfigurationId.7"`
 ResourceOwnerId            requests.Integer `position:"Query" name:"ResourceOwnerId"`
 ScalingConfigurationId4    string           `position:"Query" name:"ScalingConfigurationId.4"`
 ScalingConfigurationId5    string           `position:"Query" name:"ScalingConfigurationId.5"`
 ScalingGroupId             string           `position:"Query" name:"ScalingGroupId"`
 ScalingConfigurationId8    string           `position:"Query" name:"ScalingConfigurationId.8"`
 ScalingConfigurationId9    string           `position:"Query" name:"ScalingConfigurationId.9"`
 ScalingConfigurationId10   string           `position:"Query" name:"ScalingConfigurationId.10"`
 PageNumber                 requests.Integer `position:"Query" name:"PageNumber"`
 ScalingConfigurationName2  string           `position:"Query" name:"ScalingConfigurationName.2"`
 ScalingConfigurationName3  string           `position:"Query" name:"ScalingConfigurationName.3"`
 ScalingConfigurationName1  string           `position:"Query" name:"ScalingConfigurationName.1"`
 PageSize                   requests.Integer `position:"Query" name:"PageSize"`
 ScalingConfigurationId2    string           `position:"Query" name:"ScalingConfigurationId.2"`
 ScalingConfigurationId3    string           `position:"Query" name:"ScalingConfigurationId.3"`
 ScalingConfigurationId1    string           `position:"Query" name:"ScalingConfigurationId.1"`
 ResourceOwnerAccount       string           `position:"Query" name:"ResourceOwnerAccount"`
 OwnerAccount               string           `position:"Query" name:"OwnerAccount"`
 ScalingConfigurationName6  string           `position:"Query" name:"ScalingConfigurationName.6"`
 ScalingConfigurationName7  string           `position:"Query" name:"ScalingConfigurationName.7"`
 ScalingConfigurationName4  string           `position:"Query" name:"ScalingConfigurationName.4"`
 ScalingConfigurationName5  string           `position:"Query" name:"ScalingConfigurationName.5"`
 OwnerId                    requests.Integer `position:"Query" name:"OwnerId"`
 ScalingConfigurationName8  string           `position:"Query" name:"ScalingConfigurationName.8"`
 ScalingConfigurationName9  string           `position:"Query" name:"ScalingConfigurationName.9"`
 ScalingConfigurationName10 string           `position:"Query" name:"ScalingConfigurationName.10"`
}
type DescribeScalingConfigurationsResponse struct {
 *responses.BaseResponse
 TotalCount            int                   `json:"TotalCount" xml:"TotalCount"`
 PageNumber            int                   `json:"PageNumber" xml:"PageNumber"`
 PageSize              int                   `json:"PageSize" xml:"PageSize"`
 RequestId             string                `json:"RequestId" xml:"RequestId"`
 ScalingConfigurations ScalingConfigurations `json:"ScalingConfigurations" xml:"ScalingConfigurations"`
}

func CreateDescribeScalingConfigurationsRequest() (request *DescribeScalingConfigurationsRequest) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 request = &DescribeScalingConfigurationsRequest{RpcRequest: &requests.RpcRequest{}}
 request.InitWithApiInfo("Ess", "2014-08-28", "DescribeScalingConfigurations", "ess", "openAPI")
 return
}
func CreateDescribeScalingConfigurationsResponse() (response *DescribeScalingConfigurationsResponse) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 response = &DescribeScalingConfigurationsResponse{BaseResponse: &responses.BaseResponse{}}
 return
}

type ScalingConfigurations struct {
 ScalingConfiguration []ScalingConfiguration `json:"ScalingConfiguration" xml:"ScalingConfiguration"`
}
type ScalingConfiguration struct {
 ScalingConfigurationId      string         `json:"ScalingConfigurationId" xml:"ScalingConfigurationId"`
 ScalingConfigurationName    string         `json:"ScalingConfigurationName" xml:"ScalingConfigurationName"`
 ScalingGroupId              string         `json:"ScalingGroupId" xml:"ScalingGroupId"`
 InstanceName                string         `json:"InstanceName" xml:"InstanceName"`
 ImageId                     string         `json:"ImageId" xml:"ImageId"`
 ImageName                   string         `json:"ImageName" xml:"ImageName"`
 HostName                    string         `json:"HostName" xml:"HostName"`
 InstanceType                string         `json:"InstanceType" xml:"InstanceType"`
 InstanceGeneration          string         `json:"InstanceGeneration" xml:"InstanceGeneration"`
 SecurityGroupId             string         `json:"SecurityGroupId" xml:"SecurityGroupId"`
 IoOptimized                 string         `json:"IoOptimized" xml:"IoOptimized"`
 InternetChargeType          string         `json:"InternetChargeType" xml:"InternetChargeType"`
 InternetMaxBandwidthIn      int            `json:"InternetMaxBandwidthIn" xml:"InternetMaxBandwidthIn"`
 InternetMaxBandwidthOut     int            `json:"InternetMaxBandwidthOut" xml:"InternetMaxBandwidthOut"`
 SystemDiskCategory          string         `json:"SystemDiskCategory" xml:"SystemDiskCategory"`
 SystemDiskSize              int            `json:"SystemDiskSize" xml:"SystemDiskSize"`
 LifecycleState              string         `json:"LifecycleState" xml:"LifecycleState"`
 CreationTime                string         `json:"CreationTime" xml:"CreationTime"`
 LoadBalancerWeight          int            `json:"LoadBalancerWeight" xml:"LoadBalancerWeight"`
 UserData                    string         `json:"UserData" xml:"UserData"`
 KeyPairName                 string         `json:"KeyPairName" xml:"KeyPairName"`
 RamRoleName                 string         `json:"RamRoleName" xml:"RamRoleName"`
 DeploymentSetId             string         `json:"DeploymentSetId" xml:"DeploymentSetId"`
 SecurityEnhancementStrategy string         `json:"SecurityEnhancementStrategy" xml:"SecurityEnhancementStrategy"`
 SpotStrategy                string         `json:"SpotStrategy" xml:"SpotStrategy"`
 PasswordInherit             bool           `json:"PasswordInherit" xml:"PasswordInherit"`
 InstanceTypes               InstanceTypes  `json:"InstanceTypes" xml:"InstanceTypes"`
 DataDisks                   DataDisks      `json:"DataDisks" xml:"DataDisks"`
 Tags                        Tags           `json:"Tags" xml:"Tags"`
 SpotPriceLimit              SpotPriceLimit `json:"SpotPriceLimit" xml:"SpotPriceLimit"`
}
type InstanceTypes struct {
 InstanceType []string `json:"InstanceType" xml:"InstanceType"`
}
type DataDisks struct {
 DataDisk []DataDisk `json:"DataDisk" xml:"DataDisk"`
}
type DataDisk struct {
 Size               int    `json:"Size" xml:"Size"`
 Category           string `json:"Category" xml:"Category"`
 SnapshotId         string `json:"SnapshotId" xml:"SnapshotId"`
 Device             string `json:"Device" xml:"Device"`
 DeleteWithInstance bool   `json:"DeleteWithInstance" xml:"DeleteWithInstance"`
}
type Tags struct {
 Tag []Tag `json:"Tag" xml:"Tag"`
}
type Tag struct {
 Key   string `json:"Key" xml:"Key"`
 Value string `json:"Value" xml:"Value"`
}
type SpotPriceLimit struct {
 SpotPriceModel []SpotPriceModel `json:"SpotPriceModel" xml:"SpotPriceModel"`
}
type SpotPriceModel struct {
 InstanceType string  `json:"InstanceType" xml:"InstanceType"`
 PriceLimit   float64 `json:"PriceLimit" xml:"PriceLimit"`
}
