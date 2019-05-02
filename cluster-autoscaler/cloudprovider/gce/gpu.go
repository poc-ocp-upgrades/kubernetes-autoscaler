package gce

import (
 "strconv"
 "strings"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

var (
 gpuZones          = map[string]map[string]bool{"nvidia-tesla-k80": {"us-west1-b": true, "us-central1-c": true, "us-east1-c": true, "us-east1-d": true, "europe-west1-b": true, "europe-west1-d": true, "asia-east1-a": true, "asia-east1-b": true}, "nvidia-tesla-p100": {"us-west1-b": true, "us-central1-c": true, "us-central1-f": true, "us-east1-b": true, "us-east1-c": true, "europe-west1-b": true, "europe-west1-d": true, "asia-east1-a": true, "asia-east1-c": true, "europe-west4-a": true}, "nvidia-tesla-v100": {"us-west1-a": true, "us-west1-b": true, "us-central1-a": true, "us-central1-f": true, "europe-west4-a": true, "asia-east1-c": true}}
 maxGpuCount       = map[string]int64{"nvidia-tesla-k80": 8, "nvidia-tesla-p100": 4, "nvidia-tesla-v100": 8}
 maxCpuCount       = map[string]map[int64]int{"nvidia-tesla-k80": {1: 8, 2: 16, 4: 32, 8: 64}, "nvidia-tesla-p100": {1: 16, 2: 32, 4: 64}, "nvidia-tesla-v100": {1: 12, 2: 24, 4: 48, 8: 96}}
 supportedGpuTypes = []string{"nvidia-tesla-k80", "nvidia-tesla-p100", "nvidia-tesla-v100"}
)

func validateGpuConfig(gpuType string, gpuCount int64, zone, machineType string) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 zoneInfo, found := gpuZones[gpuType]
 if !found {
  return cloudprovider.ErrIllegalConfiguration
 }
 if allowed := zoneInfo[zone]; !allowed {
  return cloudprovider.ErrIllegalConfiguration
 }
 maxGpu, found := maxGpuCount[gpuType]
 if !found || gpuCount > maxGpu {
  return cloudprovider.ErrIllegalConfiguration
 }
 parts := strings.Split(machineType, "-")
 cpus, err := strconv.Atoi(parts[len(parts)-1])
 if err != nil {
  return cloudprovider.ErrIllegalConfiguration
 }
 maxCpuInfo, found := maxCpuCount[gpuType]
 if !found {
  return cloudprovider.ErrIllegalConfiguration
 }
 maxCpus, found := maxCpuInfo[gpuCount]
 if !found || cpus > maxCpus {
  return cloudprovider.ErrIllegalConfiguration
 }
 return nil
}
func getNormalizedGpuCount(initialCount int64) (int64, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for i := int64(1); i <= int64(8); i = 2 * i {
  if i >= initialCount {
   return i, nil
  }
 }
 return 0, cloudprovider.ErrIllegalConfiguration
}
