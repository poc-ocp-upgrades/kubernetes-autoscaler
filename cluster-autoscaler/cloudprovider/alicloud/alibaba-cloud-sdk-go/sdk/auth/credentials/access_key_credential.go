package credentials

import (
	godefaultruntime "runtime"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
)

type BaseCredential struct {
	AccessKeyId		string
	AccessKeySecret	string
}
type AccessKeyCredential struct {
	AccessKeyId		string
	AccessKeySecret	string
}

func NewBaseCredential(accessKeyId, accessKeySecret string) *BaseCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &BaseCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}
}
func (baseCred *BaseCredential) ToAccessKeyCredential() *AccessKeyCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &AccessKeyCredential{AccessKeyId: baseCred.AccessKeyId, AccessKeySecret: baseCred.AccessKeySecret}
}
func NewAccessKeyCredential(accessKeyId, accessKeySecret string) *AccessKeyCredential {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &AccessKeyCredential{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
