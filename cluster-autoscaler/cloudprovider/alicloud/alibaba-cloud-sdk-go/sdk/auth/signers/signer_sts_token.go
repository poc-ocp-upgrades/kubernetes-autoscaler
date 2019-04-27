package signers

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

type StsTokenSigner struct {
	credential *credentials.StsTokenCredential
}

func NewStsTokenSigner(credential *credentials.StsTokenCredential) (*StsTokenSigner, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &StsTokenSigner{credential: credential}, nil
}
func (*StsTokenSigner) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "HMAC-SHA1"
}
func (*StsTokenSigner) GetType() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (*StsTokenSigner) GetVersion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "1.0"
}
func (signer *StsTokenSigner) GetAccessKeyId() (accessKeyId string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return signer.credential.AccessKeyId, nil
}
func (signer *StsTokenSigner) GetExtraParam() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return map[string]string{"SecurityToken": signer.credential.AccessKeyStsToken}
}
func (signer *StsTokenSigner) Sign(stringToSign, secretSuffix string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret := signer.credential.AccessKeySecret + secretSuffix
	return ShaHmac1(stringToSign, secret)
}
func (signer *StsTokenSigner) Shutdown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
}
