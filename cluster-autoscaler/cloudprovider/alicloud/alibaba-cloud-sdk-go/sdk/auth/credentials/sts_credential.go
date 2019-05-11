package credentials

type StsTokenCredential struct {
	AccessKeyId			string
	AccessKeySecret		string
	AccessKeyStsToken	string
}

func NewStsTokenCredential(accessKeyId, accessKeySecret, accessKeyStsToken string) *StsTokenCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &StsTokenCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret, AccessKeyStsToken: accessKeyStsToken}
}
