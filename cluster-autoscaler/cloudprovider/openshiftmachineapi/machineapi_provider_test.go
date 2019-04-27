package openshiftmachineapi

import (
	"reflect"
	"testing"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

func TestProviderConstructorProperties(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceLimits := cloudprovider.ResourceLimiter{}
	controller, stop := mustCreateTestController(t)
	defer stop()
	provider, err := newProvider(ProviderName, &resourceLimits, controller)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if actual := provider.Name(); actual != ProviderName {
		t.Errorf("expected %q, got %q", ProviderName, actual)
	}
	rl, err := provider.GetResourceLimiter()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reflect.DeepEqual(rl, resourceLimits) {
		t.Errorf("expected %+v, got %+v", resourceLimits, rl)
	}
	if _, err := provider.Pricing(); err != cloudprovider.ErrNotImplemented {
		t.Errorf("expected an error")
	}
	machineTypes, err := provider.GetAvailableMachineTypes()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(machineTypes) != 0 {
		t.Errorf("expected 0, got %v", len(machineTypes))
	}
	if _, err := provider.NewNodeGroup("foo", nil, nil, nil, nil); err == nil {
		t.Error("expected an error")
	}
	if err := provider.Cleanup(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := provider.Refresh(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	nodegroups := provider.NodeGroups()
	if len(nodegroups) != 0 {
		t.Errorf("expected 0, got %v", len(nodegroups))
	}
	ng, err := provider.NodeGroupForNode(&apiv1.Node{TypeMeta: v1.TypeMeta{Kind: "Node"}, ObjectMeta: v1.ObjectMeta{Name: "missing-node"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ng != nil {
		t.Fatalf("unexpected nodegroup: %v", ng.Id())
	}
}
