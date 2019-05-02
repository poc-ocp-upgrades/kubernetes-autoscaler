package credentials

type StsRoleNameOnEcsCredential struct{ RoleName string }

func NewStsRoleNameOnEcsCredential(roleName string) *StsRoleNameOnEcsCredential {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &StsRoleNameOnEcsCredential{RoleName: roleName}
}
func (oldCred *StsRoleNameOnEcsCredential) ToEcsRamRoleCredential() *EcsRamRoleCredential {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &EcsRamRoleCredential{RoleName: oldCred.RoleName}
}

type EcsRamRoleCredential struct{ RoleName string }

func NewEcsRamRoleCredential(roleName string) *EcsRamRoleCredential {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &EcsRamRoleCredential{RoleName: roleName}
}
