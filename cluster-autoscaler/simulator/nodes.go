package simulator

import (
	"time"
	apiv1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/autoscaler/cluster-autoscaler/utils/drain"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	kube_client "k8s.io/client-go/kubernetes"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

func GetRequiredPodsForNode(nodename string, client kube_client.Interface) ([]*apiv1.Pod, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	podListResult, err := client.CoreV1().Pods(apiv1.NamespaceAll).List(metav1.ListOptions{FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": nodename}).String()})
	if err != nil {
		return []*apiv1.Pod{}, errors.ToAutoscalerError(errors.ApiCallError, err)
	}
	allPods := make([]*apiv1.Pod, 0)
	for i := range podListResult.Items {
		allPods = append(allPods, &podListResult.Items[i])
	}
	podsToRemoveList, err := drain.GetPodsForDeletionOnNodeDrain(allPods, []*policyv1.PodDisruptionBudget{}, true, false, false, false, nil, 0, time.Now())
	if err != nil {
		return []*apiv1.Pod{}, errors.ToAutoscalerError(errors.InternalError, err)
	}
	podsToRemoveMap := make(map[string]struct{})
	for _, pod := range podsToRemoveList {
		podsToRemoveMap[pod.SelfLink] = struct{}{}
	}
	podsOnNewNode := make([]*apiv1.Pod, 0)
	for _, pod := range allPods {
		if pod.DeletionTimestamp != nil {
			continue
		}
		if _, found := podsToRemoveMap[pod.SelfLink]; !found {
			podsOnNewNode = append(podsOnNewNode, pod)
		}
	}
	return podsOnNewNode, nil
}
func BuildNodeInfoForNode(node *apiv1.Node, client kube_client.Interface) (*schedulercache.NodeInfo, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	requiredPods, err := GetRequiredPodsForNode(node.Name, client)
	if err != nil {
		return nil, err
	}
	result := schedulercache.NewNodeInfo(requiredPods...)
	if err := result.SetNode(node); err != nil {
		return nil, errors.ToAutoscalerError(errors.InternalError, err)
	}
	return result, nil
}
