package adal

import (
	"net/http"
)

const (
	contentType			= "Content-Type"
	mimeTypeFormPost	= "application/x-www-form-urlencoded"
)

type Sender interface {
	Do(*http.Request) (*http.Response, error)
}
type SenderFunc func(*http.Request) (*http.Response, error)

func (sf SenderFunc) Do(r *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return sf(r)
}

type SendDecorator func(Sender) Sender

func CreateSender(decorators ...SendDecorator) Sender {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return DecorateSender(&http.Client{}, decorators...)
}
func DecorateSender(s Sender, decorators ...SendDecorator) Sender {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, decorate := range decorators {
		s = decorate(s)
	}
	return s
}
