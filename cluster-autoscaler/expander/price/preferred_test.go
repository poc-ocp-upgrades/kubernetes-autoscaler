package price

import (
	"testing"
	apiv1 "k8s.io/api/core/v1"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	"github.com/stretchr/testify/assert"
)

type testNodeLister struct{ list []*apiv1.Node }

func (n *testNodeLister) List() ([]*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.list, nil
}
func testPreferredNodeSingleCase(t *testing.T, currentNodes int, expectedNodeSize int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodes := []*apiv1.Node{}
	for i := 1; i <= currentNodes; i++ {
		nodes = append(nodes, BuildTestNode("n1", 1000, 1000))
	}
	provider := SimplePreferredNodeProvider{nodeLister: &testNodeLister{list: nodes}}
	node, err := provider.Node()
	assert.NoError(t, err)
	cpu := node.Status.Capacity[apiv1.ResourceCPU]
	assert.Equal(t, int64(expectedNodeSize), cpu.Value())
}
func TestPreferredNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testPreferredNodeSingleCase(t, 1, 1)
	testPreferredNodeSingleCase(t, 3, 2)
	testPreferredNodeSingleCase(t, 9, 4)
	testPreferredNodeSingleCase(t, 27, 8)
	testPreferredNodeSingleCase(t, 81, 16)
	testPreferredNodeSingleCase(t, 243, 32)
	testPreferredNodeSingleCase(t, 500, 32)
}
func TestSimpleNodeUnfitness(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	n1 := BuildTestNode("n1", 1000, 1000)
	n2 := BuildTestNode("n2", 2000, 2000)
	assert.Equal(t, 2.0, SimpleNodeUnfitness(n1, n2))
	assert.Equal(t, 2.0, SimpleNodeUnfitness(n2, n1))
	assert.Equal(t, 1.0, SimpleNodeUnfitness(n1, n1))
}
