package model

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

var (
	anyTime = time.Unix(0, 0)
)

func TestMergeAggregateContainerState(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	containersInitialAggregateState := ContainerNameToAggregateStateMap{}
	containersInitialAggregateState["test"] = NewAggregateContainerState()
	vpa := NewVpa(VpaID{}, nil, anyTime)
	vpa.ContainersInitialAggregateState = containersInitialAggregateState
	containerNameToAggregateStateMap := ContainerNameToAggregateStateMap{}
	vpa.MergeCheckpointedState(containerNameToAggregateStateMap)
	assert.Contains(t, containerNameToAggregateStateMap, "test")
}
