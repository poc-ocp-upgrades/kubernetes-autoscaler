package alicloud

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/metadata"
	"k8s.io/klog"
	"os"
)

const (
	accessKeyId		= "ACCESS_KEY_ID"
	accessKeyScret	= "ACCESS_KEY_SECRET"
	regionId		= "REGION_ID"
)

type cloudConfig struct {
	RegionId		string
	AccessKeyID		string
	AccessKeySecret	string
	STSEnabled		bool
}

func (cc *cloudConfig) isValid() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if cc.AccessKeyID == "" {
		cc.AccessKeyID = os.Getenv(accessKeyId)
	}
	if cc.AccessKeySecret == "" {
		cc.AccessKeySecret = os.Getenv(accessKeyScret)
	}
	if cc.RegionId == "" {
		cc.RegionId = os.Getenv(regionId)
	}
	if cc.RegionId == "" || cc.AccessKeyID == "" || cc.AccessKeySecret == "" {
		klog.V(5).Infof("Failed to get AccessKeyId:%s,AccessKeySecret:%s,RegionId:%s from cloudConfig and Env\n", cc.AccessKeyID, cc.AccessKeySecret, cc.RegionId)
		klog.V(5).Infof("Try to use sts token in metadata instead.\n")
		if cc.validateSTSToken() == true && cc.getRegion() != "" {
			cc.STSEnabled = true
			return true
		}
	} else {
		cc.STSEnabled = false
		return true
	}
	return false
}
func (cc *cloudConfig) validateSTSToken() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := metadata.NewMetaData(nil)
	r, err := m.RoleName()
	if err != nil || r == "" {
		klog.Warningf("The role name %s is not valid and error is %v", r, err)
		return false
	}
	return true
}
func (cc *cloudConfig) getSTSToken() (metadata.RoleAuth, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	m := metadata.NewMetaData(nil)
	r, err := m.RoleName()
	if err != nil {
		return metadata.RoleAuth{}, err
	}
	auth, err := m.RamRoleToken(r)
	if err != nil {
		return metadata.RoleAuth{}, err
	}
	return auth, nil
}
func (cc *cloudConfig) getRegion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if cc.RegionId != "" {
		return cc.RegionId
	}
	m := metadata.NewMetaData(nil)
	r, err := m.Region()
	if err != nil {
		klog.Errorf("Failed to get RegionId from metadata.Because of %s\n", err.Error())
	}
	return r
}
