package stack

import (
 "bufio"
 "fmt"
 "os"
 "strings"
 "k8s.io/autoscaler/tools/junitreport/pkg/api"
 "k8s.io/autoscaler/tools/junitreport/pkg/builder"
 "k8s.io/autoscaler/tools/junitreport/pkg/parser"
)

func NewParser(builder builder.TestSuitesBuilder, testParser TestDataParser, suiteParser TestSuiteDataParser, stream bool) parser.TestOutputParser {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &testOutputParser{builder: builder, testParser: testParser, suiteParser: suiteParser, stream: stream}
}

type testOutputParser struct {
 builder     builder.TestSuitesBuilder
 testParser  TestDataParser
 suiteParser TestSuiteDataParser
 stream      bool
}

func (p *testOutputParser) Parse(input *bufio.Scanner) (*api.TestSuites, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 inProgress := NewTestSuiteStack()
 var currentTest *api.TestCase
 var currentResult api.TestResult
 var currentOutput []string
 var currentMessage string
 for input.Scan() {
  line := input.Text()
  isTestOutput := true
  if p.testParser.MarksBeginning(line) {
   currentTest = &api.TestCase{}
   currentResult = api.TestResultFail
   currentOutput = []string{}
   currentMessage = ""
  }
  if name, contained := p.testParser.ExtractName(line); contained {
   currentTest.Name = name
  }
  if result, contained := p.testParser.ExtractResult(line); contained {
   currentResult = result
  }
  if duration, contained := p.testParser.ExtractDuration(line); contained {
   if err := currentTest.SetDuration(duration); err != nil {
    return nil, err
   }
  }
  if message, contained := p.testParser.ExtractMessage(line); contained {
   currentMessage = message
  }
  if p.testParser.MarksCompletion(line) {
   currentOutput = append(currentOutput, line)
   switch currentResult {
   case api.TestResultSkip:
    currentTest.MarkSkipped(currentMessage)
   case api.TestResultFail:
    output := strings.Join(currentOutput, "\n")
    currentTest.MarkFailed(currentMessage, output)
   }
   if inProgress.Peek() == nil {
    return nil, fmt.Errorf("found test case %q outside of a test suite", currentTest.Name)
   }
   inProgress.Peek().AddTestCase(currentTest)
   currentTest = &api.TestCase{}
  }
  if p.suiteParser.MarksBeginning(line) {
   inProgress.Push(&api.TestSuite{})
   isTestOutput = false
  }
  if name, contained := p.suiteParser.ExtractName(line); contained {
   inProgress.Peek().Name = name
   isTestOutput = false
  }
  if properties, contained := p.suiteParser.ExtractProperties(line); contained {
   for propertyName := range properties {
    inProgress.Peek().AddProperty(propertyName, properties[propertyName])
   }
   isTestOutput = false
  }
  if p.suiteParser.MarksCompletion(line) {
   if p.stream {
    fmt.Fprintln(os.Stdout, line)
   }
   p.builder.AddSuite(inProgress.Pop())
   isTestOutput = false
  }
  if isTestOutput {
   currentOutput = append(currentOutput, line)
  }
 }
 return p.builder.Build(), nil
}
