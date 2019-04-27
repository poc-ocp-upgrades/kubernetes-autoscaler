package openshiftmachineapi

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"testing"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/utils/pointer"
)

const (
	testNamespace = "test-namespace"
)

func TestNodeGroupNewNodeGroupConstructor(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		description	string
		annotations	map[string]string
		errors		bool
		replicas	int32
		minSize		int
		maxSize		int
		nodeCount	int
	}
	var testCases = []testCase{{description: "errors because minSize is invalid", annotations: map[string]string{nodeGroupMinSizeAnnotationKey: "-1", nodeGroupMaxSizeAnnotationKey: "0"}, errors: true}, {description: "errors because maxSize is invalid", annotations: map[string]string{nodeGroupMinSizeAnnotationKey: "0", nodeGroupMaxSizeAnnotationKey: "-1"}, errors: true}, {description: "errors because minSize > maxSize", annotations: map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "0"}, errors: true}, {description: "errors because maxSize < minSize", annotations: map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "0"}, errors: true}, {description: "no error: min=0, max=0", minSize: 0, maxSize: 0, replicas: 0, errors: false}, {description: "no error: min=0, max=1", annotations: map[string]string{nodeGroupMaxSizeAnnotationKey: "1"}, minSize: 0, maxSize: 1, replicas: 0, errors: false}, {description: "no error: min=1, max=10, replicas=5", annotations: map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}, minSize: 1, maxSize: 10, replicas: 5, nodeCount: 5, errors: false}}
	newNodeGroup := func(t *testing.T, controller *machineController, testConfig *testConfig) (*nodegroup, error) {
		if testConfig.machineDeployment != nil {
			return newNodegroupFromMachineDeployment(controller, testConfig.machineDeployment)
		}
		return newNodegroupFromMachineSet(controller, testConfig.machineSet)
	}
	test := func(t *testing.T, tc testCase, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t)
		defer stop()
		ng, err := newNodeGroup(t, controller, testConfig)
		if tc.errors && err == nil {
			t.Fatal("expected an error")
		}
		if !tc.errors && ng == nil {
			t.Fatalf("test case logic error: %v", err)
		}
		if tc.errors {
			return
		}
		if ng == nil {
			t.Fatal("expected nodegroup to be non-nil")
		}
		var expectedName string
		switch v := (ng.scalableResource).(type) {
		case *machineSetScalableResource:
			expectedName = testConfig.spec.machineSetName
		case *machineDeploymentScalableResource:
			expectedName = testConfig.spec.machineDeploymentName
		default:
			t.Fatalf("unexpected type: %T", v)
		}
		expectedID := path.Join(testConfig.spec.namespace, expectedName)
		expectedDebug := fmt.Sprintf(debugFormat, expectedID, tc.minSize, tc.maxSize, tc.replicas)
		if ng.Name() != expectedName {
			t.Errorf("expected %q, got %q", expectedName, ng.Name())
		}
		if ng.Namespace() != testConfig.spec.namespace {
			t.Errorf("expected %q, got %q", testConfig.spec.namespace, ng.Namespace())
		}
		if ng.MinSize() != tc.minSize {
			t.Errorf("expected %v, got %v", tc.minSize, ng.MinSize())
		}
		if ng.MaxSize() != tc.maxSize {
			t.Errorf("expected %v, got %v", tc.maxSize, ng.MaxSize())
		}
		if ng.Id() != expectedID {
			t.Errorf("expected %q, got %q", expectedID, ng.Id())
		}
		if ng.Debug() != expectedDebug {
			t.Errorf("expected %q, got %q", expectedDebug, ng.Debug())
		}
		if _, err := ng.TemplateNodeInfo(); err != cloudprovider.ErrNotImplemented {
			t.Error("expected error")
		}
		if exists := ng.Exist(); !exists {
			t.Errorf("expected %t, got %t", true, exists)
		}
		if _, err := ng.Create(); err != cloudprovider.ErrAlreadyExist {
			t.Error("expected error")
		}
		if err := ng.Delete(); err != cloudprovider.ErrNotImplemented {
			t.Error("expected error")
		}
		if result := ng.Autoprovisioned(); result {
			t.Errorf("expected %t, got %t", false, result)
		}
	}
	t.Run("MachineSet", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				test(t, tc, createMachineSetTestConfig(testNamespace, tc.nodeCount, tc.replicas, tc.annotations))
			})
		}
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				test(t, tc, createMachineDeploymentTestConfig(testNamespace, tc.nodeCount, tc.replicas, tc.annotations))
			})
		}
	})
}
func TestNodeGroupIncreaseSizeErrors(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		description	string
		delta		int
		initial		int32
		errorMsg	string
	}
	testCases := []testCase{{description: "errors because delta is negative", delta: -1, initial: 3, errorMsg: "size increase must be positive"}, {description: "errors because initial+delta > maxSize", delta: 8, initial: 3, errorMsg: "size increase too large - desired:11 max:10"}}
	test := func(t *testing.T, tc *testCase, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l := len(nodegroups); l != 1 {
			t.Fatalf("expected 1 nodegroup, got %d", l)
		}
		ng := nodegroups[0]
		currReplicas, err := ng.TargetSize()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if currReplicas != int(tc.initial) {
			t.Errorf("expected %v, got %v", tc.initial, currReplicas)
		}
		errors := len(tc.errorMsg) > 0
		err = ng.IncreaseSize(tc.delta)
		if errors && err == nil {
			t.Fatal("expected an error")
		}
		if !errors && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(err.Error(), tc.errorMsg) {
			t.Errorf("expected error message to contain %q, got %q", tc.errorMsg, err.Error())
		}
		switch v := (ng.scalableResource).(type) {
		case *machineSetScalableResource:
			ms, err := ng.machineapiClient.MachineSets(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(ms.Spec.Replicas, 0); actual != tc.initial {
				t.Errorf("expected %v, got %v", tc.initial, actual)
			}
		case *machineDeploymentScalableResource:
			md, err := ng.machineapiClient.MachineDeployments(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(md.Spec.Replicas, 0); actual != tc.initial {
				t.Errorf("expected %v, got %v", tc.initial, actual)
			}
		default:
			t.Errorf("unexpected type: %T", v)
		}
	}
	t.Run("MachineSet", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
				test(t, &tc, createMachineSetTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
			})
		}
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
				test(t, &tc, createMachineDeploymentTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
			})
		}
	})
}
func TestNodeGroupIncreaseSize(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		description	string
		delta		int
		initial		int32
		expected	int32
	}
	test := func(t *testing.T, tc *testCase, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l := len(nodegroups); l != 1 {
			t.Fatalf("expected 1 nodegroup, got %d", l)
		}
		ng := nodegroups[0]
		currReplicas, err := ng.TargetSize()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if currReplicas != int(tc.initial) {
			t.Errorf("initially expected %v, got %v", tc.initial, currReplicas)
		}
		if err := ng.IncreaseSize(tc.delta); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		switch v := (ng.scalableResource).(type) {
		case *machineSetScalableResource:
			ms, err := ng.machineapiClient.MachineSets(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(ms.Spec.Replicas, 0); actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		case *machineDeploymentScalableResource:
			md, err := ng.machineapiClient.MachineDeployments(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(md.Spec.Replicas, 0); actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		default:
			t.Errorf("unexpected type: %T", v)
		}
	}
	annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
	t.Run("MachineSet", func(t *testing.T) {
		tc := testCase{description: "increase by 1", initial: 3, expected: 4, delta: 1}
		test(t, &tc, createMachineSetTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		tc := testCase{description: "increase by 1", initial: 3, expected: 4, delta: 1}
		test(t, &tc, createMachineDeploymentTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
	})
}
func TestNodeGroupDecreaseTargetSize(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		description	string
		delta		int
		initial		int32
		expected	int32
	}
	test := func(t *testing.T, tc *testCase, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l := len(nodegroups); l != 1 {
			t.Fatalf("expected 1 nodegroup, got %d", l)
		}
		ng := nodegroups[0]
		currReplicas, err := ng.TargetSize()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if currReplicas != int(tc.initial) {
			t.Errorf("initially expected %v, got %v", tc.initial, currReplicas)
		}
		if err := controller.nodeInformer.GetStore().Delete(testConfig.nodes[0]); err != nil {
			t.Fatalf("failed to add new node: %v", err)
		}
		if err := controller.machineInformer.Informer().GetStore().Add(testConfig.machines[0]); err != nil {
			t.Fatalf("failed to add new machine: %v", err)
		}
		if err := ng.DecreaseTargetSize(tc.delta); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		switch v := (ng.scalableResource).(type) {
		case *machineSetScalableResource:
			ms, err := ng.machineapiClient.MachineSets(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(ms.Spec.Replicas, 0); actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		case *machineDeploymentScalableResource:
			md, err := ng.machineapiClient.MachineDeployments(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(md.Spec.Replicas, 0); actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		default:
			t.Errorf("unexpected type: %T", v)
		}
	}
	annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
	t.Run("MachineSet", func(t *testing.T) {
		tc := testCase{description: "decrease by 1", initial: 3, expected: 2, delta: -1}
		test(t, &tc, createMachineSetTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		tc := testCase{description: "decrease by 1", initial: 3, expected: 2, delta: -1}
		test(t, &tc, createMachineDeploymentTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
	})
}
func TestNodeGroupDecreaseSizeErrors(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		description	string
		delta		int
		initial		int32
		errorMsg	string
	}
	testCases := []testCase{{description: "errors because delta is positive", delta: 1, initial: 3, errorMsg: "size decrease must be negative"}, {description: "errors because initial+delta < len(nodes)", delta: -1, initial: 3, errorMsg: "attempt to delete existing nodes targetSize:3 delta:-1 existingNodes: 3"}}
	test := func(t *testing.T, tc *testCase, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l := len(nodegroups); l != 1 {
			t.Fatalf("expected 1 nodegroup, got %d", l)
		}
		ng := nodegroups[0]
		currReplicas, err := ng.TargetSize()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if currReplicas != int(tc.initial) {
			t.Errorf("expected %v, got %v", tc.initial, currReplicas)
		}
		errors := len(tc.errorMsg) > 0
		err = ng.DecreaseTargetSize(tc.delta)
		if errors && err == nil {
			t.Fatal("expected an error")
		}
		if !errors && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(err.Error(), tc.errorMsg) {
			t.Errorf("expected error message to contain %q, got %q", tc.errorMsg, err.Error())
		}
		switch v := (ng.scalableResource).(type) {
		case *machineSetScalableResource:
			ms, err := ng.machineapiClient.MachineSets(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(ms.Spec.Replicas, 0); actual != tc.initial {
				t.Errorf("expected %v, got %v", tc.initial, actual)
			}
		case *machineDeploymentScalableResource:
			md, err := ng.machineapiClient.MachineDeployments(ng.Namespace()).Get(ng.Name(), v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(md.Spec.Replicas, 0); actual != tc.initial {
				t.Errorf("expected %v, got %v", tc.initial, actual)
			}
		default:
			t.Errorf("unexpected type: %T", v)
		}
	}
	t.Run("MachineSet", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
				test(t, &tc, createMachineSetTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
			})
		}
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
				test(t, &tc, createMachineDeploymentTestConfig(testNamespace, int(tc.initial), tc.initial, annotations))
			})
		}
	})
}
func TestNodeGroupDeleteNodes(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := func(t *testing.T, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l := len(nodegroups); l != 1 {
			t.Fatalf("expected 1 nodegroup, got %d", l)
		}
		ng := nodegroups[0]
		nodeNames, err := ng.Nodes()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(nodeNames) != len(testConfig.nodes) {
			t.Fatalf("expected len=%v, got len=%v", len(testConfig.nodes), len(nodeNames))
		}
		sort.SliceStable(nodeNames, func(i, j int) bool {
			return nodeNames[i].Id < nodeNames[j].Id
		})
		for i := 0; i < len(nodeNames); i++ {
			if nodeNames[i].Id != testConfig.nodes[i].Spec.ProviderID {
				t.Fatalf("expected %q, got %q", testConfig.nodes[i].Spec.ProviderID, nodeNames[i].Id)
			}
		}
		if err := ng.DeleteNodes(testConfig.nodes[5:]); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for i := 5; i < len(testConfig.machines); i++ {
			machine, err := controller.clusterClientset.MachineV1beta1().Machines(testConfig.machines[i].Namespace).Get(testConfig.machines[i].Name, v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, found := machine.Annotations[machineDeleteAnnotationKey]; !found {
				t.Errorf("expected annotation %q on machine %s", machineDeleteAnnotationKey, machine.Name)
			}
		}
		switch v := (ng.scalableResource).(type) {
		case *machineSetScalableResource:
			updatedMachineSet, err := controller.clusterClientset.MachineV1beta1().MachineSets(testConfig.machineSet.Namespace).Get(testConfig.machineSet.Name, v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(updatedMachineSet.Spec.Replicas, 0); actual != 5 {
				t.Fatalf("expected 5 nodes, got %v", actual)
			}
		case *machineDeploymentScalableResource:
			updatedMachineDeployment, err := controller.clusterClientset.MachineV1beta1().MachineDeployments(testConfig.machineDeployment.Namespace).Get(testConfig.machineDeployment.Name, v1.GetOptions{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual := pointer.Int32PtrDerefOr(updatedMachineDeployment.Spec.Replicas, 0); actual != 5 {
				t.Fatalf("expected 5 nodes, got %v", actual)
			}
		default:
			t.Errorf("unexpected type: %T", v)
		}
	}
	t.Run("MachineSet", func(t *testing.T) {
		test(t, createMachineSetTestConfig(testNamespace, 10, 10, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}))
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		test(t, createMachineDeploymentTestConfig(testNamespace, 10, 10, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}))
	})
}
func TestNodeGroupMachineSetDeleteNodesWithMismatchedNodes(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := func(t *testing.T, expected int, testConfigs []*testConfig) {
		t.Helper()
		testConfig0, testConfig1 := testConfigs[0], testConfigs[1]
		controller, stop := mustCreateTestController(t, testConfigs...)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if l := len(nodegroups); l != expected {
			t.Fatalf("expected %d, got %d", expected, l)
		}
		ng0, err := controller.nodeGroupForNode(testConfig0.nodes[0])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		ng1, err := controller.nodeGroupForNode(testConfig1.nodes[0])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err0 := ng0.DeleteNodes(testConfig1.nodes)
		if err0 == nil {
			t.Error("expected an error")
		}
		expectedErr0 := `node "test-namespace1-machineset-0-nodeid-0" doesn't belong to node group "test-namespace0/machineset-0"`
		if testConfig0.machineDeployment != nil {
			expectedErr0 = `node "test-namespace1-machineset-0-nodeid-0" doesn't belong to node group "test-namespace0/machinedeployment-0"`
		}
		if !strings.Contains(err0.Error(), expectedErr0) {
			t.Errorf("expected: %q, got: %q", expectedErr0, err0.Error())
		}
		err1 := ng1.DeleteNodes(testConfig0.nodes)
		if err1 == nil {
			t.Error("expected an error")
		}
		expectedErr1 := `node "test-namespace0-machineset-0-nodeid-0" doesn't belong to node group "test-namespace1/machineset-0"`
		if testConfig1.machineDeployment != nil {
			expectedErr1 = `node "test-namespace0-machineset-0-nodeid-0" doesn't belong to node group "test-namespace1/machinedeployment-0"`
		}
		if !strings.Contains(err1.Error(), expectedErr1) {
			t.Errorf("expected: %q, got: %q", expectedErr1, err1.Error())
		}
		if err := ng0.DeleteNodes(testConfig0.nodes); err == nil {
			t.Error("expected error")
		}
		if err := ng1.DeleteNodes(testConfig1.nodes); err == nil {
			t.Error("expected error")
		}
	}
	annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "3"}
	t.Run("MachineSet", func(t *testing.T) {
		testConfig0 := createMachineSetTestConfigs(testNamespace+"0", 1, 2, 2, annotations)
		testConfig1 := createMachineSetTestConfigs(testNamespace+"1", 1, 2, 2, annotations)
		test(t, 2, append(testConfig0, testConfig1...))
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		testConfig0 := createMachineDeploymentTestConfigs(testNamespace+"0", 1, 2, 2, annotations)
		testConfig1 := createMachineDeploymentTestConfigs(testNamespace+"1", 1, 2, 2, annotations)
		test(t, 2, append(testConfig0, testConfig1...))
	})
}
