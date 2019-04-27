package aws

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBuildAsg(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	asgCache := &asgCache{}
	asg, err := asgCache.buildAsgFromSpec("1:5:test-asg")
	assert.NoError(t, err)
	assert.Equal(t, asg.minSize, 1)
	assert.Equal(t, asg.maxSize, 5)
	assert.Equal(t, asg.Name, "test-asg")
	_, err = asgCache.buildAsgFromSpec("a")
	assert.Error(t, err)
	_, err = asgCache.buildAsgFromSpec("a:b:c")
	assert.Error(t, err)
	_, err = asgCache.buildAsgFromSpec("1:")
	assert.Error(t, err)
	_, err = asgCache.buildAsgFromSpec("1:2:")
	assert.Error(t, err)
}
func validateAsg(t *testing.T, asg *asg, name string, minSize int, maxSize int) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	assert.Equal(t, name, asg.Name)
	assert.Equal(t, minSize, asg.minSize)
	assert.Equal(t, maxSize, asg.maxSize)
}
