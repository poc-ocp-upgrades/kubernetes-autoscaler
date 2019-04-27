package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	v1beta1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	testing "k8s.io/client-go/testing"
)

type FakeVerticalPodAutoscalers struct {
	Fake	*FakeAutoscalingV1beta1
	ns	string
}

var verticalpodautoscalersResource = schema.GroupVersionResource{Group: "autoscaling.k8s.io", Version: "v1beta1", Resource: "verticalpodautoscalers"}
var verticalpodautoscalersKind = schema.GroupVersionKind{Group: "autoscaling.k8s.io", Version: "v1beta1", Kind: "VerticalPodAutoscaler"}

func (c *FakeVerticalPodAutoscalers) Get(name string, options v1.GetOptions) (result *v1beta1.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewGetAction(verticalpodautoscalersResource, c.ns, name), &v1beta1.VerticalPodAutoscaler{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscaler), err
}
func (c *FakeVerticalPodAutoscalers) List(opts v1.ListOptions) (result *v1beta1.VerticalPodAutoscalerList, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewListAction(verticalpodautoscalersResource, verticalpodautoscalersKind, c.ns, opts), &v1beta1.VerticalPodAutoscalerList{})
	if obj == nil {
		return nil, err
	}
	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.VerticalPodAutoscalerList{}
	for _, item := range obj.(*v1beta1.VerticalPodAutoscalerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
func (c *FakeVerticalPodAutoscalers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Fake.InvokesWatch(testing.NewWatchAction(verticalpodautoscalersResource, c.ns, opts))
}
func (c *FakeVerticalPodAutoscalers) Create(verticalPodAutoscaler *v1beta1.VerticalPodAutoscaler) (result *v1beta1.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewCreateAction(verticalpodautoscalersResource, c.ns, verticalPodAutoscaler), &v1beta1.VerticalPodAutoscaler{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscaler), err
}
func (c *FakeVerticalPodAutoscalers) Update(verticalPodAutoscaler *v1beta1.VerticalPodAutoscaler) (result *v1beta1.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewUpdateAction(verticalpodautoscalersResource, c.ns, verticalPodAutoscaler), &v1beta1.VerticalPodAutoscaler{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscaler), err
}
func (c *FakeVerticalPodAutoscalers) UpdateStatus(verticalPodAutoscaler *v1beta1.VerticalPodAutoscaler) (*v1beta1.VerticalPodAutoscaler, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewUpdateSubresourceAction(verticalpodautoscalersResource, "status", c.ns, verticalPodAutoscaler), &v1beta1.VerticalPodAutoscaler{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscaler), err
}
func (c *FakeVerticalPodAutoscalers) Delete(name string, options *v1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := c.Fake.Invokes(testing.NewDeleteAction(verticalpodautoscalersResource, c.ns, name), &v1beta1.VerticalPodAutoscaler{})
	return err
}
func (c *FakeVerticalPodAutoscalers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	action := testing.NewDeleteCollectionAction(verticalpodautoscalersResource, c.ns, listOptions)
	_, err := c.Fake.Invokes(action, &v1beta1.VerticalPodAutoscalerList{})
	return err
}
func (c *FakeVerticalPodAutoscalers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VerticalPodAutoscaler, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewPatchSubresourceAction(verticalpodautoscalersResource, c.ns, name, data, subresources...), &v1beta1.VerticalPodAutoscaler{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscaler), err
}
