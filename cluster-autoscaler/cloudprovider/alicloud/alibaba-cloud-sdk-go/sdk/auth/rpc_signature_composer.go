package auth

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/utils"
	"net/url"
	"sort"
	"strings"
)

func signRpcRequest(request requests.AcsRequest, signer Signer, regionId string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = completeRpcSignParams(request, signer, regionId)
	if err != nil {
		return
	}
	if _, containsSign := request.GetQueryParams()["Signature"]; containsSign {
		delete(request.GetQueryParams(), "Signature")
	}
	stringToSign := buildRpcStringToSign(request)
	request.SetStringToSign(stringToSign)
	signature := signer.Sign(stringToSign, "&")
	request.GetQueryParams()["Signature"] = signature
	return
}
func completeRpcSignParams(request requests.AcsRequest, signer Signer, regionId string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	queryParams := request.GetQueryParams()
	queryParams["Version"] = request.GetVersion()
	queryParams["Action"] = request.GetActionName()
	queryParams["Format"] = request.GetAcceptFormat()
	queryParams["Timestamp"] = utils.GetTimeInFormatISO8601()
	queryParams["SignatureMethod"] = signer.GetName()
	queryParams["SignatureType"] = signer.GetType()
	queryParams["SignatureVersion"] = signer.GetVersion()
	queryParams["SignatureNonce"] = utils.GetUUIDV4()
	queryParams["AccessKeyId"], err = signer.GetAccessKeyId()
	if err != nil {
		return
	}
	if _, contains := queryParams["RegionId"]; !contains {
		queryParams["RegionId"] = regionId
	}
	if extraParam := signer.GetExtraParam(); extraParam != nil {
		for key, value := range extraParam {
			queryParams[key] = value
		}
	}
	request.GetHeaders()["Content-Type"] = requests.Form
	formString := utils.GetUrlFormedMap(request.GetFormParams())
	request.SetContent([]byte(formString))
	return
}
func buildRpcStringToSign(request requests.AcsRequest) (stringToSign string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	signParams := make(map[string]string)
	for key, value := range request.GetQueryParams() {
		signParams[key] = value
	}
	for key, value := range request.GetFormParams() {
		signParams[key] = value
	}
	var paramKeySlice []string
	for key := range signParams {
		paramKeySlice = append(paramKeySlice, key)
	}
	sort.Strings(paramKeySlice)
	stringToSign = utils.GetUrlFormedMap(signParams)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	stringToSign = url.QueryEscape(stringToSign)
	stringToSign = request.GetMethod() + "&%2F&" + stringToSign
	return
}
