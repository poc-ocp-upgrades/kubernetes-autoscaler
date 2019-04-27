package openshiftmachineapi

import (
	"fmt"
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"
	"github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	fakeclusterapi "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/fake"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakekube "k8s.io/client-go/kubernetes/fake"
)

type testControllerShutdownFunc func()
type testConfig struct {
	spec			*testSpec
	machineDeployment	*v1beta1.MachineDeployment
	machineSet		*v1beta1.MachineSet
	machines		[]*v1beta1.Machine
	nodes			[]*apiv1.Node
}
type testSpec struct {
	annotations		map[string]string
	machineDeploymentName	string
	machineSetName		string
	namespace		string
	nodeCount		int
	replicaCount		int32
	rootIsMachineDeployment	bool
}

func mustCreateTestController(t *testing.T, testConfigs ...*testConfig) (*machineController, testControllerShutdownFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Helper()
	nodeObjects := make([]runtime.Object, 0)
	machineObjects := make([]runtime.Object, 0)
	for _, config := range testConfigs {
		for i := range config.nodes {
			nodeObjects = append(nodeObjects, config.nodes[i])
		}
		for i := range config.machines {
			machineObjects = append(machineObjects, config.machines[i])
		}
		machineObjects = append(machineObjects, config.machineSet)
		if config.machineDeployment != nil {
			machineObjects = append(machineObjects, config.machineDeployment)
		}
	}
	kubeclientSet := fakekube.NewSimpleClientset(nodeObjects...)
	clusterclientSet := fakeclusterapi.NewSimpleClientset(machineObjects...)
	controller, err := newMachineController(kubeclientSet, clusterclientSet, true)
	if err != nil {
		t.Fatal("failed to create test controller")
	}
	stopCh := make(chan struct{})
	if err := controller.run(stopCh); err != nil {
		t.Fatalf("failed to run controller: %v", err)
	}
	return controller, func() {
		close(stopCh)
	}
}
func createMachineSetTestConfig(namespace string, nodeCount int, replicaCount int32, annotations map[string]string) *testConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return createTestConfigs(createTestSpecs(namespace, 1, nodeCount, replicaCount, false, annotations)...)[0]
}
func createMachineSetTestConfigs(namespace string, configCount, nodeCount int, replicaCount int32, annotations map[string]string) []*testConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return createTestConfigs(createTestSpecs(namespace, configCount, nodeCount, replicaCount, false, annotations)...)
}
func createMachineDeploymentTestConfig(namespace string, nodeCount int, replicaCount int32, annotations map[string]string) *testConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return createTestConfigs(createTestSpecs(namespace, 1, nodeCount, replicaCount, true, annotations)...)[0]
}
func createMachineDeploymentTestConfigs(namespace string, configCount, nodeCount int, replicaCount int32, annotations map[string]string) []*testConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return createTestConfigs(createTestSpecs(namespace, configCount, nodeCount, replicaCount, true, annotations)...)
}
func createTestSpecs(namespace string, scalableResourceCount, nodeCount int, replicaCount int32, isMachineDeployment bool, annotations map[string]string) []testSpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var specs []testSpec
	for i := 0; i < scalableResourceCount; i++ {
		specs = append(specs, testSpec{annotations: annotations, machineDeploymentName: fmt.Sprintf("machinedeployment-%d", i), machineSetName: fmt.Sprintf("machineset-%d", i), namespace: strings.ToLower(namespace), nodeCount: nodeCount, replicaCount: replicaCount, rootIsMachineDeployment: isMachineDeployment})
	}
	return specs
}
func createTestConfigs(specs ...testSpec) []*testConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var result []*testConfig
	for i, spec := range specs {
		config := &testConfig{spec: &specs[i], nodes: make([]*apiv1.Node, spec.nodeCount), machines: make([]*v1beta1.Machine, spec.nodeCount)}
		config.machineSet = &v1beta1.MachineSet{TypeMeta: v1.TypeMeta{Kind: "MachineSet"}, ObjectMeta: v1.ObjectMeta{Name: spec.machineSetName, Namespace: spec.namespace, UID: types.UID(spec.machineSetName)}}
		if !spec.rootIsMachineDeployment {
			config.machineSet.ObjectMeta.Annotations = spec.annotations
			config.machineSet.Spec.Replicas = int32ptr(spec.replicaCount)
		} else {
			config.machineDeployment = &v1beta1.MachineDeployment{TypeMeta: v1.TypeMeta{Kind: "MachineDeployment"}, ObjectMeta: v1.ObjectMeta{Name: spec.machineDeploymentName, Namespace: spec.namespace, UID: types.UID(spec.machineDeploymentName), Annotations: spec.annotations}, Spec: v1beta1.MachineDeploymentSpec{Replicas: int32ptr(spec.replicaCount)}}
			config.machineSet.OwnerReferences = make([]v1.OwnerReference, 1)
			config.machineSet.OwnerReferences[0] = v1.OwnerReference{Name: config.machineDeployment.Name, Kind: config.machineDeployment.Kind, UID: config.machineDeployment.UID}
		}
		machineOwner := v1.OwnerReference{Name: config.machineSet.Name, Kind: config.machineSet.Kind, UID: config.machineSet.UID}
		for j := 0; j < spec.nodeCount; j++ {
			config.nodes[j], config.machines[j] = makeLinkedNodeAndMachine(j, spec.namespace, machineOwner)
		}
		result = append(result, config)
	}
	return result
}
func makeLinkedNodeAndMachine(i int, namespace string, owner v1.OwnerReference) (*apiv1.Node, *v1beta1.Machine) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := &apiv1.Node{TypeMeta: v1.TypeMeta{Kind: "Node"}, ObjectMeta: v1.ObjectMeta{Name: fmt.Sprintf("%s-%s-node-%d", namespace, owner.Name, i), Annotations: map[string]string{machineAnnotationKey: fmt.Sprintf("%s/%s-%s-machine-%d", namespace, namespace, owner.Name, i)}}, Spec: apiv1.NodeSpec{ProviderID: fmt.Sprintf("%s-%s-nodeid-%d", namespace, owner.Name, i)}}
	machine := &v1beta1.Machine{TypeMeta: v1.TypeMeta{Kind: "Machine"}, ObjectMeta: v1.ObjectMeta{Name: fmt.Sprintf("%s-%s-machine-%d", namespace, owner.Name, i), Namespace: namespace, OwnerReferences: []v1.OwnerReference{{Name: owner.Name, Kind: owner.Kind, UID: owner.UID}}}, Status: v1beta1.MachineStatus{NodeRef: &apiv1.ObjectReference{Kind: node.Kind, Name: node.Name}}}
	return node, machine
}
func int32ptr(v int32) *int32 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v
}
func addTestConfigs(t *testing.T, controller *machineController, testConfigs ...*testConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Helper()
	for _, config := range testConfigs {
		if config.machineDeployment != nil {
			if err := controller.machineDeploymentInformer.Informer().GetStore().Add(config.machineDeployment); err != nil {
				return err
			}
		}
		if err := controller.machineSetInformer.Informer().GetStore().Add(config.machineSet); err != nil {
			return err
		}
		for i := range config.machines {
			if err := controller.machineInformer.Informer().GetStore().Add(config.machines[i]); err != nil {
				return err
			}
		}
		for i := range config.nodes {
			if err := controller.nodeInformer.GetStore().Add(config.nodes[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
func deleteTestConfigs(t *testing.T, controller *machineController, testConfigs ...*testConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Helper()
	for _, config := range testConfigs {
		for i := range config.nodes {
			if err := controller.nodeInformer.GetStore().Delete(config.nodes[i]); err != nil {
				return err
			}
		}
		for i := range config.machines {
			if err := controller.machineInformer.Informer().GetStore().Delete(config.machines[i]); err != nil {
				return err
			}
		}
		if err := controller.machineSetInformer.Informer().GetStore().Delete(config.machineSet); err != nil {
			return err
		}
		if config.machineDeployment != nil {
			if err := controller.machineDeploymentInformer.Informer().GetStore().Delete(config.machineDeployment); err != nil {
				return err
			}
		}
	}
	return nil
}
func TestControllerFindMachineByID(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		description	string
		name		string
		namespace	string
		lookupSucceeds	bool
	}
	var testCases = []testCase{{description: "lookup fails", lookupSucceeds: false, name: "machine-does-not-exist", namespace: "namespace-does-not-exist"}, {description: "lookup fails in valid namespace", lookupSucceeds: false, name: "machine-does-not-exist-in-existing-namespace"}, {description: "lookup succeeds", lookupSucceeds: true}}
	test := func(t *testing.T, tc testCase, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		machine, err := controller.findMachine(path.Join(tc.namespace, tc.name))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tc.lookupSucceeds && machine == nil {
			t.Error("expected success, findMachine failed")
		}
		if tc.lookupSucceeds && machine != nil {
			if machine.Name != tc.name {
				t.Errorf("expected %q, got %q", tc.name, machine.Name)
			}
			if machine.Namespace != tc.namespace {
				t.Errorf("expected %q, got %q", tc.namespace, machine.Namespace)
			}
		}
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
			if tc.name == "" {
				tc.name = testConfig.machines[0].Name
			}
			if tc.namespace == "" {
				tc.namespace = testConfig.machines[0].Namespace
			}
			test(t, tc, testConfig)
		})
	}
}
func TestControllerFindMachineOwner(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
	controller, stop := mustCreateTestController(t, testConfig)
	defer stop()
	testResult1, err := controller.findMachineOwner(testConfig.machines[0].DeepCopy())
	if err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}
	if testResult1 == nil {
		t.Fatal("expected non-nil result")
	}
	if testConfig.spec.machineSetName != testResult1.Name {
		t.Errorf("expected %q, got %q", testConfig.spec.machineSetName, testResult1.Name)
	}
	testMachine2 := testConfig.machines[0].DeepCopy()
	testMachine2.OwnerReferences[0].UID = "does-not-match-machineset"
	testResult2, err := controller.findMachineOwner(testMachine2)
	if err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}
	if testResult2 != nil {
		t.Fatal("expected nil result")
	}
	if err := controller.machineSetInformer.Informer().GetStore().Delete(testResult1); err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}
	testResult3, err := controller.findMachineOwner(testConfig.machines[0].DeepCopy())
	if err != nil {
		t.Fatalf("unexpected error, got %v", err)
	}
	if testResult3 != nil {
		t.Fatal("expected lookup to fail")
	}
}
func TestControllerFindMachineByNodeProviderID(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
	controller, stop := mustCreateTestController(t, testConfig)
	defer stop()
	machine, err := controller.findMachineByNodeProviderID(testConfig.nodes[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if machine == nil {
		t.Fatal("expected to find machine")
	}
	if !reflect.DeepEqual(machine, testConfig.machines[0]) {
		t.Fatalf("expected machines to be equal - expected %+v, got %+v", testConfig.machines[0], machine)
	}
	node := testConfig.nodes[0].DeepCopy()
	node.Spec.ProviderID = ""
	nonExistentMachine, err := controller.findMachineByNodeProviderID(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nonExistentMachine != nil {
		t.Fatal("expected find to fail")
	}
	node = testConfig.nodes[0].DeepCopy()
	delete(node.Annotations, machineAnnotationKey)
	if err := controller.nodeInformer.GetStore().Update(node); err != nil {
		t.Fatalf("unexpected error updating node, got %v", err)
	}
	nonExistentMachine, err = controller.findMachineByNodeProviderID(testConfig.nodes[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nonExistentMachine != nil {
		t.Fatal("expected find to fail")
	}
}
func TestControllerFindNodeByNodeName(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
	controller, stop := mustCreateTestController(t, testConfig)
	defer stop()
	node, err := controller.findNodeByNodeName(testConfig.nodes[0].Name)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node == nil {
		t.Fatal("expected lookup to be successful")
	}
	node, err = controller.findNodeByNodeName(testConfig.nodes[0].Name + "non-existent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node != nil {
		t.Fatal("expected lookup to fail")
	}
}
func TestControllerMachinesInMachineSet(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testConfig1 := createMachineSetTestConfig("testConfig1", 5, 5, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
	controller, stop := mustCreateTestController(t, testConfig1)
	defer stop()
	testConfig2 := createMachineSetTestConfig("testConfig2", 5, 5, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
	if err := addTestConfigs(t, controller, testConfig2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	machinesInTestObjs1, err := controller.machineInformer.Lister().Machines(testConfig1.spec.namespace).List(labels.Everything())
	if err != nil {
		t.Fatalf("error listing machines: %v", err)
	}
	machinesInTestObjs2, err := controller.machineInformer.Lister().Machines(testConfig2.spec.namespace).List(labels.Everything())
	if err != nil {
		t.Fatalf("error listing machines: %v", err)
	}
	actual := len(machinesInTestObjs1) + len(machinesInTestObjs2)
	expected := len(testConfig1.machines) + len(testConfig2.machines)
	if actual != expected {
		t.Fatalf("expected %d machines, got %d", expected, actual)
	}
	sort.Slice(machinesInTestObjs1, func(i, j int) bool {
		return machinesInTestObjs1[i].Name < machinesInTestObjs1[j].Name
	})
	sort.Slice(machinesInTestObjs2, func(i, j int) bool {
		return machinesInTestObjs2[i].Name < machinesInTestObjs2[j].Name
	})
	for i, m := range machinesInTestObjs1 {
		if m.Name != testConfig1.machines[i].Name {
			t.Errorf("expected %q, got %q", testConfig1.machines[i].Name, m.Name)
		}
		if m.Namespace != testConfig1.machines[i].Namespace {
			t.Errorf("expected %q, got %q", testConfig1.machines[i].Namespace, m.Namespace)
		}
	}
	for i, m := range machinesInTestObjs2 {
		if m.Name != testConfig2.machines[i].Name {
			t.Errorf("expected %q, got %q", testConfig2.machines[i].Name, m.Name)
		}
		if m.Namespace != testConfig2.machines[i].Namespace {
			t.Errorf("expected %q, got %q", testConfig2.machines[i].Namespace, m.Namespace)
		}
	}
	if !reflect.DeepEqual(testConfig1.machines, machinesInTestObjs1) {
		t.Fatalf("expected %+v, got %+v", testConfig1.machines, machinesInTestObjs1)
	}
	if !reflect.DeepEqual(testConfig2.machines, machinesInTestObjs2) {
		t.Fatalf("expected %+v, got %+v", testConfig2.machines, machinesInTestObjs2)
	}
}
func TestControllerLookupNodeGroupForNonExistentNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
	controller, stop := mustCreateTestController(t, testConfig)
	defer stop()
	node := testConfig.nodes[0].DeepCopy()
	node.Spec.ProviderID = "does-not-exist"
	ng, err := controller.nodeGroupForNode(node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ng != nil {
		t.Fatalf("unexpected nodegroup: %v", ng)
	}
}
func TestControllerNodeGroupForNodeWithMissingMachineOwner(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := func(t *testing.T, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		machine := testConfig.machines[0].DeepCopy()
		machine.OwnerReferences = []v1.OwnerReference{}
		if err := controller.machineInformer.Informer().GetStore().Update(machine); err != nil {
			t.Fatalf("unexpected error updating machine, got %v", err)
		}
		ng, err := controller.nodeGroupForNode(testConfig.nodes[0])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ng != nil {
			t.Fatalf("unexpected nodegroup: %v", ng)
		}
	}
	t.Run("MachineSet", func(t *testing.T) {
		testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
		test(t, testConfig)
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		testConfig := createMachineDeploymentTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"})
		test(t, testConfig)
	})
}
func TestControllerNodeGroupForNodeWithPositiveScalingBounds(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	test := func(t *testing.T, testConfig *testConfig) {
		controller, stop := mustCreateTestController(t, testConfig)
		defer stop()
		ng, err := controller.nodeGroupForNode(testConfig.nodes[0])
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ng != nil {
			t.Fatalf("unexpected nodegroup: %v", ng)
		}
	}
	t.Run("MachineSet", func(t *testing.T) {
		testConfig := createMachineSetTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "1"})
		test(t, testConfig)
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		testConfig := createMachineDeploymentTestConfig(testNamespace, 1, 1, map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "1"})
		test(t, testConfig)
	})
}
func TestControllerNodeGroups(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	assertNodegroupLen := func(t *testing.T, controller *machineController, expected int) {
		t.Helper()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := len(nodegroups); got != expected {
			t.Fatalf("expected %d, got %d", expected, got)
		}
	}
	annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "2"}
	controller, stop := mustCreateTestController(t)
	defer stop()
	assertNodegroupLen(t, controller, 0)
	machineSetConfigs := createMachineSetTestConfigs("MachineSet", 5, 1, 1, annotations)
	if err := addTestConfigs(t, controller, machineSetConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNodegroupLen(t, controller, 5)
	machineDeploymentConfigs := createMachineDeploymentTestConfigs("MachineDeployment", 2, 1, 1, annotations)
	if err := addTestConfigs(t, controller, machineDeploymentConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNodegroupLen(t, controller, 7)
	if err := deleteTestConfigs(t, controller, machineSetConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNodegroupLen(t, controller, 2)
	if err := deleteTestConfigs(t, controller, machineDeploymentConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNodegroupLen(t, controller, 0)
	annotations = map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "1"}
	machineSetConfigs = createMachineSetTestConfigs("MachineSet", 5, 1, 1, annotations)
	if err := addTestConfigs(t, controller, machineSetConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNodegroupLen(t, controller, 0)
	machineDeploymentConfigs = createMachineDeploymentTestConfigs("MachineDeployment", 2, 1, 1, annotations)
	if err := addTestConfigs(t, controller, machineDeploymentConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertNodegroupLen(t, controller, 0)
	annotations = map[string]string{nodeGroupMinSizeAnnotationKey: "-1", nodeGroupMaxSizeAnnotationKey: "1"}
	machineSetConfigs = createMachineSetTestConfigs("MachineSet", 5, 1, 1, annotations)
	if err := addTestConfigs(t, controller, machineSetConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := controller.nodeGroups(); err == nil {
		t.Fatalf("expected an error")
	}
	machineDeploymentConfigs = createMachineDeploymentTestConfigs("MachineDeployment", 2, 1, 1, annotations)
	if err := addTestConfigs(t, controller, machineDeploymentConfigs...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := controller.nodeGroups(); err == nil {
		t.Fatalf("expected an error")
	}
}
func TestControllerNodeGroupsNodeCount(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	type testCase struct {
		nodeGroups	int
		nodesPerGroup	int
	}
	var testCases = []testCase{{nodeGroups: 0, nodesPerGroup: 0}, {nodeGroups: 1, nodesPerGroup: 0}, {nodeGroups: 2, nodesPerGroup: 10}}
	test := func(t *testing.T, tc testCase, testConfigs []*testConfig) {
		controller, stop := mustCreateTestController(t, testConfigs...)
		defer stop()
		nodegroups, err := controller.nodeGroups()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := len(nodegroups); got != tc.nodeGroups {
			t.Fatalf("expected %d, got %d", tc.nodeGroups, got)
		}
		for i := range nodegroups {
			nodes, err := nodegroups[i].Nodes()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := len(nodes); got != tc.nodesPerGroup {
				t.Fatalf("expected %d, got %d", tc.nodesPerGroup, got)
			}
		}
	}
	annotations := map[string]string{nodeGroupMinSizeAnnotationKey: "1", nodeGroupMaxSizeAnnotationKey: "10"}
	t.Run("MachineSet", func(t *testing.T) {
		for _, tc := range testCases {
			test(t, tc, createMachineSetTestConfigs(testNamespace, tc.nodeGroups, tc.nodesPerGroup, int32(tc.nodesPerGroup), annotations))
		}
	})
	t.Run("MachineDeployment", func(t *testing.T) {
		for _, tc := range testCases {
			test(t, tc, createMachineDeploymentTestConfigs(testNamespace, tc.nodeGroups, tc.nodesPerGroup, int32(tc.nodesPerGroup), annotations))
		}
	})
}
