package gke

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetNormalizedGpuCount(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	gpus, err := getNormalizedGpuCount(int64(0))
	assert.Equal(t, err, nil)
	assert.Equal(t, gpus, int64(1))
	gpus, err = getNormalizedGpuCount(int64(1))
	assert.Equal(t, err, nil)
	assert.Equal(t, gpus, int64(1))
	gpus, err = getNormalizedGpuCount(int64(2))
	assert.Equal(t, err, nil)
	assert.Equal(t, gpus, int64(2))
	gpus, err = getNormalizedGpuCount(int64(3))
	assert.Equal(t, err, nil)
	assert.Equal(t, gpus, int64(4))
	gpus, err = getNormalizedGpuCount(int64(7))
	assert.Equal(t, err, nil)
	assert.Equal(t, gpus, int64(8))
	gpus, err = getNormalizedGpuCount(int64(8))
	assert.Equal(t, err, nil)
	assert.Equal(t, gpus, int64(8))
	gpus, err = getNormalizedGpuCount(int64(9))
	assert.Error(t, err)
}
func TestValidateGpuConfig(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := validateGpuConfig("nvidia-tesla-k80", int64(1), "europe-west1-b", "n1-standard-1")
	assert.Equal(t, err, nil)
	err = validateGpuConfig("nvidia-tesla-p100", int64(1), "europe-west1-b", "n1-standard-1")
	assert.Equal(t, err, nil)
	err = validateGpuConfig("nvidia-tesla-k80", int64(4), "europe-west1-b", "n1-standard-32")
	assert.Equal(t, err, nil)
	err = validateGpuConfig("duke-igthorn", int64(1), "europe-west1-b", "n1-standard-1")
	assert.Error(t, err)
	err = validateGpuConfig("nvidia-tesla-k80", int64(1), "castle-drekmore", "n1-standard-1")
	assert.Error(t, err)
	err = validateGpuConfig("nvidia-tesla-k80", int64(1), "europe-west1-b", "toadie-the-ogre")
	assert.Error(t, err)
	err = validateGpuConfig("nvidia-tesla-k80", int64(1), "europe-west1-b", "n1-standard-32")
	assert.Error(t, err)
}
