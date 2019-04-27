package autorest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type RetriableRequest struct {
	req	*http.Request
	rc	io.ReadCloser
	br	*bytes.Reader
}

func (rr *RetriableRequest) Prepare() (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if rr.req.Body != nil {
		if rr.rc != nil {
			rr.req.Body = rr.rc
		} else if rr.br != nil {
			_, err = rr.br.Seek(0, io.SeekStart)
			rr.req.Body = ioutil.NopCloser(rr.br)
		}
		if err != nil {
			return err
		}
		if rr.req.GetBody != nil {
			rr.rc, err = rr.req.GetBody()
			if err != nil {
				return err
			}
		} else if rr.br == nil {
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
	req.GetBody = nil
	req.ContentLength = 0
}
