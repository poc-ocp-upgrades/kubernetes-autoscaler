package requests

import (
	"io"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/utils"
	"strings"
)

type RpcRequest struct{ *baseRequest }

func (request *RpcRequest) init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.baseRequest = defaultBaseRequest()
	request.Method = POST
}
func (*RpcRequest) GetStyle() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return RPC
}
func (request *RpcRequest) GetBodyReader() io.Reader {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if request.FormParams != nil && len(request.FormParams) > 0 {
		formString := utils.GetUrlFormedMap(request.FormParams)
		return strings.NewReader(formString)
	}
	return strings.NewReader("")
}
func (request *RpcRequest) BuildQueries() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.queries = "/?" + utils.GetUrlFormedMap(request.QueryParams)
	return request.queries
}
func (request *RpcRequest) GetQueries() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return request.queries
}
func (request *RpcRequest) BuildUrl() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.ToLower(request.Scheme) + "://" + request.Domain + ":" + request.Port + request.BuildQueries()
}
func (request *RpcRequest) GetUrl() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.ToLower(request.Scheme) + "://" + request.Domain + request.GetQueries()
}
func (request *RpcRequest) GetVersion() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return request.version
}
func (request *RpcRequest) GetActionName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return request.actionName
}
func (request *RpcRequest) addPathParam(key, value string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	panic("not support")
}
func (request *RpcRequest) InitWithApiInfo(product, version, action, serviceCode, endpointType string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.init()
	request.product = product
	request.version = version
	request.actionName = action
	request.locationServiceCode = serviceCode
	request.locationEndpointType = endpointType
}
