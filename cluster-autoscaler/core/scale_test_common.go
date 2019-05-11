package core

import (
	"testing"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate/utils"
	"k8s.io/autoscaler/cluster-autoscaler/config"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/estimator"
	"k8s.io/autoscaler/cluster-autoscaler/expander/random"
	"k8s.io/autoscaler/cluster-autoscaler/metrics"
	"k8s.io/autoscaler/cluster-autoscaler/processors/nodegroups"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/autoscaler/cluster-autoscaler/utils/labels"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/clusterstate"
	"k8s.io/autoscaler/cluster-autoscaler/utils/backoff"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
)

type nodeConfig struct {
	name	string
	cpu		int64
	memory	int64
	gpu		int64
	ready	bool
	group	string
}
type podConfig struct {
	name	string
	cpu		int64
	memory	int64
	gpu		int64
	node	string
}
type groupSizeChange struct {
	groupName	string
	sizeChange	int
}
type scaleTestConfig struct {
	nodes					[]nodeConfig
	pods					[]podConfig
	extraPods				[]podConfig
	expectedScaleUpOptions	[]groupSizeChange
	scaleUpOptionToChoose	groupSizeChange
	expectedFinalScaleUp	groupSizeChange
	expectedScaleDowns		[]string
	options					config.AutoscalingOptions
}

func NewScaleTestAutoscalingContext(options config.AutoscalingOptions, fakeClient kube_client.Interface, provider cloudprovider.CloudProvider) context.AutoscalingContext {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fakeRecorder := kube_record.NewFakeRecorder(5)
	fakeLogRecorder, _ := utils.NewStatusMapRecorder(fakeClient, "kube-system", fakeRecorder, false)
	estimatorBuilder, _ := estimator.NewEstimatorBuilder(options.EstimatorName)
	return context.AutoscalingContext{AutoscalingOptions: options, AutoscalingKubeClients: context.AutoscalingKubeClients{ClientSet: fakeClient, Recorder: fakeRecorder, LogRecorder: fakeLogRecorder}, CloudProvider: provider, PredicateChecker: simulator.NewTestPredicateChecker(), ExpanderStrategy: random.NewStrategy(), EstimatorBuilder: estimatorBuilder}
}

type mockAutoprovisioningNodeGroupManager struct{ t *testing.T }

func (p *mockAutoprovisioningNodeGroupManager) CreateNodeGroup(context *context.AutoscalingContext, nodeGroup cloudprovider.NodeGroup) (nodegroups.CreateNodeGroupResult, errors.AutoscalerError) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	newNodeGroup, err := nodeGroup.Create()
	assert.NoError(p.t, err)
	metrics.RegisterNodeGroupCreation()
	result := nodegroups.CreateNodeGroupResult{MainCreatedNodeGroup: newNodeGroup}
	return result, nil
}
func (p *mockAutoprovisioningNodeGroupManager) RemoveUnneededNodeGroups(context *context.AutoscalingContext) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !context.AutoscalingOptions.NodeAutoprovisioningEnabled {
		return nil
	}
	nodeGroups := context.CloudProvider.NodeGroups()
	for _, nodeGroup := range nodeGroups {
		if !nodeGroup.Autoprovisioned() {
			continue
		}
		targetSize, err := nodeGroup.TargetSize()
		assert.NoError(p.t, err)
		if targetSize > 0 {
			continue
		}
		nodes, err := nodeGroup.Nodes()
		assert.NoError(p.t, err)
		if len(nodes) > 0 {
			continue
		}
		err = nodeGroup.Delete()
		assert.NoError(p.t, err)
	}
	return nil
}
func (p *mockAutoprovisioningNodeGroupManager) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}

type mockAutoprovisioningNodeGroupListProcessor struct{ t *testing.T }

func (p *mockAutoprovisioningNodeGroupListProcessor) Process(context *context.AutoscalingContext, nodeGroups []cloudprovider.NodeGroup, nodeInfos map[string]*schedulercache.NodeInfo, unschedulablePods []*apiv1.Pod) ([]cloudprovider.NodeGroup, map[string]*schedulercache.NodeInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	machines, err := context.CloudProvider.GetAvailableMachineTypes()
	assert.NoError(p.t, err)
	bestLabels := labels.BestLabelSet(unschedulablePods)
	for _, machineType := range machines {
		nodeGroup, err := context.CloudProvider.NewNodeGroup(machineType, bestLabels, map[string]string{}, []apiv1.Taint{}, map[string]resource.Quantity{})
		assert.NoError(p.t, err)
		nodeInfo, err := nodeGroup.TemplateNodeInfo()
		assert.NoError(p.t, err)
		nodeInfos[nodeGroup.Id()] = nodeInfo
		nodeGroups = append(nodeGroups, nodeGroup)
	}
	return nodeGroups, nodeInfos, nil
}
func (p *mockAutoprovisioningNodeGroupListProcessor) CleanUp() {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func newBackoff() backoff.Backoff {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return backoff.NewIdBasedExponentialBackoff(clusterstate.InitialNodeGroupBackoffDuration, clusterstate.MaxNodeGroupBackoffDuration, clusterstate.NodeGroupBackoffResetTimeout)
}
