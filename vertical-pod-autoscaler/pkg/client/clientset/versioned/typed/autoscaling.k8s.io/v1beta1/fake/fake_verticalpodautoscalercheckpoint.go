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

type FakeVerticalPodAutoscalerCheckpoints struct {
	Fake	*FakeAutoscalingV1beta1
	ns		string
}

var verticalpodautoscalercheckpointsResource = schema.GroupVersionResource{Group: "autoscaling.k8s.io", Version: "v1beta1", Resource: "verticalpodautoscalercheckpoints"}
var verticalpodautoscalercheckpointsKind = schema.GroupVersionKind{Group: "autoscaling.k8s.io", Version: "v1beta1", Kind: "VerticalPodAutoscalerCheckpoint"}

func (c *FakeVerticalPodAutoscalerCheckpoints) Get(name string, options v1.GetOptions) (result *v1beta1.VerticalPodAutoscalerCheckpoint, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewGetAction(verticalpodautoscalercheckpointsResource, c.ns, name), &v1beta1.VerticalPodAutoscalerCheckpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscalerCheckpoint), err
}
func (c *FakeVerticalPodAutoscalerCheckpoints) List(opts v1.ListOptions) (result *v1beta1.VerticalPodAutoscalerCheckpointList, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewListAction(verticalpodautoscalercheckpointsResource, verticalpodautoscalercheckpointsKind, c.ns, opts), &v1beta1.VerticalPodAutoscalerCheckpointList{})
	if obj == nil {
		return nil, err
	}
	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.VerticalPodAutoscalerCheckpointList{}
	for _, item := range obj.(*v1beta1.VerticalPodAutoscalerCheckpointList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}
func (c *FakeVerticalPodAutoscalerCheckpoints) Watch(opts v1.ListOptions) (watch.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.Fake.InvokesWatch(testing.NewWatchAction(verticalpodautoscalercheckpointsResource, c.ns, opts))
}
func (c *FakeVerticalPodAutoscalerCheckpoints) Create(verticalPodAutoscalerCheckpoint *v1beta1.VerticalPodAutoscalerCheckpoint) (result *v1beta1.VerticalPodAutoscalerCheckpoint, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewCreateAction(verticalpodautoscalercheckpointsResource, c.ns, verticalPodAutoscalerCheckpoint), &v1beta1.VerticalPodAutoscalerCheckpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscalerCheckpoint), err
}
func (c *FakeVerticalPodAutoscalerCheckpoints) Update(verticalPodAutoscalerCheckpoint *v1beta1.VerticalPodAutoscalerCheckpoint) (result *v1beta1.VerticalPodAutoscalerCheckpoint, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewUpdateAction(verticalpodautoscalercheckpointsResource, c.ns, verticalPodAutoscalerCheckpoint), &v1beta1.VerticalPodAutoscalerCheckpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscalerCheckpoint), err
}
func (c *FakeVerticalPodAutoscalerCheckpoints) Delete(name string, options *v1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := c.Fake.Invokes(testing.NewDeleteAction(verticalpodautoscalercheckpointsResource, c.ns, name), &v1beta1.VerticalPodAutoscalerCheckpoint{})
	return err
}
func (c *FakeVerticalPodAutoscalerCheckpoints) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	action := testing.NewDeleteCollectionAction(verticalpodautoscalercheckpointsResource, c.ns, listOptions)
	_, err := c.Fake.Invokes(action, &v1beta1.VerticalPodAutoscalerCheckpointList{})
	return err
}
func (c *FakeVerticalPodAutoscalerCheckpoints) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.VerticalPodAutoscalerCheckpoint, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := c.Fake.Invokes(testing.NewPatchSubresourceAction(verticalpodautoscalercheckpointsResource, c.ns, name, data, subresources...), &v1beta1.VerticalPodAutoscalerCheckpoint{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VerticalPodAutoscalerCheckpoint), err
}
