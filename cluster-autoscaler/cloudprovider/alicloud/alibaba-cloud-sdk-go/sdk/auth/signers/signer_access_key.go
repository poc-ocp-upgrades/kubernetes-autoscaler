package signers

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

type AccessKeySigner struct {
	credential *credentials.AccessKeyCredential
}

func (signer *AccessKeySigner) GetExtraParam() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func NewAccessKeySigner(credential *credentials.AccessKeyCredential) (*AccessKeySigner, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &AccessKeySigner{credential: credential}, nil
}
func (*AccessKeySigner) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "HMAC-SHA1"
}
func (*AccessKeySigner) GetType() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (*AccessKeySigner) GetVersion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "1.0"
}
func (signer *AccessKeySigner) GetAccessKeyId() (accessKeyId string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return signer.credential.AccessKeyId, nil
}
func (signer *AccessKeySigner) Sign(stringToSign, secretSuffix string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret := signer.credential.AccessKeySecret + secretSuffix
	return ShaHmac1(stringToSign, secret)
}
func (signer *AccessKeySigner) Shutdown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
