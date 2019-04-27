package glogx

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"k8s.io/klog"
)

func TestLoggingQuota(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := NewLoggingQuota(3)
	for i := 0; i < 5; i++ {
		assert.Equal(t, 3-i, q.Left())
		assert.Equal(t, i < 3, bool(UpTo(q)))
		assert.Equal(t, i >= 3, bool(Over(q)))
	}
}
func TestReset(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	q := NewLoggingQuota(3)
	for i := 0; i < 5; i++ {
		assert.Equal(t, i < 3, bool(UpTo(q)))
	}
	q.Reset()
	assert.Equal(t, 3, q.Left())
	assert.False(t, bool(Over(q)))
	assert.True(t, bool(UpTo(q)))
}
func TestVFalse(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := Verbose(false)
	q := NewLoggingQuota(3)
	assert.False(t, bool(v.UpTo(q)))
	assert.Equal(t, 3, q.Left())
}
func TestV(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := klog.Level(0); i <= 10; i++ {
		assert.Equal(t, bool(klog.V(i)), bool(V(i)))
	}
}
