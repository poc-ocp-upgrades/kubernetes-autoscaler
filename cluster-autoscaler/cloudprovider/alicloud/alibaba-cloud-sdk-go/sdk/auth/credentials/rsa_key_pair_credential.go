package credentials

type RsaKeyPairCredential struct {
	PrivateKey		string
	PublicKeyId		string
	SessionExpiration	int
}

func NewRsaKeyPairCredential(privateKey, publicKeyId string, sessionExpiration int) *RsaKeyPairCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &RsaKeyPairCredential{PrivateKey: privateKey, PublicKeyId: publicKeyId, SessionExpiration: sessionExpiration}
}
