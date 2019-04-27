package alicloud

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEcsInstanceIdFromProviderId(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	instanceName := "cn-hangzhou.instanceId"
	instanceId, err := ecsInstanceIdFromProviderId(instanceName)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, instanceId, "instanceId")
}
