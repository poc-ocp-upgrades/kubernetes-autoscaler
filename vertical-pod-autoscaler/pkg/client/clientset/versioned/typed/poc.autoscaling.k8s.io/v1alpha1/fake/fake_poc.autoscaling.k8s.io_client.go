package fake

import (
	v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/poc.autoscaling.k8s.io/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakePocV1alpha1 struct{ *testing.Fake }

func (c *FakePocV1alpha1) VerticalPodAutoscalers(namespace string) v1alpha1.VerticalPodAutoscalerInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeVerticalPodAutoscalers{c, namespace}
}
func (c *FakePocV1alpha1) VerticalPodAutoscalerCheckpoints(namespace string) v1alpha1.VerticalPodAutoscalerCheckpointInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeVerticalPodAutoscalerCheckpoints{c, namespace}
}
func (c *FakePocV1alpha1) RESTClient() rest.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var ret *rest.RESTClient
	return ret
}
