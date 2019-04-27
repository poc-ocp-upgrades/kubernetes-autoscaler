package v1alpha1

import (
	time "time"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	poc_autoscaling_k8s_io_v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/poc.autoscaling.k8s.io/v1alpha1"
	versioned "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	internalinterfaces "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/listers/poc.autoscaling.k8s.io/v1alpha1"
	cache "k8s.io/client-go/tools/cache"
)

type VerticalPodAutoscalerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.VerticalPodAutoscalerLister
}
type verticalPodAutoscalerInformer struct {
	factory			internalinterfaces.SharedInformerFactory
	tweakListOptions	internalinterfaces.TweakListOptionsFunc
	namespace		string
}

func NewVerticalPodAutoscalerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewFilteredVerticalPodAutoscalerInformer(client, namespace, resyncPeriod, indexers, nil)
}
func NewFilteredVerticalPodAutoscalerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.NewSharedIndexInformer(&cache.ListWatch{ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
		if tweakListOptions != nil {
			tweakListOptions(&options)
		}
		return client.PocV1alpha1().VerticalPodAutoscalers(namespace).List(options)
	}, WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
		if tweakListOptions != nil {
			tweakListOptions(&options)
		}
		return client.PocV1alpha1().VerticalPodAutoscalers(namespace).Watch(options)
	}}, &poc_autoscaling_k8s_io_v1alpha1.VerticalPodAutoscaler{}, resyncPeriod, indexers)
}
func (f *verticalPodAutoscalerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewFilteredVerticalPodAutoscalerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}
func (f *verticalPodAutoscalerInformer) Informer() cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f.factory.InformerFor(&poc_autoscaling_k8s_io_v1alpha1.VerticalPodAutoscaler{}, f.defaultInformer)
}
func (f *verticalPodAutoscalerInformer) Lister() v1alpha1.VerticalPodAutoscalerLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v1alpha1.NewVerticalPodAutoscalerLister(f.Informer().GetIndexer())
}
