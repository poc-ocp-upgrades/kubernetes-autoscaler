package gce

import (
	"math"
	"testing"
	"time"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
	"github.com/stretchr/testify/assert"
)

func TestGetNodePrice(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels1, _ := BuildGenericLabels(GceRef{Name: "kubernetes-minion-group", Project: "mwielgus-proj", Zone: "us-central1-b"}, "n1-standard-8", "sillyname")
	labels2, _ := BuildGenericLabels(GceRef{Name: "kubernetes-minion-group", Project: "mwielgus-proj", Zone: "us-central1-b"}, "n1-standard-8", "sillyname")
	labels2[preemptibleLabel] = "true"
	model := &GcePriceModel{}
	now := time.Now()
	node1 := BuildTestNode("sillyname1", 8000, 30*units.GiB)
	node1.Labels = labels1
	price1, err := model.NodePrice(node1, now, now.Add(time.Hour))
	assert.NoError(t, err)
	node2 := BuildTestNode("sillyname2", 8000, 30*units.GiB)
	node2.Labels = labels2
	price2, err := model.NodePrice(node2, now, now.Add(time.Hour))
	assert.NoError(t, err)
	assert.True(t, price1 > 3*price2)
	node3 := BuildTestNode("sillyname3", 8000, 30*units.GiB)
	price3, err := model.NodePrice(node3, now, now.Add(time.Hour))
	assert.NoError(t, err)
	assert.True(t, price1 < price3)
	assert.True(t, price1*1.2 > price3)
	node4 := BuildTestNode("sillyname4", 8000, 30*units.GiB)
	node4.Status.Capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(1, resource.DecimalSI)
	node4.Labels = labels1
	price4, err := model.NodePrice(node4, now, now.Add(time.Hour))
	node5 := BuildTestNode("sillyname5", 8000, 30*units.GiB)
	node5.Labels = labels2
	node5.Status.Capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(1, resource.DecimalSI)
	price5, err := model.NodePrice(node5, now, now.Add(time.Hour))
	assert.True(t, price4 > price5)
	assert.True(t, price4 < 1.5*price5)
	assert.True(t, price4 > 2*price1)
	node6 := BuildTestNode("sillyname6", 1000, 3750*units.MiB)
	price6, err := model.NodePrice(node6, now, now.Add(time.Hour))
	assert.NoError(t, err)
	assert.True(t, math.Abs(price3-8*price6) < 0.1)
}
func TestGetPodPrice(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pod1 := BuildTestPod("a1", 100, 500*units.MiB)
	pod2 := BuildTestPod("a2", 2*100, 2*500*units.MiB)
	model := &GcePriceModel{}
	now := time.Now()
	price1, err := model.PodPrice(pod1, now, now.Add(time.Hour))
	assert.NoError(t, err)
	price2, err := model.PodPrice(pod2, now, now.Add(time.Hour))
	assert.NoError(t, err)
	assert.True(t, math.Abs(price1*2-price2) < 0.001)
}
