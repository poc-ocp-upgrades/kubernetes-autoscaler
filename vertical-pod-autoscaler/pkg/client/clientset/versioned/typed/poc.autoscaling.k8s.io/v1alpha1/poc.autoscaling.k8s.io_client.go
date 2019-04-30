package v1alpha1

import (
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/poc.autoscaling.k8s.io/v1alpha1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type PocV1alpha1Interface interface {
	RESTClient() rest.Interface
	VerticalPodAutoscalersGetter
	VerticalPodAutoscalerCheckpointsGetter
}
type PocV1alpha1Client struct{ restClient rest.Interface }

func (c *PocV1alpha1Client) VerticalPodAutoscalers(namespace string) VerticalPodAutoscalerInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newVerticalPodAutoscalers(c, namespace)
}
func (c *PocV1alpha1Client) VerticalPodAutoscalerCheckpoints(namespace string) VerticalPodAutoscalerCheckpointInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newVerticalPodAutoscalerCheckpoints(c, namespace)
}
func NewForConfig(c *rest.Config) (*PocV1alpha1Client, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &PocV1alpha1Client{client}, nil
}
func NewForConfigOrDie(c *rest.Config) *PocV1alpha1Client {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}
func New(c rest.Interface) *PocV1alpha1Client {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &PocV1alpha1Client{c}
}
func setConfigDefaults(config *rest.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	return nil
}
func (c *PocV1alpha1Client) RESTClient() rest.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c == nil {
		return nil
	}
	return c.restClient
}
