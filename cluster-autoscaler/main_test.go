package main

import (
	"testing"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"github.com/stretchr/testify/assert"
)

func TestParseSingleGpuLimit(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testcase struct {
		input			string
		expectError		bool
		expectedLimits		config.GpuLimits
		expectedErrorMessage	string
	}
	testcases := []testcase{{input: "gpu:1:10", expectError: false, expectedLimits: config.GpuLimits{GpuType: "gpu", Min: 1, Max: 10}}, {input: "gpu:1", expectError: true, expectedErrorMessage: "Incorrect gpu limit specification: gpu:1"}, {input: "gpu:1:10:x", expectError: true, expectedErrorMessage: "Incorrect gpu limit specification: gpu:1:10:x"}, {input: "gpu:x:10", expectError: true, expectedErrorMessage: "Incorrect gpu limit - min is not integer: gpu:x:10"}, {input: "gpu:1:y", expectError: true, expectedErrorMessage: "Incorrect gpu limit - max is not integer: gpu:1:y"}, {input: "gpu:-1:10", expectError: true, expectedErrorMessage: "Incorrect gpu limit - min is less than 0; gpu:-1:10"}, {input: "gpu:1:-10", expectError: true, expectedErrorMessage: "Incorrect gpu limit - max is less than 0; gpu:1:-10"}, {input: "gpu:10:1", expectError: true, expectedErrorMessage: "Incorrect gpu limit - min is greater than max; gpu:10:1"}}
	for _, testcase := range testcases {
		limits, err := parseSingleGpuLimit(testcase.input)
		if testcase.expectError {
			assert.NotNil(t, err)
			if err != nil {
				assert.Equal(t, testcase.expectedErrorMessage, err.Error())
			}
		} else {
			assert.Equal(t, testcase.expectedLimits, limits)
		}
	}
}
