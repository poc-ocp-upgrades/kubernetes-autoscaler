package fake

import (
	v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAutoscalingV1beta1 struct{ *testing.Fake }

func (c *FakeAutoscalingV1beta1) VerticalPodAutoscalers(namespace string) v1beta1.VerticalPodAutoscalerInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeVerticalPodAutoscalers{c, namespace}
}
func (c *FakeAutoscalingV1beta1) VerticalPodAutoscalerCheckpoints(namespace string) v1beta1.VerticalPodAutoscalerCheckpointInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeVerticalPodAutoscalerCheckpoints{c, namespace}
}
func (c *FakeAutoscalingV1beta1) RESTClient() rest.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var ret *rest.RESTClient
	return ret
}
