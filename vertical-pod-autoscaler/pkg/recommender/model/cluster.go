package model

import (
 "fmt"
 "time"
 "github.com/golang/glog"
 apiv1 "k8s.io/api/core/v1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 labels "k8s.io/apimachinery/pkg/labels"
 vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta1"
)

type ClusterState struct {
 Pods              map[PodID]*PodState
 Vpas              map[VpaID]*Vpa
 ObservedVpas      []*vpa_types.VerticalPodAutoscaler
 aggregateStateMap aggregateContainerStatesMap
 labelSetMap       labelSetMap
}
type AggregateStateKey interface {
 Namespace() string
 ContainerName() string
 Labels() labels.Labels
}
type labelSetKey string
type labelSetMap map[labelSetKey]labels.Set
type aggregateContainerStatesMap map[AggregateStateKey]*AggregateContainerState
type PodState struct {
 ID          PodID
 labelSetKey labelSetKey
 Containers  map[string]*ContainerState
 Phase       apiv1.PodPhase
}

func NewClusterState() *ClusterState {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &ClusterState{Pods: make(map[PodID]*PodState), Vpas: make(map[VpaID]*Vpa), aggregateStateMap: make(aggregateContainerStatesMap), labelSetMap: make(labelSetMap)}
}

type ContainerUsageSampleWithKey struct {
 ContainerUsageSample
 Container ContainerID
}

func (cluster *ClusterState) AddOrUpdatePod(podID PodID, newLabels labels.Set, phase apiv1.PodPhase) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, podExists := cluster.Pods[podID]
 if !podExists {
  pod = newPod(podID)
  cluster.Pods[podID] = pod
 }
 newlabelSetKey := cluster.getLabelSetKey(newLabels)
 if !podExists || pod.labelSetKey != newlabelSetKey {
  pod.labelSetKey = newlabelSetKey
  for containerName, container := range pod.Containers {
   containerID := ContainerID{PodID: podID, ContainerName: containerName}
   container.aggregator = cluster.findOrCreateAggregateContainerState(containerID)
  }
 }
 pod.Phase = phase
}
func (cluster *ClusterState) GetContainer(containerID ContainerID) *ContainerState {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, podExists := cluster.Pods[containerID.PodID]
 if podExists {
  container, containerExists := pod.Containers[containerID.ContainerName]
  if containerExists {
   return container
  }
 }
 return nil
}
func (cluster *ClusterState) DeletePod(podID PodID) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 delete(cluster.Pods, podID)
}
func (cluster *ClusterState) AddOrUpdateContainer(containerID ContainerID, request Resources) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, podExists := cluster.Pods[containerID.PodID]
 if !podExists {
  return NewKeyError(containerID.PodID)
 }
 if container, containerExists := pod.Containers[containerID.ContainerName]; !containerExists {
  cluster.findOrCreateAggregateContainerState(containerID)
  pod.Containers[containerID.ContainerName] = NewContainerState(request, NewContainerStateAggregatorProxy(cluster, containerID))
 } else {
  container.Request = request
 }
 return nil
}
func (cluster *ClusterState) AddSample(sample *ContainerUsageSampleWithKey) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, podExists := cluster.Pods[sample.Container.PodID]
 if !podExists {
  return NewKeyError(sample.Container.PodID)
 }
 containerState, containerExists := pod.Containers[sample.Container.ContainerName]
 if !containerExists {
  return NewKeyError(sample.Container)
 }
 if !containerState.AddSample(&sample.ContainerUsageSample) {
  return fmt.Errorf("Sample discarded (invalid or out of order)")
 }
 return nil
}
func (cluster *ClusterState) RecordOOM(containerID ContainerID, timestamp time.Time, requestedMemory ResourceAmount) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, podExists := cluster.Pods[containerID.PodID]
 if !podExists {
  return NewKeyError(containerID.PodID)
 }
 containerState, containerExists := pod.Containers[containerID.ContainerName]
 if !containerExists {
  return NewKeyError(containerID.ContainerName)
 }
 err := containerState.RecordOOM(timestamp, requestedMemory)
 if err != nil {
  return fmt.Errorf("Error while recording OOM for %v, Reason: %v", containerID, err)
 }
 return nil
}
func (cluster *ClusterState) AddOrUpdateVpa(apiObject *vpa_types.VerticalPodAutoscaler) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 vpaID := VpaID{Namespace: apiObject.Namespace, VpaName: apiObject.Name}
 conditionsMap := make(vpaConditionsMap)
 for _, condition := range apiObject.Status.Conditions {
  conditionsMap[condition.Type] = condition
 }
 var currentRecommendation *vpa_types.RecommendedPodResources
 if conditionsMap[vpa_types.RecommendationProvided].Status == apiv1.ConditionTrue {
  currentRecommendation = apiObject.Status.Recommendation
 }
 selector, err := metav1.LabelSelectorAsSelector(apiObject.Spec.Selector)
 if err != nil {
  return err
 }
 vpa, vpaExists := cluster.Vpas[vpaID]
 if vpaExists && (err != nil || vpa.PodSelector.String() != selector.String()) {
  if err := cluster.DeleteVpa(vpaID); err != nil {
   return err
  }
  vpaExists = false
 }
 if !vpaExists {
  vpa = NewVpa(vpaID, selector, apiObject.CreationTimestamp.Time)
  cluster.Vpas[vpaID] = vpa
  for aggregationKey, aggregation := range cluster.aggregateStateMap {
   vpa.UseAggregationIfMatching(aggregationKey, aggregation)
  }
 }
 vpa.Conditions = conditionsMap
 vpa.Recommendation = currentRecommendation
 vpa.ResourcePolicy = apiObject.Spec.ResourcePolicy
 if apiObject.Spec.UpdatePolicy != nil {
  vpa.UpdateMode = apiObject.Spec.UpdatePolicy.UpdateMode
 }
 return nil
}
func (cluster *ClusterState) DeleteVpa(vpaID VpaID) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if _, vpaExists := cluster.Vpas[vpaID]; !vpaExists {
  return NewKeyError(vpaID)
 }
 delete(cluster.Vpas, vpaID)
 return nil
}
func newPod(id PodID) *PodState {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &PodState{ID: id, Containers: make(map[string]*ContainerState)}
}
func (cluster *ClusterState) getLabelSetKey(labelSet labels.Set) labelSetKey {
 _logClusterCodePath()
 defer _logClusterCodePath()
 labelSetKey := labelSetKey(labelSet.String())
 cluster.labelSetMap[labelSetKey] = labelSet
 return labelSetKey
}
func (cluster *ClusterState) MakeAggregateStateKey(pod *PodState, containerName string) AggregateStateKey {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return aggregateStateKey{namespace: pod.ID.Namespace, containerName: containerName, labelSetKey: pod.labelSetKey, labelSetMap: &cluster.labelSetMap}
}
func (cluster *ClusterState) aggregateStateKeyForContainerID(containerID ContainerID) AggregateStateKey {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pod, podExists := cluster.Pods[containerID.PodID]
 if !podExists {
  panic(fmt.Sprintf("Pod not present in the ClusterState: %v", containerID.PodID))
 }
 return cluster.MakeAggregateStateKey(pod, containerID.ContainerName)
}
func (cluster *ClusterState) findOrCreateAggregateContainerState(containerID ContainerID) *AggregateContainerState {
 _logClusterCodePath()
 defer _logClusterCodePath()
 aggregateStateKey := cluster.aggregateStateKeyForContainerID(containerID)
 aggregateContainerState, aggregateStateExists := cluster.aggregateStateMap[aggregateStateKey]
 if !aggregateStateExists {
  aggregateContainerState = NewAggregateContainerState()
  cluster.aggregateStateMap[aggregateStateKey] = aggregateContainerState
  for _, vpa := range cluster.Vpas {
   vpa.UseAggregationIfMatching(aggregateStateKey, aggregateContainerState)
  }
 }
 return aggregateContainerState
}
func (cluster *ClusterState) GarbageCollectAggregateCollectionStates(now time.Time) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 glog.V(1).Info("Garbage collection of AggregateCollectionStates triggered")
 keysToDelete := make([]AggregateStateKey, 0)
 for key, aggregateContainerState := range cluster.aggregateStateMap {
  if aggregateContainerState.isExpired(now) {
   keysToDelete = append(keysToDelete, key)
   glog.V(1).Infof("Removing AggregateCollectionStates for %+v", key)
  }
 }
 for _, key := range keysToDelete {
  delete(cluster.aggregateStateMap, key)
  for _, vpa := range cluster.Vpas {
   vpa.DeleteAggregation(key)
  }
 }
}

type aggregateStateKey struct {
 namespace     string
 containerName string
 labelSetKey   labelSetKey
 labelSetMap   *labelSetMap
}

func (k aggregateStateKey) Namespace() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return k.namespace
}
func (k aggregateStateKey) ContainerName() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return k.containerName
}
func (k aggregateStateKey) Labels() labels.Labels {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return (*k.labelSetMap)[k.labelSetKey]
}
