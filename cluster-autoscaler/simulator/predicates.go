package simulator

import (
 "fmt"
 "strings"
 apiv1 "k8s.io/api/core/v1"
 kube_util "k8s.io/autoscaler/cluster-autoscaler/utils/kubernetes"
 informers "k8s.io/client-go/informers"
 kube_client "k8s.io/client-go/kubernetes"
 "k8s.io/kubernetes/pkg/scheduler/algorithm"
 "k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
 "k8s.io/kubernetes/pkg/scheduler/factory"
 "k8s.io/kubernetes/pkg/scheduler/algorithmprovider"
 "k8s.io/klog"
)

const (
 affinityPredicateName = "MatchInterPodAffinity"
)

type predicateInfo struct {
 name      string
 predicate algorithm.FitPredicate
}
type PredicateChecker struct {
 predicates                []predicateInfo
 predicateMetadataProducer algorithm.PredicateMetadataProducer
 enableAffinityPredicate   bool
}

var priorityPredicates = []string{"PodFitsResources", "GeneralPredicates", "PodToleratesNodeTaints"}

func init() {
 _logClusterCodePath()
 defer _logClusterCodePath()
 algorithmprovider.ApplyFeatureGates()
}
func NewPredicateChecker(kubeClient kube_client.Interface, stop <-chan struct{}) (*PredicateChecker, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 provider, err := factory.GetAlgorithmProvider(factory.DefaultProvider)
 if err != nil {
  return nil, err
 }
 informerFactory := informers.NewSharedInformerFactory(kubeClient, 0)
 schedulerConfigFactory := factory.NewConfigFactory(&factory.ConfigFactoryArgs{SchedulerName: "cluster-autoscaler", Client: kubeClient, NodeInformer: informerFactory.Core().V1().Nodes(), PodInformer: informerFactory.Core().V1().Pods(), PvInformer: informerFactory.Core().V1().PersistentVolumes(), PvcInformer: informerFactory.Core().V1().PersistentVolumeClaims(), ReplicationControllerInformer: informerFactory.Core().V1().ReplicationControllers(), ReplicaSetInformer: informerFactory.Apps().V1().ReplicaSets(), StatefulSetInformer: informerFactory.Apps().V1().StatefulSets(), ServiceInformer: informerFactory.Core().V1().Services(), PdbInformer: informerFactory.Policy().V1beta1().PodDisruptionBudgets(), StorageClassInformer: informerFactory.Storage().V1().StorageClasses(), HardPodAffinitySymmetricWeight: apiv1.DefaultHardPodAffinitySymmetricWeight})
 informerFactory.Start(stop)
 metadataProducer, err := schedulerConfigFactory.GetPredicateMetadataProducer()
 if err != nil {
  return nil, err
 }
 predicateMap, err := schedulerConfigFactory.GetPredicates(provider.FitPredicateKeys)
 predicateMap["ready"] = isNodeReadyAndSchedulablePredicate
 if err != nil {
  return nil, err
 }
 if _, found := predicateMap["PodFitsResources"]; !found {
  predicateMap["PodFitsResources"] = predicates.PodFitsResources
 }
 predicateList := make([]predicateInfo, 0)
 for _, predicateName := range priorityPredicates {
  if predicate, found := predicateMap[predicateName]; found {
   predicateList = append(predicateList, predicateInfo{name: predicateName, predicate: predicate})
   delete(predicateMap, predicateName)
  }
 }
 for predicateName, predicate := range predicateMap {
  predicateList = append(predicateList, predicateInfo{name: predicateName, predicate: predicate})
 }
 for _, predInfo := range predicateList {
  klog.V(1).Infof("Using predicate %s", predInfo.name)
 }
 return &PredicateChecker{predicates: predicateList, predicateMetadataProducer: metadataProducer, enableAffinityPredicate: true}, nil
}
func isNodeReadyAndSchedulablePredicate(pod *apiv1.Pod, meta algorithm.PredicateMetadata, nodeInfo *schedulercache.NodeInfo) (bool, []algorithm.PredicateFailureReason, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 ready := kube_util.IsNodeReadyAndSchedulable(nodeInfo.Node())
 if !ready {
  return false, []algorithm.PredicateFailureReason{predicates.NewFailureReason("node is unready")}, nil
 }
 return true, []algorithm.PredicateFailureReason{}, nil
}
func NewTestPredicateChecker() *PredicateChecker {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &PredicateChecker{predicates: []predicateInfo{{name: "default", predicate: predicates.GeneralPredicates}, {name: "ready", predicate: isNodeReadyAndSchedulablePredicate}}, predicateMetadataProducer: func(_ *apiv1.Pod, _ map[string]*schedulercache.NodeInfo) algorithm.PredicateMetadata {
  return nil
 }}
}
func (p *PredicateChecker) SetAffinityPredicateEnabled(enable bool) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 p.enableAffinityPredicate = enable
}
func (p *PredicateChecker) IsAffinityPredicateEnabled() bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return p.enableAffinityPredicate
}
func (p *PredicateChecker) GetPredicateMetadata(pod *apiv1.Pod, nodeInfos map[string]*schedulercache.NodeInfo) algorithm.PredicateMetadata {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if !p.enableAffinityPredicate {
  return nil
 }
 return p.predicateMetadataProducer(pod, nodeInfos)
}
func (p *PredicateChecker) FitsAny(pod *apiv1.Pod, nodeInfos map[string]*schedulercache.NodeInfo) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for name, nodeInfo := range nodeInfos {
  if nodeInfo.Node().Spec.Unschedulable {
   continue
  }
  if err := p.CheckPredicates(pod, nil, nodeInfo); err == nil {
   return name, nil
  }
 }
 return "", fmt.Errorf("cannot put pod %s on any node", pod.Name)
}

type PredicateError struct {
 predicateName  string
 failureReasons []algorithm.PredicateFailureReason
 err            error
 reasons        []string
 message        string
}

func (pe *PredicateError) Error() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pe.message != "" {
  return pe.message
 }
 return "Predicates failed"
}
func (pe *PredicateError) VerboseError() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pe.message != "" {
  return pe.message
 }
 if pe.err != nil {
  pe.message = fmt.Sprintf("%s predicate error: %v", pe.predicateName, pe.err)
  return pe.message
 }
 pe.message = fmt.Sprintf("%s predicate mismatch, reason: %s", pe.predicateName, strings.Join(pe.Reasons(), ", "))
 return pe.message
}
func NewPredicateError(name string, err error, reasons []string, originalReasons []algorithm.PredicateFailureReason) *PredicateError {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &PredicateError{predicateName: name, err: err, reasons: reasons, failureReasons: originalReasons}
}
func (pe *PredicateError) Reasons() []string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if pe.reasons != nil {
  return pe.reasons
 }
 pe.reasons = make([]string, len(pe.failureReasons))
 for i, reason := range pe.failureReasons {
  pe.reasons[i] = reason.GetReason()
 }
 return pe.reasons
}
func (pe *PredicateError) OriginalReasons() []algorithm.PredicateFailureReason {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pe.failureReasons
}
func (pe *PredicateError) PredicateName() string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pe.predicateName
}
func (p *PredicateChecker) CheckPredicates(pod *apiv1.Pod, predicateMetadata algorithm.PredicateMetadata, nodeInfo *schedulercache.NodeInfo) *PredicateError {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, predInfo := range p.predicates {
  if !p.enableAffinityPredicate && predInfo.name == affinityPredicateName {
   continue
  }
  match, failureReasons, err := predInfo.predicate(pod, predicateMetadata, nodeInfo)
  if err != nil || !match {
   return &PredicateError{predicateName: predInfo.name, failureReasons: failureReasons, err: err}
  }
 }
 return nil
}
