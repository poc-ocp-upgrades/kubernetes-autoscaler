package v1alpha1

import (
 v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 types "k8s.io/apimachinery/pkg/types"
 watch "k8s.io/apimachinery/pkg/watch"
 v1alpha1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/poc.autoscaling.k8s.io/v1alpha1"
 scheme "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/scheme"
 rest "k8s.io/client-go/rest"
)

type VerticalPodAutoscalerCheckpointsGetter interface {
 VerticalPodAutoscalerCheckpoints(namespace string) VerticalPodAutoscalerCheckpointInterface
}
type VerticalPodAutoscalerCheckpointInterface interface {
 Create(*v1alpha1.VerticalPodAutoscalerCheckpoint) (*v1alpha1.VerticalPodAutoscalerCheckpoint, error)
 Update(*v1alpha1.VerticalPodAutoscalerCheckpoint) (*v1alpha1.VerticalPodAutoscalerCheckpoint, error)
 Delete(name string, options *v1.DeleteOptions) error
 DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
 Get(name string, options v1.GetOptions) (*v1alpha1.VerticalPodAutoscalerCheckpoint, error)
 List(opts v1.ListOptions) (*v1alpha1.VerticalPodAutoscalerCheckpointList, error)
 Watch(opts v1.ListOptions) (watch.Interface, error)
 Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VerticalPodAutoscalerCheckpoint, err error)
 VerticalPodAutoscalerCheckpointExpansion
}
type verticalPodAutoscalerCheckpoints struct {
 client rest.Interface
 ns     string
}

func newVerticalPodAutoscalerCheckpoints(c *PocV1alpha1Client, namespace string) *verticalPodAutoscalerCheckpoints {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &verticalPodAutoscalerCheckpoints{client: c.RESTClient(), ns: namespace}
}
func (c *verticalPodAutoscalerCheckpoints) Get(name string, options v1.GetOptions) (result *v1alpha1.VerticalPodAutoscalerCheckpoint, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1alpha1.VerticalPodAutoscalerCheckpoint{}
 err = c.client.Get().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").Name(name).VersionedParams(&options, scheme.ParameterCodec).Do().Into(result)
 return
}
func (c *verticalPodAutoscalerCheckpoints) List(opts v1.ListOptions) (result *v1alpha1.VerticalPodAutoscalerCheckpointList, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1alpha1.VerticalPodAutoscalerCheckpointList{}
 err = c.client.Get().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").VersionedParams(&opts, scheme.ParameterCodec).Do().Into(result)
 return
}
func (c *verticalPodAutoscalerCheckpoints) Watch(opts v1.ListOptions) (watch.Interface, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 opts.Watch = true
 return c.client.Get().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").VersionedParams(&opts, scheme.ParameterCodec).Watch()
}
func (c *verticalPodAutoscalerCheckpoints) Create(verticalPodAutoscalerCheckpoint *v1alpha1.VerticalPodAutoscalerCheckpoint) (result *v1alpha1.VerticalPodAutoscalerCheckpoint, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1alpha1.VerticalPodAutoscalerCheckpoint{}
 err = c.client.Post().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").Body(verticalPodAutoscalerCheckpoint).Do().Into(result)
 return
}
func (c *verticalPodAutoscalerCheckpoints) Update(verticalPodAutoscalerCheckpoint *v1alpha1.VerticalPodAutoscalerCheckpoint) (result *v1alpha1.VerticalPodAutoscalerCheckpoint, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1alpha1.VerticalPodAutoscalerCheckpoint{}
 err = c.client.Put().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").Name(verticalPodAutoscalerCheckpoint.Name).Body(verticalPodAutoscalerCheckpoint).Do().Into(result)
 return
}
func (c *verticalPodAutoscalerCheckpoints) Delete(name string, options *v1.DeleteOptions) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return c.client.Delete().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").Name(name).Body(options).Do().Error()
}
func (c *verticalPodAutoscalerCheckpoints) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return c.client.Delete().Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").VersionedParams(&listOptions, scheme.ParameterCodec).Body(options).Do().Error()
}
func (c *verticalPodAutoscalerCheckpoints) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VerticalPodAutoscalerCheckpoint, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1alpha1.VerticalPodAutoscalerCheckpoint{}
 err = c.client.Patch(pt).Namespace(c.ns).Resource("verticalpodautoscalercheckpoints").SubResource(subresources...).Name(name).Body(data).Do().Into(result)
 return
}
