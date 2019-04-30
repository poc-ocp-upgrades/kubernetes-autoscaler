package v1beta1

import (
	time "time"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	autoscaling_k8s_io_v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	versioned "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	internalinterfaces "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/informers/externalversions/internalinterfaces"
	v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/listers/autoscaling.k8s.io/v1beta1"
	cache "k8s.io/client-go/tools/cache"
)

type VerticalPodAutoscalerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.VerticalPodAutoscalerLister
}
type verticalPodAutoscalerInformer struct {
	factory			internalinterfaces.SharedInformerFactory
	tweakListOptions	internalinterfaces.TweakListOptionsFunc
	namespace		string
}

func NewVerticalPodAutoscalerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewFilteredVerticalPodAutoscalerInformer(client, namespace, resyncPeriod, indexers, nil)
}
func NewFilteredVerticalPodAutoscalerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cache.NewSharedIndexInformer(&cache.ListWatch{ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
		if tweakListOptions != nil {
			tweakListOptions(&options)
		}
		return client.AutoscalingV1beta1().VerticalPodAutoscalers(namespace).List(options)
	}, WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
		if tweakListOptions != nil {
			tweakListOptions(&options)
		}
		return client.AutoscalingV1beta1().VerticalPodAutoscalers(namespace).Watch(options)
	}}, &autoscaling_k8s_io_v1beta1.VerticalPodAutoscaler{}, resyncPeriod, indexers)
}
func (f *verticalPodAutoscalerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewFilteredVerticalPodAutoscalerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}
func (f *verticalPodAutoscalerInformer) Informer() cache.SharedIndexInformer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return f.factory.InformerFor(&autoscaling_k8s_io_v1beta1.VerticalPodAutoscaler{}, f.defaultInformer)
}
func (f *verticalPodAutoscalerInformer) Lister() v1beta1.VerticalPodAutoscalerLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v1beta1.NewVerticalPodAutoscalerLister(f.Informer().GetIndexer())
}
