package backoff

import (
	"testing"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	testprovider "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/test"
	"github.com/stretchr/testify/assert"
)

func nodeGroup(id string) cloudprovider.NodeGroup {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	provider := testprovider.NewTestCloudProvider(nil, nil)
	provider.AddNodeGroup(id, 1, 10, 1)
	return provider.GetNodeGroup(id)
}

var nodeGroup1 = nodeGroup("id1")
var nodeGroup2 = nodeGroup("id2")

func TestBackoffTwoKeys(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoff := NewIdBasedExponentialBackoff(10*time.Minute, time.Hour, 3*time.Hour)
	startTime := time.Now()
	assert.False(t, backoff.IsBackedOff(nodeGroup1, startTime))
	assert.False(t, backoff.IsBackedOff(nodeGroup2, startTime))
	backoff.Backoff(nodeGroup1, startTime.Add(time.Minute))
	assert.True(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(2*time.Minute)))
	assert.False(t, backoff.IsBackedOff(nodeGroup2, startTime))
	assert.False(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(11*time.Minute)))
}
func TestMaxBackoff(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoff := NewIdBasedExponentialBackoff(1*time.Minute, 3*time.Minute, 3*time.Hour)
	startTime := time.Now()
	backoff.Backoff(nodeGroup1, startTime)
	assert.True(t, backoff.IsBackedOff(nodeGroup1, startTime))
	assert.False(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(1*time.Minute)))
	backoff.Backoff(nodeGroup1, startTime.Add(1*time.Minute))
	assert.True(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(1*time.Minute)))
	assert.False(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(3*time.Minute)))
	backoff.Backoff(nodeGroup1, startTime.Add(3*time.Minute))
	assert.True(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(3*time.Minute)))
	assert.False(t, backoff.IsBackedOff(nodeGroup1, startTime.Add(6*time.Minute)))
}
func TestRemoveBackoff(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoff := NewIdBasedExponentialBackoff(1*time.Minute, 3*time.Minute, 3*time.Hour)
	startTime := time.Now()
	backoff.Backoff(nodeGroup1, startTime)
	assert.True(t, backoff.IsBackedOff(nodeGroup1, startTime))
	backoff.RemoveBackoff(nodeGroup1)
	assert.False(t, backoff.IsBackedOff(nodeGroup1, startTime))
}
func TestResetStaleBackoffData(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	backoff := NewIdBasedExponentialBackoff(1*time.Minute, 3*time.Minute, 3*time.Hour)
	startTime := time.Now()
	backoff.Backoff(nodeGroup1, startTime)
	backoff.Backoff(nodeGroup2, startTime.Add(time.Hour))
	backoff.RemoveStaleBackoffData(startTime.Add(time.Hour))
	assert.Equal(t, 2, len(backoff.(*exponentialBackoff).backoffInfo))
	backoff.RemoveStaleBackoffData(startTime.Add(4 * time.Hour))
	assert.Equal(t, 1, len(backoff.(*exponentialBackoff).backoffInfo))
	backoff.RemoveStaleBackoffData(startTime.Add(5 * time.Hour))
	assert.Equal(t, 0, len(backoff.(*exponentialBackoff).backoffInfo))
}
