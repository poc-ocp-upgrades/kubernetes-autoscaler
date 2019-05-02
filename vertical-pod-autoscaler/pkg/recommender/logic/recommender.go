package logic

import (
 "flag"
 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
)

var (
 safetyMarginFraction = flag.Float64("recommendation-margin-fraction", 0.15, `Fraction of usage added as the safety margin to the recommended request`)
 podMinCPUMillicores  = flag.Float64("pod-recommendation-min-cpu-millicores", 25, `Minimum CPU recommendation for a pod`)
 podMinMemoryMb       = flag.Float64("pod-recommendation-min-memory-mb", 250, `Minimum memory recommendation for a pod`)
)

type PodResourceRecommender interface {
 GetRecommendedPodResources(containerNameToAggregateStateMap model.ContainerNameToAggregateStateMap) RecommendedPodResources
}
type RecommendedPodResources map[string]RecommendedContainerResources
type RecommendedContainerResources struct {
 Target     model.Resources
 LowerBound model.Resources
 UpperBound model.Resources
}
type podResourceRecommender struct {
 targetEstimator     ResourceEstimator
 lowerBoundEstimator ResourceEstimator
 upperBoundEstimator ResourceEstimator
}

func (r *podResourceRecommender) GetRecommendedPodResources(containerNameToAggregateStateMap model.ContainerNameToAggregateStateMap) RecommendedPodResources {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var recommendation = make(RecommendedPodResources)
 if len(containerNameToAggregateStateMap) == 0 {
  return recommendation
 }
 fraction := 1.0 / float64(len(containerNameToAggregateStateMap))
 minResources := model.Resources{model.ResourceCPU: model.ScaleResource(model.CPUAmountFromCores(*podMinCPUMillicores*0.001), fraction), model.ResourceMemory: model.ScaleResource(model.MemoryAmountFromBytes(*podMinMemoryMb*1024*1024), fraction)}
 recommender := &podResourceRecommender{WithMinResources(minResources, r.targetEstimator), WithMinResources(minResources, r.lowerBoundEstimator), WithMinResources(minResources, r.upperBoundEstimator)}
 for containerName, aggregatedContainerState := range containerNameToAggregateStateMap {
  recommendation[containerName] = recommender.estimateContainerResources(aggregatedContainerState)
 }
 return recommendation
}
func (r *podResourceRecommender) estimateContainerResources(s *model.AggregateContainerState) RecommendedContainerResources {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return RecommendedContainerResources{r.targetEstimator.GetResourceEstimation(s), r.lowerBoundEstimator.GetResourceEstimation(s), r.upperBoundEstimator.GetResourceEstimation(s)}
}
func CreatePodResourceRecommender() PodResourceRecommender {
 _logClusterCodePath()
 defer _logClusterCodePath()
 targetCPUPercentile := 0.9
 lowerBoundCPUPercentile := 0.5
 upperBoundCPUPercentile := 0.95
 targetMemoryPeaksPercentile := 0.9
 lowerBoundMemoryPeaksPercentile := 0.5
 upperBoundMemoryPeaksPercentile := 0.95
 targetEstimator := NewPercentileEstimator(targetCPUPercentile, targetMemoryPeaksPercentile)
 lowerBoundEstimator := NewPercentileEstimator(lowerBoundCPUPercentile, lowerBoundMemoryPeaksPercentile)
 upperBoundEstimator := NewPercentileEstimator(upperBoundCPUPercentile, upperBoundMemoryPeaksPercentile)
 targetEstimator = WithMargin(*safetyMarginFraction, targetEstimator)
 lowerBoundEstimator = WithMargin(*safetyMarginFraction, lowerBoundEstimator)
 upperBoundEstimator = WithMargin(*safetyMarginFraction, upperBoundEstimator)
 upperBoundEstimator = WithConfidenceMultiplier(1.0, 1.0, upperBoundEstimator)
 lowerBoundEstimator = WithConfidenceMultiplier(0.001, -2.0, lowerBoundEstimator)
 return &podResourceRecommender{targetEstimator, lowerBoundEstimator, upperBoundEstimator}
}
