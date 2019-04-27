package signers

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
	"time"
)

const defaultInAdvanceScale = 0.8

type credentialUpdater struct {
	credentialExpiration	int
	lastUpdateTimestamp	int64
	inAdvanceScale		float64
	buildRequestMethod	func() (*requests.CommonRequest, error)
	responseCallBack	func(response *responses.CommonResponse) error
	refreshApi		func(request *requests.CommonRequest) (response *responses.CommonResponse, err error)
}

func (updater *credentialUpdater) needUpdateCredential() (result bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if updater.inAdvanceScale == 0 {
		updater.inAdvanceScale = defaultInAdvanceScale
	}
	return time.Now().Unix()-updater.lastUpdateTimestamp >= int64(float64(updater.credentialExpiration)*updater.inAdvanceScale)
}
func (updater *credentialUpdater) updateCredential() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request, err := updater.buildRequestMethod()
	if err != nil {
		return
	}
	response, err := updater.refreshApi(request)
	if err != nil {
		return
	}
	updater.lastUpdateTimestamp = time.Now().Unix()
	err = updater.responseCallBack(response)
	return
}
