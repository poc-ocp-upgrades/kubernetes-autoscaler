package azure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const EnvironmentFilepathName = "AZURE_ENVIRONMENT_FILEPATH"

var environments = map[string]Environment{"AZURECHINACLOUD": ChinaCloud, "AZUREGERMANCLOUD": GermanCloud, "AZUREPUBLICCLOUD": PublicCloud, "AZUREUSGOVERNMENTCLOUD": USGovernmentCloud}

type Environment struct {
	Name				string	`json:"name"`
	ManagementPortalURL		string	`json:"managementPortalURL"`
	PublishSettingsURL		string	`json:"publishSettingsURL"`
	ServiceManagementEndpoint	string	`json:"serviceManagementEndpoint"`
	ResourceManagerEndpoint		string	`json:"resourceManagerEndpoint"`
	ActiveDirectoryEndpoint		string	`json:"activeDirectoryEndpoint"`
	GalleryEndpoint			string	`json:"galleryEndpoint"`
	KeyVaultEndpoint		string	`json:"keyVaultEndpoint"`
	GraphEndpoint			string	`json:"graphEndpoint"`
	ServiceBusEndpoint		string	`json:"serviceBusEndpoint"`
	BatchManagementEndpoint		string	`json:"batchManagementEndpoint"`
	StorageEndpointSuffix		string	`json:"storageEndpointSuffix"`
	SQLDatabaseDNSSuffix		string	`json:"sqlDatabaseDNSSuffix"`
	TrafficManagerDNSSuffix		string	`json:"trafficManagerDNSSuffix"`
	KeyVaultDNSSuffix		string	`json:"keyVaultDNSSuffix"`
	ServiceBusEndpointSuffix	string	`json:"serviceBusEndpointSuffix"`
	ServiceManagementVMDNSSuffix	string	`json:"serviceManagementVMDNSSuffix"`
	ResourceManagerVMDNSSuffix	string	`json:"resourceManagerVMDNSSuffix"`
	ContainerRegistryDNSSuffix	string	`json:"containerRegistryDNSSuffix"`
	TokenAudience			string	`json:"tokenAudience"`
}

var (
	PublicCloud		= Environment{Name: "AzurePublicCloud", ManagementPortalURL: "https://manage.windowsazure.com/", PublishSettingsURL: "https://manage.windowsazure.com/publishsettings/index", ServiceManagementEndpoint: "https://management.core.windows.net/", ResourceManagerEndpoint: "https://management.azure.com/", ActiveDirectoryEndpoint: "https://login.microsoftonline.com/", GalleryEndpoint: "https://gallery.azure.com/", KeyVaultEndpoint: "https://vault.azure.net/", GraphEndpoint: "https://graph.windows.net/", ServiceBusEndpoint: "https://servicebus.windows.net/", BatchManagementEndpoint: "https://batch.core.windows.net/", StorageEndpointSuffix: "core.windows.net", SQLDatabaseDNSSuffix: "database.windows.net", TrafficManagerDNSSuffix: "trafficmanager.net", KeyVaultDNSSuffix: "vault.azure.net", ServiceBusEndpointSuffix: "servicebus.windows.net", ServiceManagementVMDNSSuffix: "cloudapp.net", ResourceManagerVMDNSSuffix: "cloudapp.azure.com", ContainerRegistryDNSSuffix: "azurecr.io", TokenAudience: "https://management.azure.com/"}
	USGovernmentCloud	= Environment{Name: "AzureUSGovernmentCloud", ManagementPortalURL: "https://manage.windowsazure.us/", PublishSettingsURL: "https://manage.windowsazure.us/publishsettings/index", ServiceManagementEndpoint: "https://management.core.usgovcloudapi.net/", ResourceManagerEndpoint: "https://management.usgovcloudapi.net/", ActiveDirectoryEndpoint: "https://login.microsoftonline.us/", GalleryEndpoint: "https://gallery.usgovcloudapi.net/", KeyVaultEndpoint: "https://vault.usgovcloudapi.net/", GraphEndpoint: "https://graph.windows.net/", ServiceBusEndpoint: "https://servicebus.usgovcloudapi.net/", BatchManagementEndpoint: "https://batch.core.usgovcloudapi.net/", StorageEndpointSuffix: "core.usgovcloudapi.net", SQLDatabaseDNSSuffix: "database.usgovcloudapi.net", TrafficManagerDNSSuffix: "usgovtrafficmanager.net", KeyVaultDNSSuffix: "vault.usgovcloudapi.net", ServiceBusEndpointSuffix: "servicebus.usgovcloudapi.net", ServiceManagementVMDNSSuffix: "usgovcloudapp.net", ResourceManagerVMDNSSuffix: "cloudapp.windowsazure.us", ContainerRegistryDNSSuffix: "azurecr.io", TokenAudience: "https://management.usgovcloudapi.net/"}
	ChinaCloud		= Environment{Name: "AzureChinaCloud", ManagementPortalURL: "https://manage.chinacloudapi.com/", PublishSettingsURL: "https://manage.chinacloudapi.com/publishsettings/index", ServiceManagementEndpoint: "https://management.core.chinacloudapi.cn/", ResourceManagerEndpoint: "https://management.chinacloudapi.cn/", ActiveDirectoryEndpoint: "https://login.chinacloudapi.cn/", GalleryEndpoint: "https://gallery.chinacloudapi.cn/", KeyVaultEndpoint: "https://vault.azure.cn/", GraphEndpoint: "https://graph.chinacloudapi.cn/", ServiceBusEndpoint: "https://servicebus.chinacloudapi.cn/", BatchManagementEndpoint: "https://batch.chinacloudapi.cn/", StorageEndpointSuffix: "core.chinacloudapi.cn", SQLDatabaseDNSSuffix: "database.chinacloudapi.cn", TrafficManagerDNSSuffix: "trafficmanager.cn", KeyVaultDNSSuffix: "vault.azure.cn", ServiceBusEndpointSuffix: "servicebus.chinacloudapi.cn", ServiceManagementVMDNSSuffix: "chinacloudapp.cn", ResourceManagerVMDNSSuffix: "cloudapp.azure.cn", ContainerRegistryDNSSuffix: "azurecr.io", TokenAudience: "https://management.chinacloudapi.cn/"}
	GermanCloud		= Environment{Name: "AzureGermanCloud", ManagementPortalURL: "http://portal.microsoftazure.de/", PublishSettingsURL: "https://manage.microsoftazure.de/publishsettings/index", ServiceManagementEndpoint: "https://management.core.cloudapi.de/", ResourceManagerEndpoint: "https://management.microsoftazure.de/", ActiveDirectoryEndpoint: "https://login.microsoftonline.de/", GalleryEndpoint: "https://gallery.cloudapi.de/", KeyVaultEndpoint: "https://vault.microsoftazure.de/", GraphEndpoint: "https://graph.cloudapi.de/", ServiceBusEndpoint: "https://servicebus.cloudapi.de/", BatchManagementEndpoint: "https://batch.cloudapi.de/", StorageEndpointSuffix: "core.cloudapi.de", SQLDatabaseDNSSuffix: "database.cloudapi.de", TrafficManagerDNSSuffix: "azuretrafficmanager.de", KeyVaultDNSSuffix: "vault.microsoftazure.de", ServiceBusEndpointSuffix: "servicebus.cloudapi.de", ServiceManagementVMDNSSuffix: "azurecloudapp.de", ResourceManagerVMDNSSuffix: "cloudapp.microsoftazure.de", ContainerRegistryDNSSuffix: "azurecr.io", TokenAudience: "https://management.microsoftazure.de/"}
)

func EnvironmentFromName(name string) (Environment, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if strings.EqualFold(name, "AZURESTACKCLOUD") {
		return EnvironmentFromFile(os.Getenv(EnvironmentFilepathName))
	}
	name = strings.ToUpper(name)
	env, ok := environments[name]
	if !ok {
		return env, fmt.Errorf("autorest/azure: There is no cloud environment matching the name %q", name)
	}
	return env, nil
}
func EnvironmentFromFile(location string) (unmarshaled Environment, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fileContents, err := ioutil.ReadFile(location)
	if err != nil {
		return
	}
	err = json.Unmarshal(fileContents, &unmarshaled)
	return
}
