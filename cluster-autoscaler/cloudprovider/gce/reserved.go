package gce

const (
	MiB				= 1024 * 1024
	GiB				= 1024 * 1024 * 1024
	KubeletEvictionHardMemory	= 100 * MiB
	kernelReservedRatio		= 64
	kernelReservedMemory		= 16 * MiB
	swiotlbReservedMemory		= 64 * MiB
	swiotlbThresholdMemory		= 3 * GiB
)

func CalculateKernelReserved(physicalMemory int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	reserved := int64(physicalMemory / kernelReservedRatio)
	reserved += kernelReservedMemory
	if physicalMemory > swiotlbThresholdMemory {
		reserved += swiotlbReservedMemory
	}
	return reserved
}
