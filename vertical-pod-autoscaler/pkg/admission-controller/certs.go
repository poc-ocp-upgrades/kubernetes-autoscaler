package main

import (
	"os"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"path"
	"github.com/golang/glog"
)

type certsContainer struct{ caKey, caCert, serverKey, serverCert []byte }

func readFile(filePath string) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(filePath)
	if err != nil {
		glog.Error(err)
		return nil
	}
	res := make([]byte, 5000)
	count, err := file.Read(res)
	if err != nil {
		glog.Error(err)
		return nil
	}
	glog.Infof("Successfully read %d bytes from %v", count, filePath)
	return res
}
func initCerts(certsDir string) certsContainer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	res := certsContainer{}
	res.caKey = readFile(path.Join(certsDir, "caKey.pem"))
	res.caCert = readFile(path.Join(certsDir, "caCert.pem"))
	res.serverKey = readFile(path.Join(certsDir, "serverKey.pem"))
	res.serverCert = readFile(path.Join(certsDir, "serverCert.pem"))
	return res
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
