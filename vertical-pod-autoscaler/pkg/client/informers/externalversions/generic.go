package externalversions

import (
	"fmt"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/poc.autoscaling.k8s.io/v1alpha1"
	cache "k8s.io/client-go/tools/cache"
)

type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}
type genericInformer struct {
	informer	cache.SharedIndexInformer
	resource	schema.GroupResource
}

func (f *genericInformer) Informer() cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f.informer
}
func (f *genericInformer) Lister() cache.GenericLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch resource {
	case v1beta1.SchemeGroupVersion.WithResource("verticalpodautoscalers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Autoscaling().V1beta1().VerticalPodAutoscalers().Informer()}, nil
	case v1beta1.SchemeGroupVersion.WithResource("verticalpodautoscalercheckpoints"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Autoscaling().V1beta1().VerticalPodAutoscalerCheckpoints().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("verticalpodautoscalers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Poc().V1alpha1().VerticalPodAutoscalers().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("verticalpodautoscalercheckpoints"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Poc().V1alpha1().VerticalPodAutoscalerCheckpoints().Informer()}, nil
	}
	return nil, fmt.Errorf("no informer found for %v", resource)
}
