package endpoints

import (
	"encoding/json"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/errors"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/responses"
	"sync"
)

const (
	ResolveEndpointUserGuideLink = ""
)

var once sync.Once
var resolvers []Resolver

type Resolver interface {
	TryResolve(param *ResolveParam) (endpoint string, support bool, err error)
}

func Resolve(param *ResolveParam) (endpoint string, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	supportedResolvers := getAllResolvers()
	for _, resolver := range supportedResolvers {
		endpoint, supported, err := resolver.TryResolve(param)
		if supported {
			return endpoint, err
		}
	}
	errorMsg := fmt.Sprintf(errors.CanNotResolveEndpointErrorMessage, param, ResolveEndpointUserGuideLink)
	err = errors.NewClientError(errors.CanNotResolveEndpointErrorCode, errorMsg, nil)
	return
}
func getAllResolvers() []Resolver {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	once.Do(func() {
		resolvers = []Resolver{&SimpleHostResolver{}, &MappingResolver{}, &LocationResolver{}, &LocalRegionalResolver{}, &LocalGlobalResolver{}}
	})
	return resolvers
}

type ResolveParam struct {
	Domain			string
	Product			string
	RegionId		string
	LocationProduct		string
	LocationEndpointType	string
	CommonApi		func(request *requests.CommonRequest) (response *responses.CommonResponse, err error)	`json:"-"`
}

func (param *ResolveParam) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	jsonBytes, err := json.Marshal(param)
	if err != nil {
		return fmt.Sprint("ResolveParam.String() process error:", err)
	}
	return string(jsonBytes)
}
