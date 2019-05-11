package azure

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	azStorage "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"golang.org/x/crypto/pkcs12"
	"k8s.io/klog"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/client-go/pkg/version"
)

const (
	customDataFieldName						= "customData"
	dependsOnFieldName						= "dependsOn"
	hardwareProfileFieldName				= "hardwareProfile"
	imageReferenceFieldName					= "imageReference"
	nameFieldName							= "name"
	osProfileFieldName						= "osProfile"
	propertiesFieldName						= "properties"
	resourcesFieldName						= "resources"
	storageProfileFieldName					= "storageProfile"
	typeFieldName							= "type"
	vmSizeFieldName							= "vmSize"
	nsgResourceType							= "Microsoft.Network/networkSecurityGroups"
	rtResourceType							= "Microsoft.Network/routeTables"
	vmResourceType							= "Microsoft.Compute/virtualMachines"
	vmExtensionType							= "Microsoft.Compute/virtualMachines/extensions"
	nsgID									= "nsgID"
	rtID									= "routeTableID"
	k8sLinuxVMNamingFormat					= "^[0-9a-zA-Z]{3}-(.+)-([0-9a-fA-F]{8})-{0,2}([0-9]+)$"
	k8sLinuxVMAgentPoolNameIndex			= 1
	k8sLinuxVMAgentClusterIDIndex			= 2
	k8sLinuxVMAgentIndexArrayIndex			= 3
	k8sWindowsOldVMNamingFormat				= "^([a-fA-F0-9]{5})([0-9a-zA-Z]{3})([9])([a-zA-Z0-9]{3,5})$"
	k8sWindowsVMNamingFormat				= "^([a-fA-F0-9]{4})([0-9a-zA-Z]{3})([0-9]{3,8})$"
	k8sWindowsVMAgentPoolPrefixIndex		= 1
	k8sWindowsVMAgentOrchestratorNameIndex	= 2
	k8sWindowsVMAgentPoolInfoIndex			= 3
)

var (
	vmnameLinuxRegexp		= regexp.MustCompile(k8sLinuxVMNamingFormat)
	vmnameWindowsRegexp		= regexp.MustCompile(k8sWindowsVMNamingFormat)
	oldvmnameWindowsRegexp	= regexp.MustCompile(k8sWindowsOldVMNamingFormat)
)

type AzUtil struct{ manager *AzureManager }

func (util *AzUtil) DeleteBlob(accountName, vhdContainer, vhdBlob string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := getContextWithCancel()
	defer cancel()
	storageKeysResult, err := util.manager.azClient.storageAccountsClient.ListKeys(ctx, util.manager.config.ResourceGroup, accountName)
	if err != nil {
		return err
	}
	keys := *storageKeysResult.Keys
	client, err := azStorage.NewBasicClientOnSovereignCloud(accountName, to.String(keys[0].Value), util.manager.env)
	if err != nil {
		return err
	}
	bs := client.GetBlobService()
	containerRef := bs.GetContainerReference(vhdContainer)
	blobRef := containerRef.GetBlobReference(vhdBlob)
	return blobRef.Delete(&azStorage.DeleteBlobOptions{})
}
func (util *AzUtil) DeleteVirtualMachine(rg string, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ctx, cancel := getContextWithCancel()
	defer cancel()
	vm, err := util.manager.azClient.virtualMachinesClient.Get(ctx, rg, name, "")
	if err != nil {
		if exists, _ := checkResourceExistsFromError(err); !exists {
			klog.V(2).Infof("VirtualMachine %s/%s has already been removed", rg, name)
			return nil
		}
		klog.Errorf("failed to get VM: %s/%s: %s", rg, name, err.Error())
		return err
	}
	vhd := vm.VirtualMachineProperties.StorageProfile.OsDisk.Vhd
	managedDisk := vm.VirtualMachineProperties.StorageProfile.OsDisk.ManagedDisk
	if vhd == nil && managedDisk == nil {
		klog.Errorf("failed to get a valid os disk URI for VM: %s/%s", rg, name)
		return fmt.Errorf("os disk does not have a VHD URI")
	}
	osDiskName := vm.VirtualMachineProperties.StorageProfile.OsDisk.Name
	var nicName string
	nicID := (*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0].ID
	if nicID == nil {
		klog.Warningf("NIC ID is not set for VM (%s/%s)", rg, name)
	} else {
		nicName, err = resourceName(*nicID)
		if err != nil {
			return err
		}
		klog.Infof("found nic name for VM (%s/%s): %s", rg, name, nicName)
	}
	klog.Infof("deleting VM: %s/%s", rg, name)
	deleteCtx, deleteCancel := getContextWithCancel()
	defer deleteCancel()
	klog.Infof("waiting for VirtualMachine deletion: %s/%s", rg, name)
	_, err = util.manager.azClient.virtualMachinesClient.Delete(deleteCtx, rg, name)
	_, realErr := checkResourceExistsFromError(err)
	if realErr != nil {
		return realErr
	}
	klog.V(2).Infof("VirtualMachine %s/%s removed", rg, name)
	if len(nicName) > 0 {
		klog.Infof("deleting nic: %s/%s", rg, nicName)
		interfaceCtx, interfaceCancel := getContextWithCancel()
		defer interfaceCancel()
		klog.Infof("waiting for nic deletion: %s/%s", rg, nicName)
		_, nicErr := util.manager.azClient.interfacesClient.Delete(interfaceCtx, rg, nicName)
		_, realErr := checkResourceExistsFromError(nicErr)
		if realErr != nil {
			return realErr
		}
		klog.V(2).Infof("interface %s/%s removed", rg, nicName)
	}
	if vhd != nil {
		accountName, vhdContainer, vhdBlob, err := splitBlobURI(*vhd.URI)
		if err != nil {
			return err
		}
		klog.Infof("found os disk storage reference: %s %s %s", accountName, vhdContainer, vhdBlob)
		klog.Infof("deleting blob: %s/%s", vhdContainer, vhdBlob)
		if err = util.DeleteBlob(accountName, vhdContainer, vhdBlob); err != nil {
			_, realErr := checkResourceExistsFromError(err)
			if realErr != nil {
				return realErr
			}
			klog.V(2).Infof("Blob %s/%s removed", rg, vhdBlob)
		}
	} else if managedDisk != nil {
		if osDiskName == nil {
			klog.Warningf("osDisk is not set for VM %s/%s", rg, name)
		} else {
			klog.Infof("deleting managed disk: %s/%s", rg, *osDiskName)
			disksCtx, disksCancel := getContextWithCancel()
			defer disksCancel()
			_, diskErr := util.manager.azClient.disksClient.Delete(disksCtx, rg, *osDiskName)
			_, realErr := checkResourceExistsFromError(diskErr)
			if realErr != nil {
				return realErr
			}
			klog.V(2).Infof("disk %s/%s removed", rg, *osDiskName)
		}
	}
	return nil
}
func decodePkcs12(pkcs []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding the PKCS#12 client certificate: %v", err)
	}
	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain a RSA private key")
	}
	return certificate, rsaPrivateKey, nil
}
func configureUserAgent(client *autorest.Client) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	k8sVersion := version.Get().GitVersion
	client.UserAgent = fmt.Sprintf("%s; cluster-autoscaler/%s", client.UserAgent, k8sVersion)
}
func normalizeForK8sVMASScalingUp(templateMap map[string]interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := normalizeMasterResourcesForScaling(templateMap); err != nil {
		return err
	}
	rtIndex := -1
	nsgIndex := -1
	resources := templateMap[resourcesFieldName].([]interface{})
	for index, resource := range resources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			klog.Warning("Template improperly formatted for resource")
			continue
		}
		resourceType, ok := resourceMap[typeFieldName].(string)
		if ok && resourceType == nsgResourceType {
			if nsgIndex != -1 {
				err := fmt.Errorf("Found 2 resources with type %s in the template. There should only be 1", nsgResourceType)
				klog.Errorf(err.Error())
				return err
			}
			nsgIndex = index
		}
		if ok && resourceType == rtResourceType {
			if rtIndex != -1 {
				err := fmt.Errorf("Found 2 resources with type %s in the template. There should only be 1", rtResourceType)
				klog.Warningf(err.Error())
				return err
			}
			rtIndex = index
		}
		dependencies, ok := resourceMap[dependsOnFieldName].([]interface{})
		if !ok {
			continue
		}
		for dIndex := len(dependencies) - 1; dIndex >= 0; dIndex-- {
			dependency := dependencies[dIndex].(string)
			if strings.Contains(dependency, nsgResourceType) || strings.Contains(dependency, nsgID) || strings.Contains(dependency, rtResourceType) || strings.Contains(dependency, rtID) {
				dependencies = append(dependencies[:dIndex], dependencies[dIndex+1:]...)
			}
		}
		if len(dependencies) > 0 {
			resourceMap[dependsOnFieldName] = dependencies
		} else {
			delete(resourceMap, dependsOnFieldName)
		}
	}
	indexesToRemove := []int{}
	if nsgIndex == -1 {
		err := fmt.Errorf("Found no resources with type %s in the template. There should have been 1", nsgResourceType)
		klog.Errorf(err.Error())
		return err
	}
	if rtIndex == -1 {
		klog.Infof("Found no resources with type %s in the template.", rtResourceType)
	} else {
		indexesToRemove = append(indexesToRemove, rtIndex)
	}
	indexesToRemove = append(indexesToRemove, nsgIndex)
	templateMap[resourcesFieldName] = removeIndexesFromArray(resources, indexesToRemove)
	return nil
}
func removeIndexesFromArray(array []interface{}, indexes []int) []interface{} {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sort.Sort(sort.Reverse(sort.IntSlice(indexes)))
	for _, index := range indexes {
		array = append(array[:index], array[index+1:]...)
	}
	return array
}
func normalizeMasterResourcesForScaling(templateMap map[string]interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resources := templateMap[resourcesFieldName].([]interface{})
	indexesToRemove := []int{}
	for index, resource := range resources {
		resourceMap, ok := resource.(map[string]interface{})
		if !ok {
			klog.Warning("Template improperly formatted")
			continue
		}
		resourceType, ok := resourceMap[typeFieldName].(string)
		if !ok || resourceType != vmResourceType {
			resourceName, ok := resourceMap[nameFieldName].(string)
			if !ok {
				klog.Warning("Template improperly formatted")
				continue
			}
			if strings.Contains(resourceName, "variables('masterVMNamePrefix')") && resourceType == vmExtensionType {
				indexesToRemove = append(indexesToRemove, index)
			}
			continue
		}
		resourceName, ok := resourceMap[nameFieldName].(string)
		if !ok {
			klog.Warning("Template improperly formatted")
			continue
		}
		if !strings.Contains(resourceName, "variables('masterVMNamePrefix')") {
			continue
		}
		resourceProperties, ok := resourceMap[propertiesFieldName].(map[string]interface{})
		if !ok {
			klog.Warning("Template improperly formatted")
			continue
		}
		hardwareProfile, ok := resourceProperties[hardwareProfileFieldName].(map[string]interface{})
		if !ok {
			klog.Warning("Template improperly formatted")
			continue
		}
		if hardwareProfile[vmSizeFieldName] != nil {
			delete(hardwareProfile, vmSizeFieldName)
		}
		if !removeCustomData(resourceProperties) || !removeImageReference(resourceProperties) {
			continue
		}
	}
	templateMap[resourcesFieldName] = removeIndexesFromArray(resources, indexesToRemove)
	return nil
}
func removeCustomData(resourceProperties map[string]interface{}) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	osProfile, ok := resourceProperties[osProfileFieldName].(map[string]interface{})
	if !ok {
		klog.Warning("Template improperly formatted")
		return ok
	}
	if osProfile[customDataFieldName] != nil {
		delete(osProfile, customDataFieldName)
	}
	return ok
}
func removeImageReference(resourceProperties map[string]interface{}) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	storageProfile, ok := resourceProperties[storageProfileFieldName].(map[string]interface{})
	if !ok {
		klog.Warningf("Template improperly formatted. Could not find: %s", storageProfileFieldName)
		return ok
	}
	if storageProfile[imageReferenceFieldName] != nil {
		delete(storageProfile, imageReferenceFieldName)
	}
	return ok
}
func resourceName(ID string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	parts := strings.Split(ID, "/")
	name := parts[len(parts)-1]
	if len(name) == 0 {
		return "", fmt.Errorf("resource name was missing from identifier")
	}
	return name, nil
}
func splitBlobURI(URI string) (string, string, string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	uri, err := url.Parse(URI)
	if err != nil {
		return "", "", "", err
	}
	accountName := strings.Split(uri.Host, ".")[0]
	urlParts := strings.Split(uri.Path, "/")
	containerName := urlParts[1]
	blobPath := strings.Join(urlParts[2:], "/")
	return accountName, containerName, blobPath, nil
}
func k8sLinuxVMNameParts(vmName string) (poolIdentifier, nameSuffix string, agentIndex int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vmNameParts := vmnameLinuxRegexp.FindStringSubmatch(vmName)
	if len(vmNameParts) != 4 {
		return "", "", -1, fmt.Errorf("resource name was missing from identifier")
	}
	vmNum, err := strconv.Atoi(vmNameParts[k8sLinuxVMAgentIndexArrayIndex])
	if err != nil {
		return "", "", -1, fmt.Errorf("Error parsing VM Name: %v", err)
	}
	return vmNameParts[k8sLinuxVMAgentPoolNameIndex], vmNameParts[k8sLinuxVMAgentClusterIDIndex], vmNum, nil
}
func windowsVMNameParts(vmName string) (poolPrefix string, orch string, poolIndex int, agentIndex int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var poolInfo string
	vmNameParts := oldvmnameWindowsRegexp.FindStringSubmatch(vmName)
	if len(vmNameParts) != 5 {
		vmNameParts = vmnameWindowsRegexp.FindStringSubmatch(vmName)
		if len(vmNameParts) != 4 {
			return "", "", -1, -1, fmt.Errorf("resource name was missing from identifier")
		}
		poolInfo = vmNameParts[3]
	} else {
		poolInfo = vmNameParts[4]
	}
	poolPrefix = vmNameParts[1]
	orch = vmNameParts[2]
	poolIndex, err = strconv.Atoi(poolInfo[:2])
	if err != nil {
		return "", "", -1, -1, fmt.Errorf("error parsing VM Name: %v", err)
	}
	agentIndex, err = strconv.Atoi(poolInfo[2:])
	if err != nil {
		return "", "", -1, -1, fmt.Errorf("error parsing VM Name: %v", err)
	}
	return poolPrefix, orch, poolIndex, agentIndex, nil
}
func GetVMNameIndex(osType compute.OperatingSystemTypes, vmName string) (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var agentIndex int
	var err error
	if osType == compute.Linux {
		_, _, agentIndex, err = k8sLinuxVMNameParts(vmName)
		if err != nil {
			return 0, err
		}
	} else if osType == compute.Windows {
		_, _, _, agentIndex, err = windowsVMNameParts(vmName)
		if err != nil {
			return 0, err
		}
	}
	return agentIndex, nil
}
func matchDiscoveryConfig(labels map[string]*string, configs []cloudprovider.LabelAutoDiscoveryConfig) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(configs) == 0 {
		return false
	}
	for _, c := range configs {
		if len(c.Selector) == 0 {
			return false
		}
		for k, v := range c.Selector {
			value, ok := labels[k]
			if !ok {
				return false
			}
			if len(v) > 0 {
				if value == nil || *value != v {
					return false
				}
			}
		}
	}
	return true
}
func validateConfig(cfg *Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if cfg.ResourceGroup == "" {
		return fmt.Errorf("resource group not set")
	}
	if cfg.SubscriptionID == "" {
		return fmt.Errorf("subscription ID not set")
	}
	if cfg.UseManagedIdentityExtension {
		return nil
	}
	if cfg.TenantID == "" {
		return fmt.Errorf("tenant ID not set")
	}
	if cfg.AADClientID == "" {
		return fmt.Errorf("ARM Client ID not set")
	}
	if cfg.VMType == vmTypeStandard {
		if cfg.Deployment == "" {
			return fmt.Errorf("deployment not set")
		}
		if len(cfg.DeploymentParameters) == 0 {
			return fmt.Errorf("deploymentParameters not set")
		}
	}
	if cfg.VMType == vmTypeACS || cfg.VMType == vmTypeAKS {
		if cfg.ClusterName == "" {
			return fmt.Errorf("Cluster name not set for type %+v", cfg.VMType)
		}
	}
	return nil
}
func getLastSegment(ID string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	parts := strings.Split(strings.TrimSpace(ID), "/")
	name := parts[len(parts)-1]
	if len(name) == 0 {
		return "", fmt.Errorf("identifier '/' not found in resource name %q", ID)
	}
	return name, nil
}
func readDeploymentParameters(paramFilePath string) (map[string]interface{}, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	contents, err := ioutil.ReadFile(paramFilePath)
	if err != nil {
		klog.Errorf("Failed to read deployment parameters from file %q: %v", paramFilePath, err)
		return nil, err
	}
	deploymentParameters := make(map[string]interface{})
	if err := json.Unmarshal(contents, &deploymentParameters); err != nil {
		klog.Errorf("Failed to unmarshal deployment parameters from file %q: %v", paramFilePath, err)
		return nil, err
	}
	if v, ok := deploymentParameters["parameters"]; ok {
		return v.(map[string]interface{}), nil
	}
	return nil, fmt.Errorf("failed to get deployment parameters from file %s", paramFilePath)
}
func getContextWithCancel() (context.Context, context.CancelFunc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return context.WithCancel(context.Background())
}
func checkResourceExistsFromError(err error) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err == nil {
		return true, nil
	}
	v, ok := err.(autorest.DetailedError)
	if !ok {
		return false, err
	}
	if v.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, v
}
func isSuccessHTTPResponse(resp *http.Response, err error) (isSuccess bool, realError error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err != nil {
		return false, err
	}
	if resp != nil {
		if 199 < resp.StatusCode && resp.StatusCode < 300 {
			return true, nil
		}
		return false, fmt.Errorf("failed with HTTP status code %d", resp.StatusCode)
	}
	return false, fmt.Errorf("failed with unknown error")
}
