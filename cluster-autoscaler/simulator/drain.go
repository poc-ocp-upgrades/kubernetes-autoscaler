package simulator

import (
 "fmt"
 "time"
 apiv1 "k8s.io/api/core/v1"
 policyv1 "k8s.io/api/policy/v1beta1"
 metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 "k8s.io/apimachinery/pkg/labels"
 "k8s.io/autoscaler/cluster-autoscaler/utils/drain"
 client "k8s.io/client-go/kubernetes"
 schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

func FastGetPodsToMove(nodeInfo *schedulercache.NodeInfo, skipNodesWithSystemPods bool, skipNodesWithLocalStorage bool, pdbs []*policyv1.PodDisruptionBudget) ([]*apiv1.Pod, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pods, err := drain.GetPodsForDeletionOnNodeDrain(nodeInfo.Pods(), pdbs, false, skipNodesWithSystemPods, skipNodesWithLocalStorage, false, nil, 0, time.Now())
 if err != nil {
  return pods, err
 }
 if err := checkPdbs(pods, pdbs); err != nil {
  return []*apiv1.Pod{}, err
 }
 return pods, nil
}
func DetailedGetPodsForMove(nodeInfo *schedulercache.NodeInfo, skipNodesWithSystemPods bool, skipNodesWithLocalStorage bool, client client.Interface, minReplicaCount int32, pdbs []*policyv1.PodDisruptionBudget) ([]*apiv1.Pod, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 pods, err := drain.GetPodsForDeletionOnNodeDrain(nodeInfo.Pods(), pdbs, false, skipNodesWithSystemPods, skipNodesWithLocalStorage, true, client, minReplicaCount, time.Now())
 if err != nil {
  return pods, err
 }
 if err := checkPdbs(pods, pdbs); err != nil {
  return []*apiv1.Pod{}, err
 }
 return pods, nil
}
func checkPdbs(pods []*apiv1.Pod, pdbs []*policyv1.PodDisruptionBudget) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, pdb := range pdbs {
  selector, err := metav1.LabelSelectorAsSelector(pdb.Spec.Selector)
  if err != nil {
   return err
  }
  for _, pod := range pods {
   if pod.Namespace == pdb.Namespace && selector.Matches(labels.Set(pod.Labels)) {
    if pdb.Status.PodDisruptionsAllowed < 1 {
     return fmt.Errorf("not enough pod disruption budget to move %s/%s", pod.Namespace, pod.Name)
    }
   }
  }
 }
 return nil
}
