package requests

import (
	"bytes"
	"io"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/utils"
	"net/url"
	"sort"
	"strings"
)

type RoaRequest struct {
	*baseRequest
	pathPattern	string
	PathParams	map[string]string
}

func (*RoaRequest) GetStyle() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ROA
}
func (request *RoaRequest) GetBodyReader() io.Reader {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if request.FormParams != nil && len(request.FormParams) > 0 {
		formString := utils.GetUrlFormedMap(request.FormParams)
		return strings.NewReader(formString)
	} else if len(request.Content) > 0 {
		return bytes.NewReader(request.Content)
	} else {
		return nil
	}
}
func (request *RoaRequest) GetQueries() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return request.queries
}
func (request *RoaRequest) BuildQueries() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return request.buildQueries(false)
}
func (request *RoaRequest) buildQueries(needParamEncode bool) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	path := request.pathPattern
	for key, value := range request.PathParams {
		path = strings.Replace(path, "["+key+"]", value, 1)
	}
	queryParams := request.QueryParams
	splitArray := strings.Split(path, "?")
	path = splitArray[0]
	if len(splitArray) > 1 && len(splitArray[1]) > 0 {
		queryParams[splitArray[1]] = ""
	}
	var queryKeys []string
	for key := range queryParams {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)
	urlBuilder := bytes.Buffer{}
	urlBuilder.WriteString(path)
	if len(queryKeys) > 0 {
		urlBuilder.WriteString("?")
	}
	for i := 0; i < len(queryKeys); i++ {
		queryKey := queryKeys[i]
		urlBuilder.WriteString(queryKey)
		if value := queryParams[queryKey]; len(value) > 0 {
			urlBuilder.WriteString("=")
			if needParamEncode {
				urlBuilder.WriteString(url.QueryEscape(value))
			} else {
				urlBuilder.WriteString(value)
			}
		}
		if i < len(queryKeys)-1 {
			urlBuilder.WriteString("&")
		}
	}
	result := urlBuilder.String()
	result = popStandardUrlencode(result)
	request.queries = result
	return request.queries
}
func popStandardUrlencode(stringToSign string) (result string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	result = strings.Replace(stringToSign, "+", "%20", -1)
	result = strings.Replace(result, "*", "%2A", -1)
	result = strings.Replace(result, "%7E", "~", -1)
	return
}
func (request *RoaRequest) GetUrl() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.ToLower(request.Scheme) + "://" + request.Domain + ":" + request.Port + request.GetQueries()
}
func (request *RoaRequest) BuildUrl() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strings.ToLower(request.Scheme) + "://" + request.Domain + ":" + request.Port + request.buildQueries(true)
}
func (request *RoaRequest) addPathParam(key, value string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.PathParams[key] = value
}
func (request *RoaRequest) InitWithApiInfo(product, version, action, uriPattern, serviceCode, endpointType string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.baseRequest = defaultBaseRequest()
	request.PathParams = make(map[string]string)
	request.Headers["x-acs-version"] = version
	request.pathPattern = uriPattern
	request.locationServiceCode = serviceCode
	request.locationEndpointType = endpointType
}
func (request *RoaRequest) initWithCommonRequest(commonRequest *CommonRequest) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	request.baseRequest = commonRequest.baseRequest
	request.PathParams = commonRequest.PathParams
	request.Headers["x-acs-version"] = commonRequest.Version
	request.pathPattern = commonRequest.PathPattern
	request.locationServiceCode = ""
	request.locationEndpointType = ""
}
