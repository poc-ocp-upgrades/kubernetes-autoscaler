package price

import (
	"fmt"
	"testing"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
	apiv1 "k8s.io/api/core/v1"
	testprovider "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/test"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

type testPricingModel struct {
	nodePrice	map[string]float64
	podPrice	map[string]float64
}

func (tpm *testPricingModel) NodePrice(node *apiv1.Node, startTime time.Time, endTime time.Time) (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if price, found := tpm.nodePrice[node.Name]; found {
		return price, nil
	}
	return 0.0, fmt.Errorf("price for node %v not found", node.Name)
}
func (tpm *testPricingModel) PodPrice(node *apiv1.Pod, startTime time.Time, endTime time.Time) (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if price, found := tpm.podPrice[node.Name]; found {
		return price, nil
	}
	return 0.0, fmt.Errorf("price for pod %v not found", node.Name)
}

type testPreferredNodeProvider struct{ preferred *apiv1.Node }

func (tpnp *testPreferredNodeProvider) Node() (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tpnp.preferred, nil
}
func TestPriceExpander(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("n1", 1000, 1000)
	n2 := BuildTestNode("n2", 4000, 1000)
	n3 := BuildTestNode("n3", 4000, 1000)
	p1 := BuildTestPod("p1", 1000, 0)
	p2 := BuildTestPod("p2", 500, 0)
	provider := testprovider.NewTestCloudProvider(nil, nil)
	provider.AddNodeGroup("ng1", 1, 10, 1)
	provider.AddNodeGroup("ng2", 1, 10, 1)
	provider.AddNode("ng1", n1)
	provider.AddNode("ng2", n2)
	ng1, _ := provider.NodeGroupForNode(n1)
	ng2, _ := provider.NodeGroupForNode(n2)
	ng3, _ := provider.NewNodeGroup("MT1", nil, nil, nil, nil)
	ni1 := schedulercache.NewNodeInfo()
	ni1.SetNode(n1)
	ni2 := schedulercache.NewNodeInfo()
	ni2.SetNode(n2)
	ni3 := schedulercache.NewNodeInfo()
	ni3.SetNode(n3)
	nodeInfosForGroups := map[string]*schedulercache.NodeInfo{"ng1": ni1, "ng2": ni2}
	options := []expander.Option{{NodeGroup: ng1, NodeCount: 2, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng1"}, {NodeGroup: ng2, NodeCount: 1, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng2"}}
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 20.0, "n2": 200.0}}, &testPreferredNodeProvider{preferred: buildNode(2000, units.GiB)}, SimpleNodeUnfitness).BestOption(options, nodeInfosForGroups).Debug, "ng1")
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 50.0, "n2": 200.0}}, &testPreferredNodeProvider{preferred: buildNode(4000, units.GiB)}, SimpleNodeUnfitness).BestOption(options, nodeInfosForGroups).Debug, "ng2")
	options1b := []expander.Option{{NodeGroup: ng1, NodeCount: 80, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng1"}, {NodeGroup: ng2, NodeCount: 40, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng2"}}
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 20.0, "n2": 200.0}}, &testPreferredNodeProvider{preferred: buildNode(4000, units.GiB)}, SimpleNodeUnfitness).BestOption(options1b, nodeInfosForGroups).Debug, "ng1")
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 200.0, "n2": 100.0}}, &testPreferredNodeProvider{preferred: buildNode(2000, units.GiB)}, SimpleNodeUnfitness).BestOption(options, nodeInfosForGroups).Debug, "ng2")
	options2 := []expander.Option{{NodeGroup: ng1, NodeCount: 2, Pods: []*apiv1.Pod{p1}, Debug: "ng1"}, {NodeGroup: ng2, NodeCount: 1, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng2"}}
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 200.0, "n2": 200.0}}, &testPreferredNodeProvider{preferred: buildNode(2000, units.GiB)}, SimpleNodeUnfitness).BestOption(options2, nodeInfosForGroups).Debug, "ng2")
	assert.Nil(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{}, nodePrice: map[string]float64{}}, &testPreferredNodeProvider{preferred: buildNode(2000, units.GiB)}, SimpleNodeUnfitness).BestOption(options2, nodeInfosForGroups))
	nodeInfosForGroups["autoprovisioned-MT1"] = ni3
	options3 := []expander.Option{{NodeGroup: ng1, NodeCount: 2, Pods: []*apiv1.Pod{p1}, Debug: "ng1"}, {NodeGroup: ng2, NodeCount: 1, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng2"}, {NodeGroup: ng3, NodeCount: 1, Pods: []*apiv1.Pod{p1, p2}, Debug: "ng3"}}
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 200.0, "n2": 200.0, "n3": 200.0}}, &testPreferredNodeProvider{preferred: buildNode(2000, units.GiB)}, SimpleNodeUnfitness).BestOption(options3, nodeInfosForGroups).Debug, "ng2")
	assert.Contains(t, NewStrategy(&testPricingModel{podPrice: map[string]float64{"p1": 20.0, "p2": 10.0, "stabilize": 10}, nodePrice: map[string]float64{"n1": 200.0, "n2": 200.0, "n3": 90.0}}, &testPreferredNodeProvider{preferred: buildNode(2000, units.GiB)}, SimpleNodeUnfitness).BestOption(options3, nodeInfosForGroups).Debug, "ng3")
}
