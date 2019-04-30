package model

import (
	"sort"
	"time"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

type vpaConditionsMap map[vpa_types.VerticalPodAutoscalerConditionType]vpa_types.VerticalPodAutoscalerCondition

func (conditionsMap *vpaConditionsMap) Set(conditionType vpa_types.VerticalPodAutoscalerConditionType, status bool, reason string, message string) *vpaConditionsMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	oldCondition, alreadyPresent := (*conditionsMap)[conditionType]
	condition := vpa_types.VerticalPodAutoscalerCondition{Type: conditionType, Reason: reason, Message: message}
	if status {
		condition.Status = apiv1.ConditionTrue
	} else {
		condition.Status = apiv1.ConditionFalse
	}
	if alreadyPresent && oldCondition.Status == condition.Status {
		condition.LastTransitionTime = oldCondition.LastTransitionTime
	} else {
		condition.LastTransitionTime = metav1.Now()
	}
	(*conditionsMap)[conditionType] = condition
	return conditionsMap
}
func (conditionsMap *vpaConditionsMap) AsList() []vpa_types.VerticalPodAutoscalerCondition {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conditions := make([]vpa_types.VerticalPodAutoscalerCondition, 0, len(*conditionsMap))
	for _, condition := range *conditionsMap {
		conditions = append(conditions, condition)
	}
	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].Type < conditions[j].Type
	})
	return conditions
}

type Vpa struct {
	ID				VpaID
	PodSelector			labels.Selector
	Conditions			vpaConditionsMap
	Recommendation			*vpa_types.RecommendedPodResources
	aggregateContainerStates	aggregateContainerStatesMap
	ResourcePolicy			*vpa_types.PodResourcePolicy
	ContainersInitialAggregateState	ContainerNameToAggregateStateMap
	UpdateMode			*vpa_types.UpdateMode
	Created				time.Time
	CheckpointWritten		time.Time
}

func NewVpa(id VpaID, selector labels.Selector, created time.Time) *Vpa {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vpa := &Vpa{ID: id, PodSelector: selector, aggregateContainerStates: make(aggregateContainerStatesMap), ContainersInitialAggregateState: make(ContainerNameToAggregateStateMap), Created: created, Conditions: make(vpaConditionsMap)}
	return vpa
}
func (vpa *Vpa) UseAggregationIfMatching(aggregationKey AggregateStateKey, aggregation *AggregateContainerState) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !vpa.UsesAggregation(aggregationKey) && vpa.matchesAggregation(aggregationKey) {
		vpa.aggregateContainerStates[aggregationKey] = aggregation
	}
}
func (vpa *Vpa) UsesAggregation(aggregationKey AggregateStateKey) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, exists := vpa.aggregateContainerStates[aggregationKey]
	return exists
}
func (vpa *Vpa) DeleteAggregation(aggregationKey AggregateStateKey) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	delete(vpa.aggregateContainerStates, aggregationKey)
}
func (vpa *Vpa) MergeCheckpointedState(aggregateContainerStateMap ContainerNameToAggregateStateMap) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for containerName, aggregation := range vpa.ContainersInitialAggregateState {
		aggregateContainerState, found := aggregateContainerStateMap[containerName]
		if !found {
			aggregateContainerState = NewAggregateContainerState()
			aggregateContainerStateMap[containerName] = aggregateContainerState
		}
		aggregateContainerState.MergeContainerState(aggregation)
	}
}
func (vpa *Vpa) AggregateStateByContainerName() ContainerNameToAggregateStateMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	containerNameToAggregateStateMap := AggregateStateByContainerName(vpa.aggregateContainerStates)
	vpa.MergeCheckpointedState(containerNameToAggregateStateMap)
	return containerNameToAggregateStateMap
}
func (vpa *Vpa) HasRecommendation() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return (vpa.Recommendation != nil) && len(vpa.Recommendation.ContainerRecommendations) > 0
}
func (vpa *Vpa) matchesAggregation(aggregationKey AggregateStateKey) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if vpa.ID.Namespace != aggregationKey.Namespace() {
		return false
	}
	return vpa.PodSelector != nil && vpa.PodSelector.Matches(aggregationKey.Labels())
}
