package daemonset

import (
	"strings"
	"testing"
	"time"
	"k8s.io/autoscaler/cluster-autoscaler/simulator"
	. "k8s.io/autoscaler/cluster-autoscaler/utils/test"
	apiv1 "k8s.io/api/core/v1"
	extensionsv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schedulercache "k8s.io/kubernetes/pkg/scheduler/cache"
	"github.com/stretchr/testify/assert"
)

func TestGetDaemonSetPodsForNode(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	node := BuildTestNode("node", 1000, 1000)
	SetNodeReadyState(node, true, time.Now())
	nodeInfo := schedulercache.NewNodeInfo()
	nodeInfo.SetNode(node)
	predicateChecker := simulator.NewTestPredicateChecker()
	ds1 := newDaemonSet("ds1")
	ds2 := newDaemonSet("ds2")
	ds2.Spec.Template.Spec.NodeSelector = map[string]string{"foo": "bar"}
	pods := GetDaemonSetPodsForNode(nodeInfo, []*extensionsv1.DaemonSet{ds1, ds2}, predicateChecker)
	assert.Equal(t, 1, len(pods))
	assert.True(t, strings.HasPrefix(pods[0].Name, "ds1"))
	assert.Equal(t, 1, len(GetDaemonSetPodsForNode(nodeInfo, []*extensionsv1.DaemonSet{ds1}, predicateChecker)))
	assert.Equal(t, 0, len(GetDaemonSetPodsForNode(nodeInfo, []*extensionsv1.DaemonSet{ds2}, predicateChecker)))
	assert.Equal(t, 0, len(GetDaemonSetPodsForNode(nodeInfo, []*extensionsv1.DaemonSet{}, predicateChecker)))
}
func newDaemonSet(name string) *extensionsv1.DaemonSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &extensionsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: metav1.NamespaceDefault}, Spec: extensionsv1.DaemonSetSpec{Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"name": "simple-daemon", "type": "production"}}, Template: apiv1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"name": "simple-daemon", "type": "production"}}, Spec: apiv1.PodSpec{Containers: []apiv1.Container{{Image: "foo/bar"}}}}}}
}
