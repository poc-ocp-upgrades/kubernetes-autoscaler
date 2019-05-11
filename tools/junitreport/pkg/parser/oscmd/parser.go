package oscmd

import (
	"k8s.io/autoscaler/tools/junitreport/pkg/builder"
	"k8s.io/autoscaler/tools/junitreport/pkg/parser"
	"k8s.io/autoscaler/tools/junitreport/pkg/parser/stack"
)

func NewParser(builder builder.TestSuitesBuilder, stream bool) parser.TestOutputParser {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return stack.NewParser(builder, newTestDataParser(), newTestSuiteDataParser(), stream)
}
