package nodegroupset

import (
	"testing"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	testprovider "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/test"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func TestIsGkeNodeInfoSimilar(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("node1", 1000, 2000)
	n1.ObjectMeta.Labels["test-label"] = "test-value"
	n1.ObjectMeta.Labels["character"] = "winnie the pooh"
	n2 := BuildTestNode("node2", 1000, 2000)
	n2.ObjectMeta.Labels["test-label"] = "test-value"
	checkNodesSimilar(t, n1, n2, IsGkeNodeInfoSimilar, false)
	n1.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = ""
	n2.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = ""
	checkNodesSimilar(t, n1, n2, IsGkeNodeInfoSimilar, false)
	n1.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = ""
	n2.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah"
	checkNodesSimilar(t, n1, n2, IsGkeNodeInfoSimilar, false)
	delete(n1.ObjectMeta.Labels, "cloud.google.com/gke-nodepool")
	n2.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah"
	checkNodesSimilar(t, n1, n2, IsGkeNodeInfoSimilar, false)
	n1.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah1"
	n2.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah2"
	checkNodesSimilar(t, n1, n2, IsGkeNodeInfoSimilar, false)
	n1.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah"
	n2.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah"
	checkNodesSimilar(t, n1, n2, IsGkeNodeInfoSimilar, true)
}
func TestFindSimilarNodeGroupsGkeBasic(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	processor := &BalancingNodeGroupSetProcessor{Comparator: IsGkeNodeInfoSimilar}
	basicSimilarNodeGroupsTest(t, processor)
}
func TestFindSimilarNodeGroupsGkeByLabel(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	processor := &BalancingNodeGroupSetProcessor{Comparator: IsGkeNodeInfoSimilar}
	context := &context.AutoscalingContext{}
	n1 := BuildTestNode("n1", 1000, 1000)
	n2 := BuildTestNode("n2", 2000, 2000)
	provider := testprovider.NewTestCloudProvider(nil, nil)
	provider.AddNodeGroup("ng1", 1, 10, 1)
	provider.AddNodeGroup("ng2", 1, 10, 1)
	provider.AddNode("ng1", n1)
	provider.AddNode("ng2", n2)
	ni1 := schedulercache.NewNodeInfo()
	ni1.SetNode(n1)
	ni2 := schedulercache.NewNodeInfo()
	ni2.SetNode(n2)
	nodeInfosForGroups := map[string]*schedulercache.NodeInfo{"ng1": ni1, "ng2": ni2}
	ng1, _ := provider.NodeGroupForNode(n1)
	ng2, _ := provider.NodeGroupForNode(n2)
	context.CloudProvider = provider
	similar, err := processor.FindSimilarNodeGroups(context, ng1, nodeInfosForGroups)
	assert.NoError(t, err)
	assert.Equal(t, similar, []cloudprovider.NodeGroup{})
	n1.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah"
	n2.ObjectMeta.Labels["cloud.google.com/gke-nodepool"] = "blah"
	similar, err = processor.FindSimilarNodeGroups(context, ng1, nodeInfosForGroups)
	assert.NoError(t, err)
	assert.Equal(t, similar, []cloudprovider.NodeGroup{ng2})
}
