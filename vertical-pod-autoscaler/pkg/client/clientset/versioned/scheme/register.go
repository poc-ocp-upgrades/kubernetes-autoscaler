package scheme

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	autoscalingv1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	pocv1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/poc.autoscaling.k8s.io/v1alpha1"
)

var Scheme = runtime.NewScheme()
var Codecs = serializer.NewCodecFactory(Scheme)
var ParameterCodec = runtime.NewParameterCodec(Scheme)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	AddToScheme(Scheme)
}
func AddToScheme(scheme *runtime.Scheme) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	autoscalingv1beta1.AddToScheme(scheme)
	pocv1alpha1.AddToScheme(scheme)
}
