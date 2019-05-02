package gce

import (
 "fmt"
 "math/rand"
 "regexp"
 "strings"
 gce "google.golang.org/api/compute/v1"
 apiv1 "k8s.io/api/core/v1"
 "k8s.io/apimachinery/pkg/api/resource"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
 "k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
 kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
 "github.com/ghodss/yaml"
 "k8s.io/klog"
)

type GceTemplateBuilder struct{}

func (t *GceTemplateBuilder) getAcceleratorCount(accelerators []*gce.AcceleratorConfig) int64 {
 _logClusterCodePath()
 defer _logClusterCodePath()
 count := int64(0)
 for _, accelerator := range accelerators {
  if strings.HasPrefix(accelerator.AcceleratorType, "nvidia-") {
   count += accelerator.AcceleratorCount
  }
 }
 return count
}
func (t *GceTemplateBuilder) BuildCapacity(cpu int64, mem int64, accelerators []*gce.AcceleratorConfig) (apiv1.ResourceList, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 capacity := apiv1.ResourceList{}
 capacity[apiv1.ResourcePods] = *resource.NewQuantity(110, resource.DecimalSI)
 capacity[apiv1.ResourceCPU] = *resource.NewQuantity(cpu, resource.DecimalSI)
 memTotal := mem - CalculateKernelReserved(mem)
 capacity[apiv1.ResourceMemory] = *resource.NewQuantity(memTotal, resource.DecimalSI)
 if accelerators != nil && len(accelerators) > 0 {
  capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(t.getAcceleratorCount(accelerators), resource.DecimalSI)
 }
 return capacity, nil
}
func (t *GceTemplateBuilder) BuildAllocatableFromKubeEnv(capacity apiv1.ResourceList, kubeEnv string) (apiv1.ResourceList, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 kubeReserved, err := extractKubeReservedFromKubeEnv(kubeEnv)
 if err != nil {
  return nil, err
 }
 reserved, err := parseKubeReserved(kubeReserved)
 if err != nil {
  return nil, err
 }
 return t.CalculateAllocatable(capacity, reserved), nil
}
func (t *GceTemplateBuilder) CalculateAllocatable(capacity, kubeReserved apiv1.ResourceList) apiv1.ResourceList {
 _logClusterCodePath()
 defer _logClusterCodePath()
 allocatable := apiv1.ResourceList{}
 for key, value := range capacity {
  quantity := *value.Copy()
  if reservedQuantity, found := kubeReserved[key]; found {
   quantity.Sub(reservedQuantity)
  }
  if key == apiv1.ResourceMemory {
   quantity = *resource.NewQuantity(quantity.Value()-KubeletEvictionHardMemory, resource.BinarySI)
  }
  allocatable[key] = quantity
 }
 return allocatable
}
func (t *GceTemplateBuilder) BuildNodeFromTemplate(mig Mig, template *gce.InstanceTemplate, cpu int64, mem int64) (*apiv1.Node, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if template.Properties == nil {
  return nil, fmt.Errorf("instance template %s has no properties", template.Name)
 }
 node := apiv1.Node{}
 nodeName := fmt.Sprintf("%s-template-%d", template.Name, rand.Int63())
 node.ObjectMeta = metav1.ObjectMeta{Name: nodeName, SelfLink: fmt.Sprintf("/api/v1/nodes/%s", nodeName), Labels: map[string]string{}}
 capacity, err := t.BuildCapacity(cpu, mem, template.Properties.GuestAccelerators)
 if err != nil {
  return nil, err
 }
 node.Status = apiv1.NodeStatus{Capacity: capacity}
 var nodeAllocatable apiv1.ResourceList
 if template.Properties.Metadata == nil {
  return nil, fmt.Errorf("instance template %s has no metadata", template.Name)
 }
 for _, item := range template.Properties.Metadata.Items {
  if item.Key == "kube-env" {
   if item.Value == nil {
    return nil, fmt.Errorf("no kube-env content in metadata")
   }
   kubeEnvLabels, err := extractLabelsFromKubeEnv(*item.Value)
   if err != nil {
    return nil, err
   }
   node.Labels = cloudprovider.JoinStringMaps(node.Labels, kubeEnvLabels)
   kubeEnvTaints, err := extractTaintsFromKubeEnv(*item.Value)
   if err != nil {
    return nil, err
   }
   node.Spec.Taints = append(node.Spec.Taints, kubeEnvTaints...)
   if allocatable, err := t.BuildAllocatableFromKubeEnv(node.Status.Capacity, *item.Value); err == nil {
    nodeAllocatable = allocatable
   }
  }
 }
 if nodeAllocatable == nil {
  klog.Warningf("could not extract kube-reserved from kubeEnv for mig %q, setting allocatable to capacity.", mig.GceRef().Name)
  node.Status.Allocatable = node.Status.Capacity
 } else {
  node.Status.Allocatable = nodeAllocatable
 }
 labels, err := BuildGenericLabels(mig.GceRef(), template.Properties.MachineType, nodeName)
 if err != nil {
  return nil, err
 }
 node.Labels = cloudprovider.JoinStringMaps(node.Labels, labels)
 node.Status.Conditions = cloudprovider.BuildReadyConditions()
 return &node, nil
}
func BuildGenericLabels(ref GceRef, machineType string, nodeName string) (map[string]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := make(map[string]string)
 result[kubeletapis.LabelArch] = cloudprovider.DefaultArch
 result[kubeletapis.LabelOS] = cloudprovider.DefaultOS
 result[kubeletapis.LabelInstanceType] = machineType
 ix := strings.LastIndex(ref.Zone, "-")
 if ix == -1 {
  return nil, fmt.Errorf("unexpected zone: %s", ref.Zone)
 }
 result[kubeletapis.LabelZoneRegion] = ref.Zone[:ix]
 result[kubeletapis.LabelZoneFailureDomain] = ref.Zone
 result[kubeletapis.LabelHostname] = nodeName
 return result, nil
}
func parseKubeReserved(kubeReserved string) (apiv1.ResourceList, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 resourcesMap, err := parseKeyValueListToMap(kubeReserved)
 if err != nil {
  return nil, fmt.Errorf("failed to extract kube-reserved from kube-env: %q", err)
 }
 reservedResources := apiv1.ResourceList{}
 for name, quantity := range resourcesMap {
  switch apiv1.ResourceName(name) {
  case apiv1.ResourceCPU, apiv1.ResourceMemory, apiv1.ResourceEphemeralStorage:
   if q, err := resource.ParseQuantity(quantity); err == nil && q.Sign() >= 0 {
    reservedResources[apiv1.ResourceName(name)] = q
   }
  default:
   klog.Warningf("ignoring resource from kube-reserved: %q", name)
  }
 }
 return reservedResources, nil
}
func extractLabelsFromKubeEnv(kubeEnv string) (map[string]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 labels, err := extractAutoscalerVarFromKubeEnv(kubeEnv, "node_labels")
 if err != nil {
  klog.Errorf("node_labels not found via AUTOSCALER_ENV_VARS due to error, will try NODE_LABELS: %v", err)
  labels, err = extractFromKubeEnv(kubeEnv, "NODE_LABELS")
  if err != nil {
   return nil, err
  }
 }
 return parseKeyValueListToMap(labels)
}
func extractTaintsFromKubeEnv(kubeEnv string) ([]apiv1.Taint, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 taints, err := extractAutoscalerVarFromKubeEnv(kubeEnv, "node_taints")
 if err != nil {
  klog.Errorf("node_taints not found via AUTOSCALER_ENV_VARS due to error, will try NODE_TAINTS: %v", err)
  taints, err = extractFromKubeEnv(kubeEnv, "NODE_TAINTS")
  if err != nil {
   return nil, err
  }
 }
 taintMap, err := parseKeyValueListToMap(taints)
 if err != nil {
  return nil, err
 }
 return buildTaints(taintMap)
}
func extractKubeReservedFromKubeEnv(kubeEnv string) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 kubeReserved, err := extractAutoscalerVarFromKubeEnv(kubeEnv, "kube_reserved")
 if err != nil {
  klog.Errorf("kube_reserved not found via AUTOSCALER_ENV_VARS due to error, will try kube-reserved in KUBELET_TEST_ARGS: %v", err)
  kubeletArgs, err := extractFromKubeEnv(kubeEnv, "KUBELET_TEST_ARGS")
  if err != nil {
   return "", err
  }
  resourcesRegexp := regexp.MustCompile(`--kube-reserved=([^ ]+)`)
  matches := resourcesRegexp.FindStringSubmatch(kubeletArgs)
  if len(matches) > 1 {
   return matches[1], nil
  }
  return "", fmt.Errorf("kube-reserved not in kubelet args in kube-env: %q", kubeletArgs)
 }
 return kubeReserved, nil
}
func extractAutoscalerVarFromKubeEnv(kubeEnv, name string) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 const autoscalerVars = "AUTOSCALER_ENV_VARS"
 autoscalerVals, err := extractFromKubeEnv(kubeEnv, autoscalerVars)
 if err != nil {
  return "", err
 }
 for _, val := range strings.Split(autoscalerVals, ";") {
  val = strings.Trim(val, " ")
  items := strings.SplitN(val, "=", 2)
  if len(items) != 2 {
   return "", fmt.Errorf("malformed autoscaler var: %s", val)
  }
  if strings.Trim(items[0], " ") == name {
   return strings.Trim(items[1], " \"'"), nil
  }
 }
 return "", fmt.Errorf("var %s not found in %s: %v", name, autoscalerVars, autoscalerVals)
}
func extractFromKubeEnv(kubeEnv, resource string) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 kubeEnvMap := make(map[string]string)
 err := yaml.Unmarshal([]byte(kubeEnv), &kubeEnvMap)
 if err != nil {
  return "", fmt.Errorf("Error unmarshalling kubeEnv: %v", err)
 }
 return kubeEnvMap[resource], nil
}
func parseKeyValueListToMap(kvList string) (map[string]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 result := make(map[string]string)
 if len(kvList) == 0 {
  return result, nil
 }
 for _, keyValue := range strings.Split(kvList, ",") {
  kvItems := strings.SplitN(keyValue, "=", 2)
  if len(kvItems) != 2 {
   return nil, fmt.Errorf("error while parsing key-value list, val: %s", keyValue)
  }
  result[kvItems[0]] = kvItems[1]
 }
 return result, nil
}
func buildTaints(kubeEnvTaints map[string]string) ([]apiv1.Taint, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 taints := make([]apiv1.Taint, 0)
 for key, value := range kubeEnvTaints {
  values := strings.SplitN(value, ":", 2)
  if len(values) != 2 {
   return nil, fmt.Errorf("error while parsing node taint value and effect: %s", value)
  }
  taints = append(taints, apiv1.Taint{Key: key, Value: values[0], Effect: apiv1.TaintEffect(values[1])})
 }
 return taints, nil
}
