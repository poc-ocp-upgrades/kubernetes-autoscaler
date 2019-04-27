package endpoints

import (
	"fmt"
	"strings"
)

const keyFormatter = "%s::%s"

var endpointMapping = make(map[string]string)

func AddEndpointMapping(regionId, productId, endpoint string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key := fmt.Sprintf(keyFormatter, strings.ToLower(regionId), strings.ToLower(productId))
	endpointMapping[key] = endpoint
	return nil
}

type MappingResolver struct{}

func (resolver *MappingResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	key := fmt.Sprintf(keyFormatter, strings.ToLower(param.RegionId), strings.ToLower(param.Product))
	endpoint, contains := endpointMapping[key]
	return endpoint, contains, nil
}
