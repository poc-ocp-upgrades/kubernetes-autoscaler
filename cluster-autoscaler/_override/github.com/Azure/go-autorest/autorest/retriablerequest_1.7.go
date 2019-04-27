package autorest

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type RetriableRequest struct {
	req	*http.Request
	br	*bytes.Reader
}

func (rr *RetriableRequest) Prepare() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rr.req.Body != nil {
		if rr.br != nil {
			_, err = rr.br.Seek(0, 0)
			rr.req.Body = ioutil.NopCloser(rr.br)
		}
		if err != nil {
			return err
		}
		if rr.br == nil {
			err = rr.prepareFromByteReader()
		}
	}
	return err
}
func removeRequestBody(req *http.Request) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	req.Body = nil
	req.ContentLength = 0
}
