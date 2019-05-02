package v1beta1

import (
 v1 "k8s.io/api/core/v1"
 meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *ContainerResourcePolicy) DeepCopyInto(out *ContainerResourcePolicy) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.Mode != nil {
  in, out := &in.Mode, &out.Mode
  if *in == nil {
   *out = nil
  } else {
   *out = new(ContainerScalingMode)
   **out = **in
  }
 }
 if in.MinAllowed != nil {
  in, out := &in.MinAllowed, &out.MinAllowed
  *out = make(v1.ResourceList, len(*in))
  for key, val := range *in {
   (*out)[key] = val.DeepCopy()
  }
 }
 if in.MaxAllowed != nil {
  in, out := &in.MaxAllowed, &out.MaxAllowed
  *out = make(v1.ResourceList, len(*in))
  for key, val := range *in {
   (*out)[key] = val.DeepCopy()
  }
 }
 return
}
func (in *ContainerResourcePolicy) DeepCopy() *ContainerResourcePolicy {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(ContainerResourcePolicy)
 in.DeepCopyInto(out)
 return out
}
func (in *HistogramCheckpoint) DeepCopyInto(out *HistogramCheckpoint) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 in.ReferenceTimestamp.DeepCopyInto(&out.ReferenceTimestamp)
 if in.BucketWeights != nil {
  in, out := &in.BucketWeights, &out.BucketWeights
  *out = make(map[int]uint32, len(*in))
  for key, val := range *in {
   (*out)[key] = val
  }
 }
 return
}
func (in *HistogramCheckpoint) DeepCopy() *HistogramCheckpoint {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(HistogramCheckpoint)
 in.DeepCopyInto(out)
 return out
}
func (in *PodResourcePolicy) DeepCopyInto(out *PodResourcePolicy) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.ContainerPolicies != nil {
  in, out := &in.ContainerPolicies, &out.ContainerPolicies
  *out = make([]ContainerResourcePolicy, len(*in))
  for i := range *in {
   (*in)[i].DeepCopyInto(&(*out)[i])
  }
 }
 return
}
func (in *PodResourcePolicy) DeepCopy() *PodResourcePolicy {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(PodResourcePolicy)
 in.DeepCopyInto(out)
 return out
}
func (in *PodUpdatePolicy) DeepCopyInto(out *PodUpdatePolicy) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.UpdateMode != nil {
  in, out := &in.UpdateMode, &out.UpdateMode
  if *in == nil {
   *out = nil
  } else {
   *out = new(UpdateMode)
   **out = **in
  }
 }
 return
}
func (in *PodUpdatePolicy) DeepCopy() *PodUpdatePolicy {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(PodUpdatePolicy)
 in.DeepCopyInto(out)
 return out
}
func (in *RecommendedContainerResources) DeepCopyInto(out *RecommendedContainerResources) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.Target != nil {
  in, out := &in.Target, &out.Target
  *out = make(v1.ResourceList, len(*in))
  for key, val := range *in {
   (*out)[key] = val.DeepCopy()
  }
 }
 if in.LowerBound != nil {
  in, out := &in.LowerBound, &out.LowerBound
  *out = make(v1.ResourceList, len(*in))
  for key, val := range *in {
   (*out)[key] = val.DeepCopy()
  }
 }
 if in.UpperBound != nil {
  in, out := &in.UpperBound, &out.UpperBound
  *out = make(v1.ResourceList, len(*in))
  for key, val := range *in {
   (*out)[key] = val.DeepCopy()
  }
 }
 if in.UncappedTarget != nil {
  in, out := &in.UncappedTarget, &out.UncappedTarget
  *out = make(v1.ResourceList, len(*in))
  for key, val := range *in {
   (*out)[key] = val.DeepCopy()
  }
 }
 return
}
func (in *RecommendedContainerResources) DeepCopy() *RecommendedContainerResources {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(RecommendedContainerResources)
 in.DeepCopyInto(out)
 return out
}
func (in *RecommendedPodResources) DeepCopyInto(out *RecommendedPodResources) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.ContainerRecommendations != nil {
  in, out := &in.ContainerRecommendations, &out.ContainerRecommendations
  *out = make([]RecommendedContainerResources, len(*in))
  for i := range *in {
   (*in)[i].DeepCopyInto(&(*out)[i])
  }
 }
 return
}
func (in *RecommendedPodResources) DeepCopy() *RecommendedPodResources {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(RecommendedPodResources)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscaler) DeepCopyInto(out *VerticalPodAutoscaler) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 out.TypeMeta = in.TypeMeta
 in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
 in.Spec.DeepCopyInto(&out.Spec)
 in.Status.DeepCopyInto(&out.Status)
 return
}
func (in *VerticalPodAutoscaler) DeepCopy() *VerticalPodAutoscaler {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscaler)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscaler) DeepCopyObject() runtime.Object {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if c := in.DeepCopy(); c != nil {
  return c
 }
 return nil
}
func (in *VerticalPodAutoscalerCheckpoint) DeepCopyInto(out *VerticalPodAutoscalerCheckpoint) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 out.TypeMeta = in.TypeMeta
 in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
 out.Spec = in.Spec
 in.Status.DeepCopyInto(&out.Status)
 return
}
func (in *VerticalPodAutoscalerCheckpoint) DeepCopy() *VerticalPodAutoscalerCheckpoint {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerCheckpoint)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerCheckpoint) DeepCopyObject() runtime.Object {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if c := in.DeepCopy(); c != nil {
  return c
 }
 return nil
}
func (in *VerticalPodAutoscalerCheckpointList) DeepCopyInto(out *VerticalPodAutoscalerCheckpointList) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 out.TypeMeta = in.TypeMeta
 out.ListMeta = in.ListMeta
 if in.Items != nil {
  in, out := &in.Items, &out.Items
  *out = make([]VerticalPodAutoscalerCheckpoint, len(*in))
  for i := range *in {
   (*in)[i].DeepCopyInto(&(*out)[i])
  }
 }
 return
}
func (in *VerticalPodAutoscalerCheckpointList) DeepCopy() *VerticalPodAutoscalerCheckpointList {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerCheckpointList)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerCheckpointList) DeepCopyObject() runtime.Object {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if c := in.DeepCopy(); c != nil {
  return c
 }
 return nil
}
func (in *VerticalPodAutoscalerCheckpointSpec) DeepCopyInto(out *VerticalPodAutoscalerCheckpointSpec) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 return
}
func (in *VerticalPodAutoscalerCheckpointSpec) DeepCopy() *VerticalPodAutoscalerCheckpointSpec {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerCheckpointSpec)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerCheckpointStatus) DeepCopyInto(out *VerticalPodAutoscalerCheckpointStatus) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
 in.CPUHistogram.DeepCopyInto(&out.CPUHistogram)
 in.MemoryHistogram.DeepCopyInto(&out.MemoryHistogram)
 in.FirstSampleStart.DeepCopyInto(&out.FirstSampleStart)
 in.LastSampleStart.DeepCopyInto(&out.LastSampleStart)
 return
}
func (in *VerticalPodAutoscalerCheckpointStatus) DeepCopy() *VerticalPodAutoscalerCheckpointStatus {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerCheckpointStatus)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerCondition) DeepCopyInto(out *VerticalPodAutoscalerCondition) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
 return
}
func (in *VerticalPodAutoscalerCondition) DeepCopy() *VerticalPodAutoscalerCondition {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerCondition)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerList) DeepCopyInto(out *VerticalPodAutoscalerList) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 out.TypeMeta = in.TypeMeta
 out.ListMeta = in.ListMeta
 if in.Items != nil {
  in, out := &in.Items, &out.Items
  *out = make([]VerticalPodAutoscaler, len(*in))
  for i := range *in {
   (*in)[i].DeepCopyInto(&(*out)[i])
  }
 }
 return
}
func (in *VerticalPodAutoscalerList) DeepCopy() *VerticalPodAutoscalerList {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerList)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerList) DeepCopyObject() runtime.Object {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if c := in.DeepCopy(); c != nil {
  return c
 }
 return nil
}
func (in *VerticalPodAutoscalerSpec) DeepCopyInto(out *VerticalPodAutoscalerSpec) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.Selector != nil {
  in, out := &in.Selector, &out.Selector
  if *in == nil {
   *out = nil
  } else {
   *out = new(meta_v1.LabelSelector)
   (*in).DeepCopyInto(*out)
  }
 }
 if in.UpdatePolicy != nil {
  in, out := &in.UpdatePolicy, &out.UpdatePolicy
  if *in == nil {
   *out = nil
  } else {
   *out = new(PodUpdatePolicy)
   (*in).DeepCopyInto(*out)
  }
 }
 if in.ResourcePolicy != nil {
  in, out := &in.ResourcePolicy, &out.ResourcePolicy
  if *in == nil {
   *out = nil
  } else {
   *out = new(PodResourcePolicy)
   (*in).DeepCopyInto(*out)
  }
 }
 return
}
func (in *VerticalPodAutoscalerSpec) DeepCopy() *VerticalPodAutoscalerSpec {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerSpec)
 in.DeepCopyInto(out)
 return out
}
func (in *VerticalPodAutoscalerStatus) DeepCopyInto(out *VerticalPodAutoscalerStatus) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 *out = *in
 if in.Recommendation != nil {
  in, out := &in.Recommendation, &out.Recommendation
  if *in == nil {
   *out = nil
  } else {
   *out = new(RecommendedPodResources)
   (*in).DeepCopyInto(*out)
  }
 }
 if in.Conditions != nil {
  in, out := &in.Conditions, &out.Conditions
  *out = make([]VerticalPodAutoscalerCondition, len(*in))
  for i := range *in {
   (*in)[i].DeepCopyInto(&(*out)[i])
  }
 }
 return
}
func (in *VerticalPodAutoscalerStatus) DeepCopy() *VerticalPodAutoscalerStatus {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if in == nil {
  return nil
 }
 out := new(VerticalPodAutoscalerStatus)
 in.DeepCopyInto(out)
 return out
}
