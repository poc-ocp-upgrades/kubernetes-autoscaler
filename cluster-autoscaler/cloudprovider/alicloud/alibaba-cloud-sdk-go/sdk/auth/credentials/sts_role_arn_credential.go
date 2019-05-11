package credentials

type StsRoleArnCredential struct {
	AccessKeyId				string
	AccessKeySecret			string
	RoleArn					string
	RoleSessionName			string
	RoleSessionExpiration	int
}
type RamRoleArnCredential struct {
	AccessKeyId				string
	AccessKeySecret			string
	RoleArn					string
	RoleSessionName			string
	RoleSessionExpiration	int
}

func NewStsRoleArnCredential(accessKeyId, accessKeySecret, roleArn, roleSessionName string, roleSessionExpiration int) *StsRoleArnCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &StsRoleArnCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret, RoleArn: roleArn, RoleSessionName: roleSessionName, RoleSessionExpiration: roleSessionExpiration}
}
func (oldCred *StsRoleArnCredential) ToRamRoleArnCredential() *RamRoleArnCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &RamRoleArnCredential{AccessKeyId: oldCred.AccessKeyId, AccessKeySecret: oldCred.AccessKeySecret, RoleArn: oldCred.RoleArn, RoleSessionName: oldCred.RoleSessionName, RoleSessionExpiration: oldCred.RoleSessionExpiration}
}
func NewRamRoleArnCredential(accessKeyId, accessKeySecret, roleArn, roleSessionName string, roleSessionExpiration int) *RamRoleArnCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &RamRoleArnCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret, RoleArn: roleArn, RoleSessionName: roleSessionName, RoleSessionExpiration: roleSessionExpiration}
}
