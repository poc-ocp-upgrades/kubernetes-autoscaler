package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: "poc.autoscaling.k8s.io", Version: "v1alpha1"}

func Resource(resource string) schema.GroupResource {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder		runtime.SchemeBuilder
	localSchemeBuilder	= &SchemeBuilder
	AddToScheme		= localSchemeBuilder.AddToScheme
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	localSchemeBuilder.Register(addKnownTypes)
}
func addKnownTypes(scheme *runtime.Scheme) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scheme.AddKnownTypes(SchemeGroupVersion, &VerticalPodAutoscaler{}, &VerticalPodAutoscalerList{}, &VerticalPodAutoscalerCheckpoint{}, &VerticalPodAutoscalerCheckpointList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
