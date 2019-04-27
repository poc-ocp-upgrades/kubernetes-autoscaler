package gke

import (
	"testing"
)

func TestPredictKubeReserved(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		name			string
		function		func(capacity int64) int64
		capacity		int64
		expectedReserved	int64
	}
	testCases := []testCase{{name: "zero memory capacity", function: PredictKubeReservedMemory, capacity: 0, expectedReserved: 0}, {name: "f1-micro", function: PredictKubeReservedMemory, capacity: 600 * MiB, expectedReserved: 255 * MiB}, {name: "between memory thresholds", function: PredictKubeReservedMemory, capacity: 2000 * MiB, expectedReserved: 500 * MiB}, {name: "at a memory threshold boundary", function: PredictKubeReservedMemory, capacity: 8000 * MiB, expectedReserved: 1800 * MiB}, {name: "exceeds highest memory threshold", function: PredictKubeReservedMemory, capacity: 200 * 1000 * MiB, expectedReserved: 10760 * MiB}, {name: "cpu sanity check", function: PredictKubeReservedCpuMillicores, capacity: 4000, expectedReserved: 80}}
	for _, tc := range testCases {
		if actualReserved := tc.function(tc.capacity); actualReserved != tc.expectedReserved {
			t.Errorf("Test case: %s, Got f(%d Mb) = %d.  Want %d", tc.name, tc.capacity, actualReserved, tc.expectedReserved)
		}
	}
}
func TestCalculateReserved(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		name			string
		function		func(capacity int64) int64
		capacity		int64
		expectedReserved	int64
	}
	testCases := []testCase{{name: "zero memory capacity", function: memoryReservedMiB, capacity: 0, expectedReserved: 0}, {name: "f1-micro", function: memoryReservedMiB, capacity: 600, expectedReserved: 255}, {name: "between memory thresholds", function: memoryReservedMiB, capacity: 2000, expectedReserved: 500}, {name: "at a memory threshold boundary", function: memoryReservedMiB, capacity: 8000, expectedReserved: 1800}, {name: "exceeds highest memory threshold", function: memoryReservedMiB, capacity: 200 * 1000, expectedReserved: 10760}, {name: "cpu sanity check", function: cpuReservedMillicores, capacity: 4 * millicoresPerCore, expectedReserved: 80}}
	for _, tc := range testCases {
		if actualReserved := tc.function(tc.capacity); actualReserved != tc.expectedReserved {
			t.Errorf("Test case: %s, Got f(%d Mb) = %d.  Want %d", tc.name, tc.capacity, actualReserved, tc.expectedReserved)
		}
	}
}
