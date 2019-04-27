package history

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	correctResponse = `{
		"status":"success",
		"data":{
			"resultType": "matrix",
		        "result": []}}`
)

type mockHTTPGetter struct{ mock.Mock }

func (m mockHTTPGetter) Get(url string) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	args := m.Called(url)
	var returnArg http.Response
	if args.Get(0) != nil {
		returnArg = args.Get(0).(http.Response)
	}
	return &returnArg, args.Error(1)
}

type readerPseudoCloser struct{ *strings.Reader }

func (r readerPseudoCloser) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Errorf("readerPseudoCloser cannot really close anything")
}
func newReaderPseudoCloser(s string) readerPseudoCloser {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return readerPseudoCloser{strings.NewReader(s)}
}
func TestUrl(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	retryDelay = time.Hour
	mockGetter := mockHTTPGetter{}
	client := NewPrometheusClient(&mockGetter, "https://1.1.1.1")
	mockGetter.On("Get", "https://1.1.1.1/api/v1/query?query=up%7Ba%3Db%7D%5B2d%5D").Times(1).Return(http.Response{StatusCode: http.StatusOK, Body: newReaderPseudoCloser(correctResponse)}, nil)
	tss, err := client.GetTimeseries("up{a=b}[2d]")
	assert.Nil(t, err)
	assert.NotNil(t, tss)
	assert.Empty(t, tss)
}
func TestSuccessfulRetry(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	retryDelay = 100 * time.Millisecond
	mockGetter := mockHTTPGetter{}
	client := NewPrometheusClient(&mockGetter, "http://bla.com")
	mockGetter.On("Get", mock.AnythingOfType("string")).Times(1).Return(http.Response{StatusCode: http.StatusInternalServerError}, nil)
	mockGetter.On("Get", mock.AnythingOfType("string")).Times(1).Return(http.Response{StatusCode: http.StatusOK, Body: newReaderPseudoCloser(correctResponse)}, nil)
	tss, err := client.GetTimeseries("up")
	assert.Nil(t, err)
	assert.NotNil(t, tss)
	assert.Empty(t, tss)
}
func TestUnsuccessfulRetries(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	retryDelay = 10 * time.Millisecond
	mockGetter := mockHTTPGetter{}
	client := NewPrometheusClient(&mockGetter, "http://bla.com")
	mockGetter.On("Get", mock.AnythingOfType("string")).Times(numRetries).Return(http.Response{StatusCode: http.StatusInternalServerError}, nil)
	_, err := client.GetTimeseries("up")
	assert.NotNil(t, err)
}
