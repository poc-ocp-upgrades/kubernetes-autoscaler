package responses

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/errors"
	"net/http"
	"strings"
)

type AcsResponse interface {
	IsSuccess() bool
	GetHttpStatus() int
	GetHttpHeaders() map[string][]string
	GetHttpContentString() string
	GetHttpContentBytes() []byte
	GetOriginHttpResponse() *http.Response
	parseFromHttpResponse(httpResponse *http.Response) error
}

func Unmarshal(response AcsResponse, httpResponse *http.Response, format string) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = response.parseFromHttpResponse(httpResponse)
	if err != nil {
		return
	}
	if !response.IsSuccess() {
		err = errors.NewServerError(response.GetHttpStatus(), response.GetHttpContentString(), "")
		return
	}
	if _, isCommonResponse := response.(CommonResponse); isCommonResponse {
		return
	}
	if len(response.GetHttpContentBytes()) == 0 {
		return
	}
	if strings.ToUpper(format) == "JSON" {
		initJsonParserOnce()
		err = jsonParser.Unmarshal(response.GetHttpContentBytes(), response)
		if err != nil {
			err = errors.NewClientError(errors.JsonUnmarshalErrorCode, errors.JsonUnmarshalErrorMessage, err)
		}
	} else if strings.ToUpper(format) == "XML" {
		err = xml.Unmarshal(response.GetHttpContentBytes(), response)
	}
	return
}

type BaseResponse struct {
	httpStatus		int
	httpHeaders		map[string][]string
	httpContentString	string
	httpContentBytes	[]byte
	originHttpResponse	*http.Response
}

func (baseResponse *BaseResponse) GetHttpStatus() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return baseResponse.httpStatus
}
func (baseResponse *BaseResponse) GetHttpHeaders() map[string][]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return baseResponse.httpHeaders
}
func (baseResponse *BaseResponse) GetHttpContentString() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return baseResponse.httpContentString
}
func (baseResponse *BaseResponse) GetHttpContentBytes() []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return baseResponse.httpContentBytes
}
func (baseResponse *BaseResponse) GetOriginHttpResponse() *http.Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return baseResponse.originHttpResponse
}
func (baseResponse *BaseResponse) IsSuccess() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if baseResponse.GetHttpStatus() >= 200 && baseResponse.GetHttpStatus() < 300 {
		return true
	}
	return false
}
func (baseResponse *BaseResponse) parseFromHttpResponse(httpResponse *http.Response) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}
	baseResponse.httpStatus = httpResponse.StatusCode
	baseResponse.httpHeaders = httpResponse.Header
	baseResponse.httpContentBytes = body
	baseResponse.httpContentString = string(body)
	baseResponse.originHttpResponse = httpResponse
	return
}
func (baseResponse *BaseResponse) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	resultBuilder := bytes.Buffer{}
	resultBuilder.WriteString("\n")
	resultBuilder.WriteString(fmt.Sprintf("%s %s\n", baseResponse.originHttpResponse.Proto, baseResponse.originHttpResponse.Status))
	for key, value := range baseResponse.httpHeaders {
		resultBuilder.WriteString(key + ": " + strings.Join(value, ";") + "\n")
	}
	resultBuilder.WriteString("\n")
	resultBuilder.WriteString(baseResponse.httpContentString + "\n")
	return resultBuilder.String()
}

type CommonResponse struct{ *BaseResponse }

func NewCommonResponse() (response *CommonResponse) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &CommonResponse{BaseResponse: &BaseResponse{}}
}
