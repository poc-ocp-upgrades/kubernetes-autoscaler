package price

import (
	"fmt"
	"math"
	"time"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/expander"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	"k8s.io/autoscaler/cluster-autoscaler/utils/units"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"k8s.io/klog"
)

type priceBased struct {
	pricingModel		cloudprovider.PricingModel
	preferredNodeProvider	PreferredNodeProvider
	nodeUnfitness		NodeUnfitness
}

var (
	defaultPreferredNode	= buildNode(4*1000, 4*4*units.GiB)
	priceStabilizationPod	= buildPod("stabilize", 500, 500*units.MiB)
	notExistCoeficient	= 2.0
	gpuUnfitnessOverride	= 1000.0
)

func NewStrategy(pricingModel cloudprovider.PricingModel, preferredNodeProvider PreferredNodeProvider, nodeUnfitness NodeUnfitness) expander.Strategy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &priceBased{pricingModel: pricingModel, preferredNodeProvider: preferredNodeProvider, nodeUnfitness: nodeUnfitness}
}
func (p *priceBased) BestOption(expansionOptions []expander.Option, nodeInfos map[string]*schedulercache.NodeInfo) *expander.Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var bestOption *expander.Option
	bestOptionScore := 0.0
	now := time.Now()
	then := now.Add(time.Hour)
	preferredNode, err := p.preferredNodeProvider.Node()
	if err != nil {
		klog.Errorf("Failed to get preferred node, switching to default: %v", err)
		preferredNode = defaultPreferredNode
	}
	stabilizationPrice, err := p.pricingModel.PodPrice(priceStabilizationPod, now, then)
	if err != nil {
		klog.Errorf("Failed to get price for stabilization pod: %v", err)
	}
nextoption:
	for _, option := range expansionOptions {
		nodeInfo, found := nodeInfos[option.NodeGroup.Id()]
		if !found {
			klog.Warningf("No node info for %s", option.NodeGroup.Id())
			continue
		}
		nodePrice, err := p.pricingModel.NodePrice(nodeInfo.Node(), now, then)
		if err != nil {
			klog.Warningf("Failed to calculate node price for %s: %v", option.NodeGroup.Id(), err)
			continue
		}
		totalNodePrice := nodePrice * float64(option.NodeCount)
		totalPodPrice := 0.0
		for _, pod := range option.Pods {
			podPrice, err := p.pricingModel.PodPrice(pod, now, then)
			if err != nil {
				klog.Warningf("Failed to calculate pod price for %s/%s: %v", pod.Namespace, pod.Name, err)
				continue nextoption
			}
			totalPodPrice += podPrice
		}
		priceSubScore := (totalNodePrice + stabilizationPrice) / (totalPodPrice + stabilizationPrice)
		nodeUnfitness := p.nodeUnfitness(preferredNode, nodeInfo.Node())
		supressedUnfitness := (nodeUnfitness-1.0)*(1.0-math.Tanh(float64(option.NodeCount-1)/15.0)) + 1.0
		if gpu.NodeHasGpu(nodeInfo.Node()) {
			klog.V(4).Infof("Price expander overriding unfitness for node group with GPU %s", option.NodeGroup.Id())
			supressedUnfitness = gpuUnfitnessOverride
		}
		optionScore := supressedUnfitness * priceSubScore
		if !option.NodeGroup.Exist() {
			optionScore *= notExistCoeficient
		}
		debug := fmt.Sprintf("all_nodes_price=%f pods_price=%f stabilized_ratio=%f unfitness=%f suppressed=%f final_score=%f", totalNodePrice, totalPodPrice, priceSubScore, nodeUnfitness, supressedUnfitness, optionScore)
		klog.V(5).Infof("Price expander for %s: %s", option.NodeGroup.Id(), debug)
		if bestOption == nil || bestOptionScore > optionScore {
			bestOption = &expander.Option{NodeGroup: option.NodeGroup, NodeCount: option.NodeCount, Debug: fmt.Sprintf("%s | price-expander: %s", option.Debug, debug), Pods: option.Pods}
			bestOptionScore = optionScore
		}
	}
	return bestOption
}
func buildPod(name string, millicpu int64, mem int64) *apiv1.Pod {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &apiv1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: name, SelfLink: fmt.Sprintf("/api/v1/namespaces/default/pods/%s", name)}, Spec: apiv1.PodSpec{Containers: []apiv1.Container{{Resources: apiv1.ResourceRequirements{Requests: apiv1.ResourceList{apiv1.ResourceCPU: *resource.NewMilliQuantity(millicpu, resource.DecimalSI), apiv1.ResourceMemory: *resource.NewQuantity(mem, resource.DecimalSI)}}}}}}
}
