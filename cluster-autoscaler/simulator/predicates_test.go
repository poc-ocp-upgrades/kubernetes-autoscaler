package simulator

import (
	"strings"
	"testing"
	"time"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func TestPredicates(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p1 := BuildTestPod("p1", 450, 500000)
	p2 := BuildTestPod("p2", 600, 500000)
	p3 := BuildTestPod("p3", 8000, 0)
	p4 := BuildTestPod("p4", 500, 500000)
	ni1 := schedulercache.NewNodeInfo(p1)
	ni2 := schedulercache.NewNodeInfo()
	nodeInfos := map[string]*schedulercache.NodeInfo{"n1": ni1, "n2": ni2}
	node1 := BuildTestNode("n1", 1000, 2000000)
	node2 := BuildTestNode("n2", 1000, 2000000)
	SetNodeReadyState(node1, true, time.Time{})
	SetNodeReadyState(node2, true, time.Time{})
	ni1.SetNode(node1)
	ni2.SetNode(node2)
	predicateChecker := NewTestPredicateChecker()
	r1, err := predicateChecker.FitsAny(p2, nodeInfos)
	assert.NoError(t, err)
	assert.Equal(t, "n2", r1)
	_, err = predicateChecker.FitsAny(p3, nodeInfos)
	assert.Error(t, err)
	predicateErr := predicateChecker.CheckPredicates(p2, nil, ni1)
	assert.NotNil(t, predicateErr)
	assert.True(t, strings.Contains(predicateErr.Error(), "Predicates failed"))
	assert.True(t, strings.Contains(predicateErr.VerboseError(), "Insufficient cpu"))
	assert.NotNil(t, predicateChecker.CheckPredicates(p2, nil, ni1))
	assert.Nil(t, predicateChecker.CheckPredicates(p4, nil, ni1))
	assert.Nil(t, predicateChecker.CheckPredicates(p2, nil, ni2))
	assert.Nil(t, predicateChecker.CheckPredicates(p4, nil, ni2))
	assert.NotNil(t, predicateChecker.CheckPredicates(p3, nil, ni2))
}
