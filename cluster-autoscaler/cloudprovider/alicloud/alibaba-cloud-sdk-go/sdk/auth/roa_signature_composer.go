package auth

import (
	"bytes"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/requests"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/utils"
	"sort"
	"strings"
)

func signRoaRequest(request requests.AcsRequest, signer Signer, regionId string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	completeROASignParams(request, signer, regionId)
	stringToSign := buildRoaStringToSign(request)
	request.SetStringToSign(stringToSign)
	signature := signer.Sign(stringToSign, "")
	accessKeyId, err := signer.GetAccessKeyId()
	if err != nil {
		return nil
	}
	request.GetHeaders()["Authorization"] = "acs " + accessKeyId + ":" + signature
	return
}
func completeROASignParams(request requests.AcsRequest, signer Signer, regionId string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	headerParams := request.GetHeaders()
	queryParams := request.GetQueryParams()
	if extraParam := signer.GetExtraParam(); extraParam != nil {
		for key, value := range extraParam {
			if key == "SecurityToken" {
				headerParams["x-acs-security-token"] = value
				continue
			}
			queryParams[key] = value
		}
	}
	headerParams["Date"] = utils.GetTimeInFormatRFC2616()
	headerParams["x-acs-signature-method"] = signer.GetName()
	headerParams["x-acs-signature-version"] = signer.GetVersion()
	if request.GetFormParams() != nil && len(request.GetFormParams()) > 0 {
		formString := utils.GetUrlFormedMap(request.GetFormParams())
		request.SetContent([]byte(formString))
		headerParams["Content-Type"] = requests.Form
	}
	contentMD5 := utils.GetMD5Base64(request.GetContent())
	headerParams["Content-MD5"] = contentMD5
	if _, contains := headerParams["Content-Type"]; !contains {
		headerParams["Content-Type"] = requests.Raw
	}
	switch format := request.GetAcceptFormat(); format {
	case "JSON":
		headerParams["Accept"] = requests.Json
	case "XML":
		headerParams["Accept"] = requests.Xml
	default:
		headerParams["Accept"] = requests.Raw
	}
}
func buildRoaStringToSign(request requests.AcsRequest) (stringToSign string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	headers := request.GetHeaders()
	stringToSignBuilder := bytes.Buffer{}
	stringToSignBuilder.WriteString(request.GetMethod())
	stringToSignBuilder.WriteString(requests.HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Accept", requests.HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Content-MD5", requests.HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Content-Type", requests.HeaderSeparator)
	appendIfContain(headers, &stringToSignBuilder, "Date", requests.HeaderSeparator)
	var acsHeaders []string
	for key := range headers {
		if strings.HasPrefix(key, "x-acs-") {
			acsHeaders = append(acsHeaders, key)
		}
	}
	sort.Strings(acsHeaders)
	for _, key := range acsHeaders {
		stringToSignBuilder.WriteString(key + ":" + headers[key])
		stringToSignBuilder.WriteString(requests.HeaderSeparator)
	}
	stringToSignBuilder.WriteString(request.BuildQueries())
	stringToSign = stringToSignBuilder.String()
	return
}
func appendIfContain(sourceMap map[string]string, target *bytes.Buffer, key, separator string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if value, contain := sourceMap[key]; contain && len(value) > 0 {
		target.WriteString(sourceMap[key])
		target.WriteString(separator)
	}
}
