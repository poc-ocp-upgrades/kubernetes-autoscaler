package signers

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

type SignerV2 struct {
	credential *credentials.RsaKeyPairCredential
}

func (signer *SignerV2) GetExtraParam() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func NewSignerV2(credential *credentials.RsaKeyPairCredential) (*SignerV2, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &SignerV2{credential: credential}, nil
}
func (*SignerV2) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "SHA256withRSA"
}
func (*SignerV2) GetType() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "PRIVATEKEY"
}
func (*SignerV2) GetVersion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "1.0"
}
func (signer *SignerV2) GetAccessKeyId() (accessKeyId string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return signer.credential.PublicKeyId, err
}
func (signer *SignerV2) Sign(stringToSign, secretSuffix string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret := signer.credential.PrivateKey
	return Sha256WithRsa(stringToSign, secret)
}
func (signer *SignerV2) Shutdown() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
