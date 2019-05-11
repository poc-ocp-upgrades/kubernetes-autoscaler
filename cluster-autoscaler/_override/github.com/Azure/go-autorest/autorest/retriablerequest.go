package autorest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func NewRetriableRequest(req *http.Request) *RetriableRequest {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &RetriableRequest{req: req}
}
func (rr *RetriableRequest) Request() *http.Request {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return rr.req
}
func (rr *RetriableRequest) prepareFromByteReader() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	b := []byte{}
	if rr.req.ContentLength > 0 {
		b = make([]byte, rr.req.ContentLength)
		_, err = io.ReadFull(rr.req.Body, b)
		if err != nil {
			return err
		}
	} else {
		b, err = ioutil.ReadAll(rr.req.Body)
		if err != nil {
			return err
		}
	}
	rr.br = bytes.NewReader(b)
	rr.req.Body = ioutil.NopCloser(rr.br)
	return err
}
