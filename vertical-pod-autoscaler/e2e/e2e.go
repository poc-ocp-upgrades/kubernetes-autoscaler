package autoscaling

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
	"github.com/golang/glog"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeutils "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/util/logs"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/kubernetes/pkg/version"
	commontest "k8s.io/kubernetes/test/e2e/common"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/ginkgowrapper"
	"k8s.io/kubernetes/test/e2e/framework/metrics"
)

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	c, err := framework.LoadClientset()
	if err != nil {
		glog.Fatal("Error loading client: ", err)
	}
	if framework.TestContext.CleanStart {
		deleted, err := framework.DeleteNamespaces(c, nil, []string{metav1.NamespaceSystem, metav1.NamespaceDefault, metav1.NamespacePublic})
		if err != nil {
			framework.Failf("Error deleting orphaned namespaces: %v", err)
		}
		glog.Infof("Waiting for deletion of the following namespaces: %v", deleted)
		if err := framework.WaitForNamespacesDeleted(c, deleted, framework.NamespaceCleanupTimeout); err != nil {
			framework.Failf("Failed to delete orphaned namespaces %v: %v", deleted, err)
		}
	}
	framework.ExpectNoError(framework.WaitForAllNodesSchedulable(c, framework.TestContext.NodeSchedulableTimeout))
	podStartupTimeout := framework.TestContext.SystemPodsStartupTimeout
	if err := framework.WaitForPodsRunningReady(c, metav1.NamespaceSystem, int32(framework.TestContext.MinStartupPods), int32(framework.TestContext.AllowedNotReadyNodes), podStartupTimeout, framework.ImagePullerLabels); err != nil {
		framework.DumpAllNamespaceInfo(c, metav1.NamespaceSystem)
		framework.LogFailedContainers(c, metav1.NamespaceSystem, framework.Logf)
		framework.Failf("Error waiting for all pods to be running and ready: %v", err)
	}
	if err := framework.WaitForPodsSuccess(c, metav1.NamespaceSystem, framework.ImagePullerLabels, framework.ImagePrePullingTimeout); err != nil {
		framework.Logf("WARNING: Image pulling pods failed to enter success in %v: %v", framework.ImagePrePullingTimeout, err)
	}
	if framework.TestContext.DumpLogsOnFailure {
		logFunc := framework.Logf
		if framework.TestContext.ReportDir != "" {
			filePath := path.Join(framework.TestContext.ReportDir, "nethealth.txt")
			file, err := os.Create(filePath)
			if err != nil {
				framework.Logf("Failed to create a file with network health data %v: %v\nPrinting to stdout", filePath, err)
			} else {
				defer file.Close()
				if err = file.Chmod(0644); err != nil {
					framework.Logf("Failed to chmod to 644 of %v: %v", filePath, err)
				}
				logFunc = framework.GetLogToFileFunc(file)
				framework.Logf("Dumping network health container logs from all nodes to file %v", filePath)
			}
		} else {
			framework.Logf("Dumping network health container logs from all nodes...")
		}
		framework.LogContainersInPodsWithLabels(c, metav1.NamespaceSystem, framework.ImagePullerLabels, "nethealth", logFunc)
	}
	framework.Logf("e2e test version: %s", version.Get().GitVersion)
	dc := c.DiscoveryClient
	serverVersion, serverErr := dc.ServerVersion()
	if serverErr != nil {
		framework.Logf("Unexpected server error retrieving version: %v", serverErr)
	}
	if serverVersion != nil {
		framework.Logf("kube-apiserver version: %s", serverVersion.GitVersion)
	}
	commontest.CurrentSuite = commontest.E2E
	return nil
}, func(data []byte) {
	framework.Logf("No cloud config support.")
})
var _ = ginkgo.SynchronizedAfterSuite(func() {
	framework.Logf("Running AfterSuite actions on all node")
	framework.RunCleanupActions()
}, func() {
	framework.Logf("Running AfterSuite actions on node 1")
	if framework.TestContext.ReportDir != "" {
		framework.CoreDump(framework.TestContext.ReportDir)
	}
	if framework.TestContext.GatherSuiteMetricsAfterTest {
		if err := gatherTestSuiteMetrics(); err != nil {
			framework.Logf("Error gathering metrics: %v", err)
		}
	}
})

func gatherTestSuiteMetrics() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	framework.Logf("Gathering metrics")
	c, err := framework.LoadClientset()
	if err != nil {
		return fmt.Errorf("error loading client: %v", err)
	}
	grabber, err := metrics.NewMetricsGrabber(c, nil, !framework.ProviderIs("kubemark"), true, true, true, framework.TestContext.IncludeClusterAutoscalerMetrics)
	if err != nil {
		return fmt.Errorf("failed to create MetricsGrabber: %v", err)
	}
	received, err := grabber.Grab()
	if err != nil {
		return fmt.Errorf("failed to grab metrics: %v", err)
	}
	metricsForE2E := (*framework.MetricsForE2E)(&received)
	metricsJSON := metricsForE2E.PrintJSON()
	if framework.TestContext.ReportDir != "" {
		filePath := path.Join(framework.TestContext.ReportDir, "MetricsForE2ESuite_"+time.Now().Format(time.RFC3339)+".json")
		if err := ioutil.WriteFile(filePath, []byte(metricsJSON), 0644); err != nil {
			return fmt.Errorf("error writing to %q: %v", filePath, err)
		}
	} else {
		framework.Logf("\n\nTest Suite Metrics:\n%s\n\n", metricsJSON)
	}
	return nil
}
func RunE2ETests(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	runtimeutils.ReallyCrash = true
	logs.InitLogs()
	defer logs.FlushLogs()
	gomega.RegisterFailHandler(ginkgowrapper.Fail)
	if config.GinkgoConfig.FocusString == "" && config.GinkgoConfig.SkipString == "" {
		config.GinkgoConfig.SkipString = `\[Flaky\]|\[Feature:.+\]`
	}
	var r []ginkgo.Reporter
	if framework.TestContext.ReportDir != "" {
		if err := os.MkdirAll(framework.TestContext.ReportDir, 0755); err != nil {
			glog.Errorf("Failed creating report directory: %v", err)
		} else {
			r = append(r, reporters.NewJUnitReporter(path.Join(framework.TestContext.ReportDir, fmt.Sprintf("junit_%v%02d.xml", framework.TestContext.ReportPrefix, config.GinkgoConfig.ParallelNode))))
		}
	}
	glog.Infof("Starting e2e run %q on Ginkgo node %d", framework.RunId, config.GinkgoConfig.ParallelNode)
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Kubernetes e2e suite", r)
}
