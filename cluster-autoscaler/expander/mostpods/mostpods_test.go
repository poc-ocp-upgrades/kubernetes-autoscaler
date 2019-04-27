package mostpods

import (
	"testing"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
)

func TestMostPods(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	eo0 := expander.Option{Debug: "EO0"}
	e := NewStrategy()
	ret := e.BestOption([]expander.Option{eo0}, nil)
	assert.Equal(t, *ret, eo0)
	eo1 := expander.Option{Debug: "EO1", Pods: []*apiv1.Pod{nil}}
	ret = e.BestOption([]expander.Option{eo0, eo1}, nil)
	assert.Equal(t, *ret, eo1)
	eo1b := expander.Option{Debug: "EO1b", Pods: []*apiv1.Pod{nil}}
	ret = e.BestOption([]expander.Option{eo0, eo1, eo1b}, nil)
	assert.NotEqual(t, *ret, eo0)
	assert.True(t, assert.ObjectsAreEqual(*ret, eo1) || assert.ObjectsAreEqual(*ret, eo1b))
}
