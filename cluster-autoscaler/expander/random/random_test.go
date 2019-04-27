package random

import (
	"testing"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	"github.com/stretchr/testify/assert"
)

func TestRandomExpander(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	eo1a := expander.Option{Debug: "EO1a"}
	e := NewStrategy()
	ret := e.BestOption([]expander.Option{eo1a}, nil)
	assert.Equal(t, *ret, eo1a)
	eo1b := expander.Option{Debug: "EO1b"}
	ret = e.BestOption([]expander.Option{eo1a, eo1b}, nil)
	assert.True(t, assert.ObjectsAreEqual(*ret, eo1a) || assert.ObjectsAreEqual(*ret, eo1b))
}
