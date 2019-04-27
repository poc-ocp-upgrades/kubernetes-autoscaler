package azure

import (
	"fmt"
	"net/http"
	"testing"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/stretchr/testify/assert"
)

func TestSplitBlobURI(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expectedAccountName := "vhdstorage8h8pjybi9hbsl6"
	expectedContainerName := "vhds"
	expectedBlobPath := "osdisks/disk1234.vhd"
	accountName, containerName, blobPath, err := splitBlobURI("https://vhdstorage8h8pjybi9hbsl6.blob.core.windows.net/vhds/osdisks/disk1234.vhd")
	if accountName != expectedAccountName {
		t.Fatalf("incorrect account name. expected=%s actual=%s", expectedAccountName, accountName)
	}
	if containerName != expectedContainerName {
		t.Fatalf("incorrect account name. expected=%s actual=%s", expectedContainerName, containerName)
	}
	if blobPath != expectedBlobPath {
		t.Fatalf("incorrect account name. expected=%s actual=%s", expectedBlobPath, blobPath)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
func TestK8sLinuxVMNameParts(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	data := []struct {
		poolIdentifier, nameSuffix	string
		agentIndex			int
	}{{"agentpool1", "38988164", 10}, {"agent-pool1", "38988164", 8}, {"agent-pool-1", "38988164", 0}}
	for _, el := range data {
		vmName := fmt.Sprintf("k8s-%s-%s-%d", el.poolIdentifier, el.nameSuffix, el.agentIndex)
		poolIdentifier, nameSuffix, agentIndex, err := k8sLinuxVMNameParts(vmName)
		if poolIdentifier != el.poolIdentifier {
			t.Fatalf("incorrect poolIdentifier. expected=%s actual=%s", el.poolIdentifier, poolIdentifier)
		}
		if nameSuffix != el.nameSuffix {
			t.Fatalf("incorrect nameSuffix. expected=%s actual=%s", el.nameSuffix, nameSuffix)
		}
		if agentIndex != el.agentIndex {
			t.Fatalf("incorrect agentIndex. expected=%d actual=%d", el.agentIndex, agentIndex)
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}
}
func TestWindowsVMNameParts(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	data := []struct {
		VMName, expectedPoolPrefix, expectedOrch	string
		expectedPoolIndex, expectedAgentIndex		int
	}{{"38988k8s90312", "38988", "k8s", 3, 12}, {"4506k8s010", "4506", "k8s", 1, 0}, {"2314k8s03000001", "2314", "k8s", 3, 1}, {"2314k8s0310", "2314", "k8s", 3, 10}}
	for _, d := range data {
		poolPrefix, orch, poolIndex, agentIndex, err := windowsVMNameParts(d.VMName)
		if poolPrefix != d.expectedPoolPrefix {
			t.Fatalf("incorrect poolPrefix. expected=%s actual=%s", d.expectedPoolPrefix, poolPrefix)
		}
		if orch != d.expectedOrch {
			t.Fatalf("incorrect acs string. expected=%s actual=%s", d.expectedOrch, orch)
		}
		if poolIndex != d.expectedPoolIndex {
			t.Fatalf("incorrect poolIndex. expected=%d actual=%d", d.expectedPoolIndex, poolIndex)
		}
		if agentIndex != d.expectedAgentIndex {
			t.Fatalf("incorrect agentIndex. expected=%d actual=%d", d.expectedAgentIndex, agentIndex)
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}
}
func TestGetVMNameIndexLinux(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expectedAgentIndex := 65
	agentIndex, err := GetVMNameIndex(compute.Linux, "k8s-agentpool1-38988164-65")
	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
func TestGetVMNameIndexWindows(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	expectedAgentIndex := 20
	agentIndex, err := GetVMNameIndex(compute.Windows, "38988k8s90320")
	if agentIndex != expectedAgentIndex {
		t.Fatalf("incorrect agentIndex. expected=%d actual=%d", expectedAgentIndex, agentIndex)
	}
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
func TestIsSuccessResponse(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		name		string
		resp		*http.Response
		err		error
		expected	bool
		expectedError	error
	}{{name: "both resp and err nil should report error", expected: false, expectedError: fmt.Errorf("failed with unknown error")}, {name: "http.StatusNotFound should report error", resp: &http.Response{StatusCode: http.StatusNotFound}, expected: false, expectedError: fmt.Errorf("failed with HTTP status code %d", http.StatusNotFound)}, {name: "http.StatusInternalServerError should report error", resp: &http.Response{StatusCode: http.StatusInternalServerError}, expected: false, expectedError: fmt.Errorf("failed with HTTP status code %d", http.StatusInternalServerError)}, {name: "http.StatusOK shouldn't report error", resp: &http.Response{StatusCode: http.StatusOK}, expected: true}, {name: "non-nil response error with http.StatusOK should report error", resp: &http.Response{StatusCode: http.StatusOK}, err: fmt.Errorf("test error"), expected: false, expectedError: fmt.Errorf("test error")}, {name: "non-nil response error with http.StatusInternalServerError should report error", resp: &http.Response{StatusCode: http.StatusInternalServerError}, err: fmt.Errorf("test error"), expected: false, expectedError: fmt.Errorf("test error")}}
	for _, test := range tests {
		result, realError := isSuccessHTTPResponse(test.resp, test.err)
		assert.Equal(t, test.expected, result, "[%s] expected: %v, saw: %v", test.name, result, test.expected)
		assert.Equal(t, test.expectedError, realError, "[%s] expected: %v, saw: %v", test.name, realError, test.expectedError)
	}
}
