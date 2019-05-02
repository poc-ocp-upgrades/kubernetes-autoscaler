package v1beta1

import (
 v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 types "k8s.io/apimachinery/pkg/types"
 watch "k8s.io/apimachinery/pkg/watch"
 v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
 scheme "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/scheme"
 rest "k8s.io/client-go/rest"
)

type VerticalPodAutoscalersGetter interface {
 VerticalPodAutoscalers(namespace string) VerticalPodAutoscalerInterface
}
type VerticalPodAutoscalerInterface interface {
 Create(*v1beta1.VerticalPodAutoscaler) (*v1beta1.VerticalPodAutoscaler, error)
 Update(*v1beta1.VerticalPodAutoscaler) (*v1beta1.VerticalPodAutoscaler, error)
 UpdateStatus(*v1beta1.VerticalPodAutoscaler) (*v1beta1.VerticalPodAutoscaler, error)
 Delete(name string, options *v1.DeleteOptions) error
 DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
 Get(name string, options v1.GetOptions) (*v1beta1.VerticalPodAutoscaler, error)
 List(opts v1.ListOptions) (*v1beta1.VerticalPodAutoscalerList, error)
 Watch(opts v1.ListOptions) (watch.Interface, error)
 Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VerticalPodAutoscaler, err error)
 VerticalPodAutoscalerExpansion
}
type verticalPodAutoscalers struct {
 client rest.Interface
 ns     string
}

func newVerticalPodAutoscalers(c *AutoscalingV1beta1Client, namespace string) *verticalPodAutoscalers {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &verticalPodAutoscalers{client: c.RESTClient(), ns: namespace}
}
func (c *verticalPodAutoscalers) Get(name string, options v1.GetOptions) (result *v1beta1.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1beta1.VerticalPodAutoscaler{}
 err = c.client.Get().Namespace(c.ns).Resource("verticalpodautoscalers").Name(name).VersionedParams(&options, scheme.ParameterCodec).Do().Into(result)
 return
}
func (c *verticalPodAutoscalers) List(opts v1.ListOptions) (result *v1beta1.VerticalPodAutoscalerList, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1beta1.VerticalPodAutoscalerList{}
 err = c.client.Get().Namespace(c.ns).Resource("verticalpodautoscalers").VersionedParams(&opts, scheme.ParameterCodec).Do().Into(result)
 return
}
func (c *verticalPodAutoscalers) Watch(opts v1.ListOptions) (watch.Interface, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 opts.Watch = true
 return c.client.Get().Namespace(c.ns).Resource("verticalpodautoscalers").VersionedParams(&opts, scheme.ParameterCodec).Watch()
}
func (c *verticalPodAutoscalers) Create(verticalPodAutoscaler *v1beta1.VerticalPodAutoscaler) (result *v1beta1.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1beta1.VerticalPodAutoscaler{}
 err = c.client.Post().Namespace(c.ns).Resource("verticalpodautoscalers").Body(verticalPodAutoscaler).Do().Into(result)
 return
}
func (c *verticalPodAutoscalers) Update(verticalPodAutoscaler *v1beta1.VerticalPodAutoscaler) (result *v1beta1.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1beta1.VerticalPodAutoscaler{}
 err = c.client.Put().Namespace(c.ns).Resource("verticalpodautoscalers").Name(verticalPodAutoscaler.Name).Body(verticalPodAutoscaler).Do().Into(result)
 return
}
func (c *verticalPodAutoscalers) UpdateStatus(verticalPodAutoscaler *v1beta1.VerticalPodAutoscaler) (result *v1beta1.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1beta1.VerticalPodAutoscaler{}
 err = c.client.Put().Namespace(c.ns).Resource("verticalpodautoscalers").Name(verticalPodAutoscaler.Name).SubResource("status").Body(verticalPodAutoscaler).Do().Into(result)
 return
}
func (c *verticalPodAutoscalers) Delete(name string, options *v1.DeleteOptions) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return c.client.Delete().Namespace(c.ns).Resource("verticalpodautoscalers").Name(name).Body(options).Do().Error()
}
func (c *verticalPodAutoscalers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return c.client.Delete().Namespace(c.ns).Resource("verticalpodautoscalers").VersionedParams(&listOptions, scheme.ParameterCodec).Body(options).Do().Error()
}
func (c *verticalPodAutoscalers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VerticalPodAutoscaler, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result = &v1beta1.VerticalPodAutoscaler{}
 err = c.client.Patch(pt).Namespace(c.ns).Resource("verticalpodautoscalers").SubResource(subresources...).Name(name).Body(data).Do().Into(result)
 return
}
