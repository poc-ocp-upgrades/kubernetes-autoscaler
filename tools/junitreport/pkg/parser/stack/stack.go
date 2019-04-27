package stack

import (
	"fmt"
	"k8s.io/autoscaler/tools/junitreport/pkg/api"
)

type TestSuiteStack interface {
	Push(pkg *api.TestSuite)
	Pop() *api.TestSuite
	Peek() *api.TestSuite
	IsEmpty() bool
}

func NewTestSuiteStack() TestSuiteStack {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &testSuiteStack{head: nil}
}

type testSuiteStack struct{ head *testSuiteNode }

func (s *testSuiteStack) Push(data *api.TestSuite) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	newNode := &testSuiteNode{Member: data, Next: s.head}
	s.head = newNode
}
func (s *testSuiteStack) Pop() *api.TestSuite {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.IsEmpty() {
		return nil
	}
	oldNode := s.head
	s.head = s.head.Next
	return oldNode.Member
}
func (s *testSuiteStack) Peek() *api.TestSuite {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.IsEmpty() {
		return nil
	}
	return s.head.Member
}
func (s *testSuiteStack) IsEmpty() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.head == nil
}

type testSuiteNode struct {
	Member	*api.TestSuite
	Next	*testSuiteNode
}

func (n *testSuiteNode) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("{Member: %s, Next: %s}", n.Member, n.Next.String())
}
