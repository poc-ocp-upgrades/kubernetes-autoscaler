package alicloud

import (
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/services/ecs"
	"k8s.io/klog"
	"time"
)

type ecsInstance interface {
	DescribeInstanceTypes(req *ecs.DescribeInstanceTypesRequest) (*ecs.DescribeInstanceTypesResponse, error)
}
type instanceType struct {
	instanceTypeID	string
	vcpu			int64
	memoryInBytes	int64
	gpu				int64
}
type instanceTypeModel struct{ instanceType }
type instanceWrapper struct {
	ecsInstance
	InstanceTypeCache	map[string]*instanceTypeModel
}

func (iw *instanceWrapper) getInstanceTypeById(typeId string) (*instanceType, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if instanceTypeModel := iw.FindInstanceType(typeId); instanceTypeModel != nil {
		return &instanceTypeModel.instanceType, nil
	}
	err := iw.RefreshCache()
	if err != nil {
		klog.Errorf("failed to refresh instance type cache,because of %s", err.Error())
		return nil, err
	}
	if instanceTypeModel := iw.FindInstanceType(typeId); instanceTypeModel != nil {
		return &instanceTypeModel.instanceType, nil
	}
	return nil, fmt.Errorf("failed to find the specific instance type by Id: %s", typeId)
}
func (iw *instanceWrapper) FindInstanceType(typeId string) *instanceTypeModel {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if iw.InstanceTypeCache == nil || iw.InstanceTypeCache[typeId] == nil {
		return nil
	}
	return iw.InstanceTypeCache[typeId]
}
func (iw *instanceWrapper) RefreshCache() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req := ecs.CreateDescribeInstanceTypesRequest()
	resp, err := iw.DescribeInstanceTypes(req)
	if err != nil {
		return err
	}
	if iw.InstanceTypeCache == nil {
		iw.InstanceTypeCache = make(map[string]*instanceTypeModel)
	}
	types := resp.InstanceTypes.InstanceType
	for _, item := range types {
		iw.InstanceTypeCache[item.InstanceTypeId] = &instanceTypeModel{instanceType{instanceTypeID: item.InstanceTypeId, vcpu: int64(item.CpuCoreCount), memoryInBytes: int64(item.MemorySize * 1024 * 1024 * 1024), gpu: int64(item.GPUAmount)}}
	}
	return nil
}
func newInstanceWrapper(cfg *cloudConfig) (*instanceWrapper, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if cfg.isValid() == false {
		return nil, fmt.Errorf("your cloud config is not valid")
	}
	iw := &instanceWrapper{}
	if cfg.STSEnabled == true {
		go func(iw *instanceWrapper, cfg *cloudConfig) {
			timer := time.NewTicker(refreshClientInterval)
			defer timer.Stop()
			for {
				select {
				case <-timer.C:
					client, err := getEcsClient(cfg)
					if err == nil {
						iw.ecsInstance = client
					}
				}
			}
		}(iw, cfg)
	}
	client, err := getEcsClient(cfg)
	if err == nil {
		iw.ecsInstance = client
	}
	return iw, err
}
func getEcsClient(cfg *cloudConfig) (client *ecs.Client, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	region := cfg.getRegion()
	if cfg.STSEnabled == true {
		auth, err := cfg.getSTSToken()
		if err != nil {
			klog.Errorf("failed to get sts token from metadata,because of %s", err.Error())
			return nil, err
		}
		client, err = ecs.NewClientWithStsToken(region, auth.AccessKeyId, auth.AccessKeySecret, auth.SecurityToken)
		if err != nil {
			klog.Errorf("failed to create client with sts in metadata,because of %s", err.Error())
		}
	} else {
		client, err = ecs.NewClientWithAccessKey(region, cfg.AccessKeyID, cfg.AccessKeySecret)
		if err != nil {
			klog.Errorf("failed to create ecs client with AccessKeyId and AccessKeySecret,because of %s", err.Error())
		}
	}
	return
}
