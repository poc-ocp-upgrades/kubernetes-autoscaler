package gce

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateKernelReserved(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		physicalMemory	int64
		reservedMemory	int64
	}
	testCases := []testCase{{physicalMemory: 256 * MiB, reservedMemory: 4*MiB + kernelReservedMemory}, {physicalMemory: 2 * GiB, reservedMemory: 32*MiB + kernelReservedMemory}, {physicalMemory: 3 * GiB, reservedMemory: 48*MiB + kernelReservedMemory}, {physicalMemory: 3.25 * GiB, reservedMemory: 52*MiB + kernelReservedMemory + swiotlbReservedMemory}, {physicalMemory: 4 * GiB, reservedMemory: 64*MiB + kernelReservedMemory + swiotlbReservedMemory}, {physicalMemory: 128 * GiB, reservedMemory: 2*GiB + kernelReservedMemory + swiotlbReservedMemory}}
	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%v", idx), func(t *testing.T) {
			reserved := CalculateKernelReserved(tc.physicalMemory)
			assert.Equal(t, tc.reservedMemory, reserved)
		})
	}
}
