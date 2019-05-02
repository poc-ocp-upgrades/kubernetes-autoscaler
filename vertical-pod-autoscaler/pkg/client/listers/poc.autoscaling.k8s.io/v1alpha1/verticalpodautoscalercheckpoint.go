package v1alpha1

import (
 "k8s.io/apimachinery/pkg/api/errors"
 "k8s.io/apimachinery/pkg/labels"
 v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/poc.autoscaling.k8s.io/v1alpha1"
 "k8s.io/client-go/tools/cache"
)

type VerticalPodAutoscalerCheckpointLister interface {
 List(selector labels.Selector) (ret []*v1alpha1.VerticalPodAutoscalerCheckpoint, err error)
 VerticalPodAutoscalerCheckpoints(namespace string) VerticalPodAutoscalerCheckpointNamespaceLister
 VerticalPodAutoscalerCheckpointListerExpansion
}
type verticalPodAutoscalerCheckpointLister struct{ indexer cache.Indexer }

func NewVerticalPodAutoscalerCheckpointLister(indexer cache.Indexer) VerticalPodAutoscalerCheckpointLister {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &verticalPodAutoscalerCheckpointLister{indexer: indexer}
}
func (s *verticalPodAutoscalerCheckpointLister) List(selector labels.Selector) (ret []*v1alpha1.VerticalPodAutoscalerCheckpoint, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 err = cache.ListAll(s.indexer, selector, func(m interface{}) {
  ret = append(ret, m.(*v1alpha1.VerticalPodAutoscalerCheckpoint))
 })
 return ret, err
}
func (s *verticalPodAutoscalerCheckpointLister) VerticalPodAutoscalerCheckpoints(namespace string) VerticalPodAutoscalerCheckpointNamespaceLister {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return verticalPodAutoscalerCheckpointNamespaceLister{indexer: s.indexer, namespace: namespace}
}

type VerticalPodAutoscalerCheckpointNamespaceLister interface {
 List(selector labels.Selector) (ret []*v1alpha1.VerticalPodAutoscalerCheckpoint, err error)
 Get(name string) (*v1alpha1.VerticalPodAutoscalerCheckpoint, error)
 VerticalPodAutoscalerCheckpointNamespaceListerExpansion
}
type verticalPodAutoscalerCheckpointNamespaceLister struct {
 indexer   cache.Indexer
 namespace string
}

func (s verticalPodAutoscalerCheckpointNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.VerticalPodAutoscalerCheckpoint, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
  ret = append(ret, m.(*v1alpha1.VerticalPodAutoscalerCheckpoint))
 })
 return ret, err
}
func (s verticalPodAutoscalerCheckpointNamespaceLister) Get(name string) (*v1alpha1.VerticalPodAutoscalerCheckpoint, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
 if err != nil {
  return nil, err
 }
 if !exists {
  return nil, errors.NewNotFound(v1alpha1.Resource("verticalpodautoscalercheckpoint"), name)
 }
 return obj.(*v1alpha1.VerticalPodAutoscalerCheckpoint), nil
}
