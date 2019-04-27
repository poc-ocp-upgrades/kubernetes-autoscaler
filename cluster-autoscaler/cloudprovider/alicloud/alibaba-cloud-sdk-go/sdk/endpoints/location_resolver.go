package endpoints

import (
	"encoding/json"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"sync"
	"time"
)

const (
	EndpointCacheExpireTime = 3600
)

var lastClearTimePerProduct = struct {
	sync.RWMutex
	cache	map[string]int64
}{cache: make(map[string]int64)}
var endpointCache = struct {
	sync.RWMutex
	cache	map[string]string
}{cache: make(map[string]string)}

type LocationResolver struct{}

func (resolver *LocationResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(param.LocationProduct) <= 0 {
		support = false
		return
	}
	cacheKey := param.Product + "#" + param.RegionId
	if endpointCache.cache != nil && len(endpointCache.cache[cacheKey]) > 0 && !CheckCacheIsExpire(cacheKey) {
		endpoint = endpointCache.cache[cacheKey]
		support = true
		return
	}
	getEndpointRequest := requests.NewCommonRequest()
	getEndpointRequest.Product = "Location"
	getEndpointRequest.Version = "2015-06-12"
	getEndpointRequest.ApiName = "DescribeEndpoints"
	getEndpointRequest.Domain = "location.aliyuncs.com"
	getEndpointRequest.Method = "GET"
	getEndpointRequest.Scheme = requests.HTTPS
	getEndpointRequest.QueryParams["Id"] = param.RegionId
	getEndpointRequest.QueryParams["ServiceCode"] = param.LocationProduct
	if len(param.LocationEndpointType) > 0 {
		getEndpointRequest.QueryParams["Type"] = param.LocationEndpointType
	} else {
		getEndpointRequest.QueryParams["Type"] = "openAPI"
	}
	response, err := param.CommonApi(getEndpointRequest)
	var getEndpointResponse GetEndpointResponse
	if !response.IsSuccess() {
		support = false
		return
	}
	json.Unmarshal([]byte(response.GetHttpContentString()), &getEndpointResponse)
	if !getEndpointResponse.Success || getEndpointResponse.Endpoints == nil {
		support = false
		return
	}
	if len(getEndpointResponse.Endpoints.Endpoint) <= 0 {
		support = false
		return
	}
	if len(getEndpointResponse.Endpoints.Endpoint[0].Endpoint) > 0 {
		endpoint = getEndpointResponse.Endpoints.Endpoint[0].Endpoint
		endpointCache.Lock()
		endpointCache.cache[cacheKey] = endpoint
		endpointCache.Unlock()
		lastClearTimePerProduct.Lock()
		lastClearTimePerProduct.cache[cacheKey] = time.Now().Unix()
		lastClearTimePerProduct.Unlock()
		support = true
		return
	}
	support = false
	return
}
func CheckCacheIsExpire(cacheKey string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lastClearTime := lastClearTimePerProduct.cache[cacheKey]
	if lastClearTime <= 0 {
		lastClearTime = time.Now().Unix()
		lastClearTimePerProduct.Lock()
		lastClearTimePerProduct.cache[cacheKey] = lastClearTime
		lastClearTimePerProduct.Unlock()
	}
	now := time.Now().Unix()
	elapsedTime := now - lastClearTime
	if elapsedTime > EndpointCacheExpireTime {
		return true
	}
	return false
}

type GetEndpointResponse struct {
	Endpoints	*EndpointsObj
	RequestId	string
	Success		bool
}
type EndpointsObj struct{ Endpoint []EndpointObj }
type EndpointObj struct {
	Protocols	map[string]string
	Type		string
	Namespace	string
	Id		string
	SerivceCode	string
	Endpoint	string
}
