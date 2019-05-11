package routines

import (
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
	"k8s.io/autoscaler/vertical-pod-autoscaler/pkg/recommender/model"
	api_utils "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/utils/vpa"
)

func GetContainerNameToAggregateStateMap(vpa *model.Vpa) model.ContainerNameToAggregateStateMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	containerNameToAggregateStateMap := vpa.AggregateStateByContainerName()
	filteredContainerNameToAggregateStateMap := make(model.ContainerNameToAggregateStateMap)
	for containerName, aggregatedContainerState := range containerNameToAggregateStateMap {
		containerResourcePolicy := api_utils.GetContainerResourcePolicy(containerName, vpa.ResourcePolicy)
		autoscalingDisabled := containerResourcePolicy != nil && containerResourcePolicy.Mode != nil && *containerResourcePolicy.Mode == vpa_types.ContainerScalingModeOff
		if !autoscalingDisabled && aggregatedContainerState.TotalSamplesCount > 0 {
			filteredContainerNameToAggregateStateMap[containerName] = aggregatedContainerState
		}
	}
	return filteredContainerNameToAggregateStateMap
}
