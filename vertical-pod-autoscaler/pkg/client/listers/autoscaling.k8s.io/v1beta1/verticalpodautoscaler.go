package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/client-go/tools/cache"
)

type VerticalPodAutoscalerLister interface {
	List(selector labels.Selector) (ret []*v1beta1.VerticalPodAutoscaler, err error)
	VerticalPodAutoscalers(namespace string) VerticalPodAutoscalerNamespaceLister
	VerticalPodAutoscalerListerExpansion
}
type verticalPodAutoscalerLister struct{ indexer cache.Indexer }

func NewVerticalPodAutoscalerLister(indexer cache.Indexer) VerticalPodAutoscalerLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &verticalPodAutoscalerLister{indexer: indexer}
}
func (s *verticalPodAutoscalerLister) List(selector labels.Selector) (ret []*v1beta1.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.VerticalPodAutoscaler))
	})
	return ret, err
}
func (s *verticalPodAutoscalerLister) VerticalPodAutoscalers(namespace string) VerticalPodAutoscalerNamespaceLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return verticalPodAutoscalerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

type VerticalPodAutoscalerNamespaceLister interface {
	List(selector labels.Selector) (ret []*v1beta1.VerticalPodAutoscaler, err error)
	Get(name string) (*v1beta1.VerticalPodAutoscaler, error)
	VerticalPodAutoscalerNamespaceListerExpansion
}
type verticalPodAutoscalerNamespaceLister struct {
	indexer		cache.Indexer
	namespace	string
}

func (s verticalPodAutoscalerNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.VerticalPodAutoscaler))
	})
	return ret, err
}
func (s verticalPodAutoscalerNamespaceLister) Get(name string) (*v1beta1.VerticalPodAutoscaler, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("verticalpodautoscaler"), name)
	}
	return obj.(*v1beta1.VerticalPodAutoscaler), nil
}
