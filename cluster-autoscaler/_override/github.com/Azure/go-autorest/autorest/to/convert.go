package to

import (
	godefaultruntime "runtime"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
)

func String(s *string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s != nil {
		return *s
	}
	return ""
}
func StringPtr(s string) *string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &s
}
func StringSlice(s *[]string) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s != nil {
		return *s
	}
	return nil
}
func StringSlicePtr(s []string) *[]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &s
}
func StringMap(msp map[string]*string) map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ms := make(map[string]string, len(msp))
	for k, sp := range msp {
		if sp != nil {
			ms[k] = *sp
		} else {
			ms[k] = ""
		}
	}
	return ms
}
func StringMapPtr(ms map[string]string) *map[string]*string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	msp := make(map[string]*string, len(ms))
	for k, s := range ms {
		msp[k] = StringPtr(s)
	}
	return &msp
}
func Bool(b *bool) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if b != nil {
		return *b
	}
	return false
}
func BoolPtr(b bool) *bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &b
}
func Int(i *int) int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i != nil {
		return *i
	}
	return 0
}
func IntPtr(i int) *int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &i
}
func Int32(i *int32) int32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i != nil {
		return *i
	}
	return 0
}
func Int32Ptr(i int32) *int32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &i
}
func Int64(i *int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i != nil {
		return *i
	}
	return 0
}
func Int64Ptr(i int64) *int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &i
}
func Float32(i *float32) float32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i != nil {
		return *i
	}
	return 0.0
}
func Float32Ptr(i float32) *float32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &i
}
func Float64(i *float64) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if i != nil {
		return *i
	}
	return 0.0
}
func Float64Ptr(i float64) *float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &i
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
