package gke

const (
	MiB			= 1024 * 1024
	mbPerGB			= 1000
	millicoresPerCore	= 1000
)

func PredictKubeReservedMemory(physicalMemory int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return memoryReservedMiB(physicalMemory/MiB) * MiB
}
func PredictKubeReservedCpuMillicores(physicalCpuMillicores int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return cpuReservedMillicores(physicalCpuMillicores)
}

type allocatableBracket struct {
	threshold		int64
	marginalReservedRate	float64
}

func memoryReservedMiB(memoryCapacityMiB int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if memoryCapacityMiB <= mbPerGB {
		if memoryCapacityMiB <= 0 {
			return 0
		}
		return 255
	}
	return calculateReserved(memoryCapacityMiB, []allocatableBracket{{threshold: 0, marginalReservedRate: 0.25}, {threshold: 4 * mbPerGB, marginalReservedRate: 0.2}, {threshold: 8 * mbPerGB, marginalReservedRate: 0.1}, {threshold: 16 * mbPerGB, marginalReservedRate: 0.06}, {threshold: 128 * mbPerGB, marginalReservedRate: 0.02}})
}
func cpuReservedMillicores(cpuCapacityMillicores int64) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return calculateReserved(cpuCapacityMillicores, []allocatableBracket{{threshold: 0, marginalReservedRate: 0.06}, {threshold: 1 * millicoresPerCore, marginalReservedRate: 0.01}, {threshold: 2 * millicoresPerCore, marginalReservedRate: 0.005}, {threshold: 4 * millicoresPerCore, marginalReservedRate: 0.0025}})
}
func calculateReserved(capacity int64, brackets []allocatableBracket) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var reserved float64
	for i, bracket := range brackets {
		c := capacity
		if i < len(brackets)-1 && brackets[i+1].threshold < capacity {
			c = brackets[i+1].threshold
		}
		additionalReserved := float64(c-bracket.threshold) * bracket.marginalReservedRate
		if additionalReserved > 0 {
			reserved += additionalReserved
		}
	}
	return int64(reserved)
}
