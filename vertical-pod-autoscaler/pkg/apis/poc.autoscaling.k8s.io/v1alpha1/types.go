package v1alpha1

import (
 "k8s.io/api/core/v1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VerticalPodAutoscalerList struct {
 metav1.TypeMeta `json:",inline"`
 metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
 Items           []VerticalPodAutoscaler `json:"items" protobuf:"bytes,2,rep,name=items"`
}
type VerticalPodAutoscaler struct {
 metav1.TypeMeta   `json:",inline"`
 metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
 Spec              VerticalPodAutoscalerSpec   `json:"spec" protobuf:"bytes,2,name=spec"`
 Status            VerticalPodAutoscalerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
type VerticalPodAutoscalerSpec struct {
 Selector       *metav1.LabelSelector `json:"selector" protobuf:"bytes,1,name=selector"`
 UpdatePolicy   *PodUpdatePolicy      `json:"updatePolicy,omitempty" protobuf:"bytes,2,opt,name=updatePolicy"`
 ResourcePolicy *PodResourcePolicy    `json:"resourcePolicy,omitempty" protobuf:"bytes,3,opt,name=resourcePolicy"`
}
type PodUpdatePolicy struct {
 UpdateMode *UpdateMode `json:"updateMode,omitempty" protobuf:"bytes,1,opt,name=updateMode"`
}
type UpdateMode string

const (
 UpdateModeOff      UpdateMode = "Off"
 UpdateModeInitial  UpdateMode = "Initial"
 UpdateModeRecreate UpdateMode = "Recreate"
 UpdateModeAuto     UpdateMode = "Auto"
)

type PodResourcePolicy struct {
 ContainerPolicies []ContainerResourcePolicy `json:"containerPolicies,omitempty" patchStrategy:"merge" patchMergeKey:"containerName" protobuf:"bytes,1,rep,name=containerPolicies"`
}
type ContainerResourcePolicy struct {
 ContainerName string                `json:"containerName,omitempty" protobuf:"bytes,1,opt,name=containerName"`
 Mode          *ContainerScalingMode `json:"mode,omitempty" protobuf:"bytes,2,opt,name=mode"`
 MinAllowed    v1.ResourceList       `json:"minAllowed,omitempty" protobuf:"bytes,3,rep,name=minAllowed,casttype=ResourceList,castkey=ResourceName"`
 MaxAllowed    v1.ResourceList       `json:"maxAllowed,omitempty" protobuf:"bytes,4,rep,name=maxAllowed,casttype=ResourceList,castkey=ResourceName"`
}

const (
 DefaultContainerResourcePolicy = "*"
)

type ContainerScalingMode string

const (
 ContainerScalingModeAuto ContainerScalingMode = "Auto"
 ContainerScalingModeOff  ContainerScalingMode = "Off"
)

type VerticalPodAutoscalerStatus struct {
 Recommendation *RecommendedPodResources         `json:"recommendation,omitempty" protobuf:"bytes,1,opt,name=recommendation"`
 Conditions     []VerticalPodAutoscalerCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}
type RecommendedPodResources struct {
 ContainerRecommendations []RecommendedContainerResources `json:"containerRecommendations,omitempty" protobuf:"bytes,1,rep,name=containerRecommendations"`
}
type RecommendedContainerResources struct {
 ContainerName string          `json:"containerName,omitempty" protobuf:"bytes,1,opt,name=containerName"`
 Target        v1.ResourceList `json:"target" protobuf:"bytes,2,rep,name=target,casttype=ResourceList,castkey=ResourceName"`
 LowerBound    v1.ResourceList `json:"lowerBound,omitempty" protobuf:"bytes,3,rep,name=lowerBound,casttype=ResourceList,castkey=ResourceName"`
 UpperBound    v1.ResourceList `json:"upperBound,omitempty" protobuf:"bytes,4,rep,name=upperBound,casttype=ResourceList,castkey=ResourceName"`
}
type VerticalPodAutoscalerConditionType string

var (
 RecommendationProvided VerticalPodAutoscalerConditionType = "RecommendationProvided"
 LowConfidence          VerticalPodAutoscalerConditionType = "LowConfidence"
 NoPodsMatched          VerticalPodAutoscalerConditionType = "NoPodsMatched"
 FetchingHistory        VerticalPodAutoscalerConditionType = "FetchingHistory"
)

type VerticalPodAutoscalerCondition struct {
 Type               VerticalPodAutoscalerConditionType `json:"type" protobuf:"bytes,1,name=type"`
 Status             v1.ConditionStatus                 `json:"status" protobuf:"bytes,2,name=status"`
 LastTransitionTime metav1.Time                        `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
 Reason             string                             `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
 Message            string                             `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}
type VerticalPodAutoscalerCheckpoint struct {
 metav1.TypeMeta   `json:",inline"`
 metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
 Spec              VerticalPodAutoscalerCheckpointSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
 Status            VerticalPodAutoscalerCheckpointStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
type VerticalPodAutoscalerCheckpointList struct {
 metav1.TypeMeta `json:",inline"`
 metav1.ListMeta `json:"metadata"`
 Items           []VerticalPodAutoscalerCheckpoint `json:"items"`
}
type VerticalPodAutoscalerCheckpointSpec struct {
 VPAObjectName string `json:"vpaObjectName,omitempty" protobuf:"bytes,1,opt,name=vpaObjectName"`
 ContainerName string `json:"containerName,omitempty" protobuf:"bytes,2,opt,name=containerName"`
}
type VerticalPodAutoscalerCheckpointStatus struct {
 LastUpdateTime    metav1.Time         `json:"lastUpdateTime,omitempty" protobuf:"bytes,1,opt,name=lastUpdateTime"`
 Version           string              `json:"version,omitempty" protobuf:"bytes,2,opt,name=version"`
 CPUHistogram      HistogramCheckpoint `json:"cpuHistogram,omitempty" protobuf:"bytes,3,rep,name=cpuHistograms"`
 MemoryHistogram   HistogramCheckpoint `json:"memoryHistogram,omitempty" protobuf:"bytes,4,rep,name=memoryHistogram"`
 FirstSampleStart  metav1.Time         `json:"firstSampleStart,omitempty" protobuf:"bytes,5,opt,name=firstSampleStart"`
 LastSampleStart   metav1.Time         `json:"lastSampleStart,omitempty" protobuf:"bytes,6,opt,name=lastSampleStart"`
 TotalSamplesCount int                 `json:"totalSamplesCount,omitempty" protobuf:"bytes,7,opt,name=totalSamplesCount"`
}
type HistogramCheckpoint struct {
 ReferenceTimestamp metav1.Time    `json:"referenceTimestamp,omitempty" protobuf:"bytes,1,opt,name=referenceTimestamp"`
 BucketWeights      map[int]uint32 `json:"bucketWeights,omitempty" protobuf:"bytes,2,opt,name=bucketWeights"`
 TotalWeight        float64        `json:"totalWeight,omitempty" protobuf:"bytes,3,opt,name=totalWeight"`
}
