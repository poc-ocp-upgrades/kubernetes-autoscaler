package azure

import (
 "context"
 "fmt"
 "io/ioutil"
 "net/http"
 "time"
 "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
 "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
 "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
 "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
 "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
 "github.com/Azure/go-autorest/autorest"
 "github.com/Azure/go-autorest/autorest/adal"
 "github.com/Azure/go-autorest/autorest/azure"
 "k8s.io/klog"
)

type VirtualMachineScaleSetsClient interface {
 Get(ctx context.Context, resourceGroupName string, vmScaleSetName string) (result compute.VirtualMachineScaleSet, err error)
 CreateOrUpdate(ctx context.Context, resourceGroupName string, name string, parameters compute.VirtualMachineScaleSet) (resp *http.Response, err error)
 DeleteInstances(ctx context.Context, resourceGroupName string, vmScaleSetName string, vmInstanceIDs compute.VirtualMachineScaleSetVMInstanceRequiredIDs) (resp *http.Response, err error)
 List(ctx context.Context, resourceGroupName string) (result []compute.VirtualMachineScaleSet, err error)
}
type VirtualMachineScaleSetVMsClient interface {
 Get(ctx context.Context, resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVM, err error)
 List(ctx context.Context, resourceGroupName string, virtualMachineScaleSetName string, filter string, selectParameter string, expand string) (result []compute.VirtualMachineScaleSetVM, err error)
}
type VirtualMachinesClient interface {
 Get(ctx context.Context, resourceGroupName string, VMName string, expand compute.InstanceViewTypes) (result compute.VirtualMachine, err error)
 Delete(ctx context.Context, resourceGroupName string, VMName string) (resp *http.Response, err error)
 List(ctx context.Context, resourceGroupName string) (result []compute.VirtualMachine, err error)
}
type InterfacesClient interface {
 Delete(ctx context.Context, resourceGroupName string, networkInterfaceName string) (resp *http.Response, err error)
}
type DeploymentsClient interface {
 Get(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExtended, err error)
 ExportTemplate(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExportResult, err error)
 CreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.Deployment) (resp *http.Response, err error)
}
type DisksClient interface {
 Delete(ctx context.Context, resourceGroupName string, diskName string) (resp *http.Response, err error)
}
type AccountsClient interface {
 ListKeys(ctx context.Context, resourceGroupName string, accountName string) (result storage.AccountListKeysResult, err error)
}
type azVirtualMachineScaleSetsClient struct {
 client compute.VirtualMachineScaleSetsClient
}

func newAzVirtualMachineScaleSetsClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azVirtualMachineScaleSetsClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 virtualMachineScaleSetsClient := compute.NewVirtualMachineScaleSetsClient(subscriptionID)
 virtualMachineScaleSetsClient.BaseURI = endpoint
 virtualMachineScaleSetsClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 virtualMachineScaleSetsClient.PollingDelay = 5 * time.Second
 configureUserAgent(&virtualMachineScaleSetsClient.Client)
 return &azVirtualMachineScaleSetsClient{client: virtualMachineScaleSetsClient}
}
func (az *azVirtualMachineScaleSetsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, VMScaleSetName string, parameters compute.VirtualMachineScaleSet) (resp *http.Response, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachineScaleSetsClient.CreateOrUpdate(%q,%q): start", resourceGroupName, VMScaleSetName)
 defer func() {
  klog.V(10).Infof("azVirtualMachineScaleSetsClient.CreateOrUpdate(%q,%q): end", resourceGroupName, VMScaleSetName)
 }()
 future, err := az.client.CreateOrUpdate(ctx, resourceGroupName, VMScaleSetName, parameters)
 if err != nil {
  return future.Response(), err
 }
 err = future.WaitForCompletion(ctx, az.client.Client)
 return future.Response(), err
}
func (az *azVirtualMachineScaleSetsClient) Get(ctx context.Context, resourceGroupName string, VMScaleSetName string) (result compute.VirtualMachineScaleSet, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachineScaleSetsClient.Get(%q,%q): start", resourceGroupName, VMScaleSetName)
 defer func() {
  klog.V(10).Infof("azVirtualMachineScaleSetsClient.Get(%q,%q): end", resourceGroupName, VMScaleSetName)
 }()
 return az.client.Get(ctx, resourceGroupName, VMScaleSetName)
}
func (az *azVirtualMachineScaleSetsClient) List(ctx context.Context, resourceGroupName string) (result []compute.VirtualMachineScaleSet, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachineScaleSetsClient.List(%q): start", resourceGroupName)
 defer func() {
  klog.V(10).Infof("azVirtualMachineScaleSetsClient.List(%q): end", resourceGroupName)
 }()
 iterator, err := az.client.ListComplete(ctx, resourceGroupName)
 if err != nil {
  return nil, err
 }
 result = make([]compute.VirtualMachineScaleSet, 0)
 for ; iterator.NotDone(); err = iterator.Next() {
  if err != nil {
   return nil, err
  }
  result = append(result, iterator.Value())
 }
 return result, nil
}
func (az *azVirtualMachineScaleSetsClient) DeleteInstances(ctx context.Context, resourceGroupName string, vmScaleSetName string, vmInstanceIDs compute.VirtualMachineScaleSetVMInstanceRequiredIDs) (resp *http.Response, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachineScaleSetsClient.DeleteInstances(%q,%q,%v): start", resourceGroupName, vmScaleSetName, vmInstanceIDs)
 defer func() {
  klog.V(10).Infof("azVirtualMachineScaleSetsClient.DeleteInstances(%q,%q,%v): end", resourceGroupName, vmScaleSetName, vmInstanceIDs)
 }()
 future, err := az.client.DeleteInstances(ctx, resourceGroupName, vmScaleSetName, vmInstanceIDs)
 if err != nil {
  return future.Response(), err
 }
 err = future.WaitForCompletion(ctx, az.client.Client)
 return future.Response(), err
}

type azVirtualMachineScaleSetVMsClient struct {
 client compute.VirtualMachineScaleSetVMsClient
}

func newAzVirtualMachineScaleSetVMsClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azVirtualMachineScaleSetVMsClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 virtualMachineScaleSetVMsClient := compute.NewVirtualMachineScaleSetVMsClient(subscriptionID)
 virtualMachineScaleSetVMsClient.BaseURI = endpoint
 virtualMachineScaleSetVMsClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 virtualMachineScaleSetVMsClient.PollingDelay = 5 * time.Second
 configureUserAgent(&virtualMachineScaleSetVMsClient.Client)
 return &azVirtualMachineScaleSetVMsClient{client: virtualMachineScaleSetVMsClient}
}
func (az *azVirtualMachineScaleSetVMsClient) Get(ctx context.Context, resourceGroupName string, VMScaleSetName string, instanceID string) (result compute.VirtualMachineScaleSetVM, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachineScaleSetVMsClient.Get(%q,%q,%q): start", resourceGroupName, VMScaleSetName, instanceID)
 defer func() {
  klog.V(10).Infof("azVirtualMachineScaleSetVMsClient.Get(%q,%q,%q): end", resourceGroupName, VMScaleSetName, instanceID)
 }()
 return az.client.Get(ctx, resourceGroupName, VMScaleSetName, instanceID)
}
func (az *azVirtualMachineScaleSetVMsClient) List(ctx context.Context, resourceGroupName string, virtualMachineScaleSetName string, filter string, selectParameter string, expand string) (result []compute.VirtualMachineScaleSetVM, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachineScaleSetVMsClient.List(%q,%q,%q): start", resourceGroupName, virtualMachineScaleSetName, filter)
 defer func() {
  klog.V(10).Infof("azVirtualMachineScaleSetVMsClient.List(%q,%q,%q): end", resourceGroupName, virtualMachineScaleSetName, filter)
 }()
 iterator, err := az.client.ListComplete(ctx, resourceGroupName, virtualMachineScaleSetName, filter, selectParameter, expand)
 if err != nil {
  return nil, err
 }
 result = make([]compute.VirtualMachineScaleSetVM, 0)
 for ; iterator.NotDone(); err = iterator.Next() {
  if err != nil {
   return nil, err
  }
  result = append(result, iterator.Value())
 }
 return result, nil
}

type azVirtualMachinesClient struct{ client compute.VirtualMachinesClient }

func newAzVirtualMachinesClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azVirtualMachinesClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 virtualMachinesClient := compute.NewVirtualMachinesClient(subscriptionID)
 virtualMachinesClient.BaseURI = endpoint
 virtualMachinesClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 virtualMachinesClient.PollingDelay = 5 * time.Second
 configureUserAgent(&virtualMachinesClient.Client)
 return &azVirtualMachinesClient{client: virtualMachinesClient}
}
func (az *azVirtualMachinesClient) Get(ctx context.Context, resourceGroupName string, VMName string, expand compute.InstanceViewTypes) (result compute.VirtualMachine, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachinesClient.Get(%q,%q,%q): start", resourceGroupName, VMName, expand)
 defer func() {
  klog.V(10).Infof("azVirtualMachinesClient.Get(%q,%q,%q): end", resourceGroupName, VMName, expand)
 }()
 return az.client.Get(ctx, resourceGroupName, VMName, expand)
}
func (az *azVirtualMachinesClient) Delete(ctx context.Context, resourceGroupName string, VMName string) (resp *http.Response, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachinesClient.Delete(%q,%q): start", resourceGroupName, VMName)
 defer func() {
  klog.V(10).Infof("azVirtualMachinesClient.Delete(%q,%q): end", resourceGroupName, VMName)
 }()
 future, err := az.client.Delete(ctx, resourceGroupName, VMName)
 if err != nil {
  return future.Response(), err
 }
 err = future.WaitForCompletion(ctx, az.client.Client)
 return future.Response(), err
}
func (az *azVirtualMachinesClient) List(ctx context.Context, resourceGroupName string) (result []compute.VirtualMachine, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azVirtualMachinesClient.List(%q): start", resourceGroupName)
 defer func() {
  klog.V(10).Infof("azVirtualMachinesClient.List(%q): end", resourceGroupName)
 }()
 iterator, err := az.client.ListComplete(ctx, resourceGroupName)
 if err != nil {
  return nil, err
 }
 result = make([]compute.VirtualMachine, 0)
 for ; iterator.NotDone(); err = iterator.Next() {
  if err != nil {
   return nil, err
  }
  result = append(result, iterator.Value())
 }
 return result, nil
}

type azInterfacesClient struct{ client network.InterfacesClient }

func newAzInterfacesClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azInterfacesClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 interfacesClient := network.NewInterfacesClient(subscriptionID)
 interfacesClient.BaseURI = endpoint
 interfacesClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 interfacesClient.PollingDelay = 5 * time.Second
 configureUserAgent(&interfacesClient.Client)
 return &azInterfacesClient{client: interfacesClient}
}
func (az *azInterfacesClient) Delete(ctx context.Context, resourceGroupName string, networkInterfaceName string) (resp *http.Response, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azInterfacesClient.Delete(%q,%q): start", resourceGroupName, networkInterfaceName)
 defer func() {
  klog.V(10).Infof("azInterfacesClient.Delete(%q,%q): end", resourceGroupName, networkInterfaceName)
 }()
 future, err := az.client.Delete(ctx, resourceGroupName, networkInterfaceName)
 if err != nil {
  return future.Response(), err
 }
 err = future.WaitForCompletion(ctx, az.client.Client)
 return future.Response(), err
}

type azDeploymentsClient struct{ client resources.DeploymentsClient }

func newAzDeploymentsClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azDeploymentsClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 deploymentsClient := resources.NewDeploymentsClient(subscriptionID)
 deploymentsClient.BaseURI = endpoint
 deploymentsClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 deploymentsClient.PollingDelay = 5 * time.Second
 configureUserAgent(&deploymentsClient.Client)
 return &azDeploymentsClient{client: deploymentsClient}
}
func (az *azDeploymentsClient) Get(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExtended, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azDeploymentsClient.Get(%q,%q): start", resourceGroupName, deploymentName)
 defer func() {
  klog.V(10).Infof("azDeploymentsClient.Get(%q,%q): end", resourceGroupName, deploymentName)
 }()
 return az.client.Get(ctx, resourceGroupName, deploymentName)
}
func (az *azDeploymentsClient) ExportTemplate(ctx context.Context, resourceGroupName string, deploymentName string) (result resources.DeploymentExportResult, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azDeploymentsClient.ExportTemplate(%q,%q): start", resourceGroupName, deploymentName)
 defer func() {
  klog.V(10).Infof("azDeploymentsClient.ExportTemplate(%q,%q): end", resourceGroupName, deploymentName)
 }()
 return az.client.ExportTemplate(ctx, resourceGroupName, deploymentName)
}
func (az *azDeploymentsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, deploymentName string, parameters resources.Deployment) (resp *http.Response, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azDeploymentsClient.CreateOrUpdate(%q,%q): start", resourceGroupName, deploymentName)
 defer func() {
  klog.V(10).Infof("azDeploymentsClient.CreateOrUpdate(%q,%q): end", resourceGroupName, deploymentName)
 }()
 future, err := az.client.CreateOrUpdate(ctx, resourceGroupName, deploymentName, parameters)
 if err != nil {
  return future.Response(), err
 }
 err = future.WaitForCompletion(ctx, az.client.Client)
 return future.Response(), err
}

type azDisksClient struct{ client compute.DisksClient }

func newAzDisksClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azDisksClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 disksClient := compute.NewDisksClient(subscriptionID)
 disksClient.BaseURI = endpoint
 disksClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 disksClient.PollingDelay = 5 * time.Second
 configureUserAgent(&disksClient.Client)
 return &azDisksClient{client: disksClient}
}
func (az *azDisksClient) Delete(ctx context.Context, resourceGroupName string, diskName string) (resp *http.Response, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azDisksClient.Delete(%q,%q): start", resourceGroupName, diskName)
 defer func() {
  klog.V(10).Infof("azDisksClient.Delete(%q,%q): end", resourceGroupName, diskName)
 }()
 future, err := az.client.Delete(ctx, resourceGroupName, diskName)
 if err != nil {
  return future.Response(), err
 }
 err = future.WaitForCompletion(ctx, az.client.Client)
 return future.Response(), err
}

type azAccountsClient struct{ client storage.AccountsClient }

func newAzAccountsClient(subscriptionID, endpoint string, servicePrincipalToken *adal.ServicePrincipalToken) *azAccountsClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 accountsClient := storage.NewAccountsClient(subscriptionID)
 accountsClient.BaseURI = endpoint
 accountsClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipalToken)
 accountsClient.PollingDelay = 5 * time.Second
 configureUserAgent(&accountsClient.Client)
 return &azAccountsClient{client: accountsClient}
}
func (az *azAccountsClient) ListKeys(ctx context.Context, resourceGroupName string, accountName string) (result storage.AccountListKeysResult, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 klog.V(10).Infof("azAccountsClient.ListKeys(%q,%q): start", resourceGroupName, accountName)
 defer func() {
  klog.V(10).Infof("azAccountsClient.ListKeys(%q,%q): end", resourceGroupName, accountName)
 }()
 return az.client.ListKeys(ctx, resourceGroupName, accountName)
}

type azClient struct {
 virtualMachineScaleSetsClient   VirtualMachineScaleSetsClient
 virtualMachineScaleSetVMsClient VirtualMachineScaleSetVMsClient
 virtualMachinesClient           VirtualMachinesClient
 deploymentsClient               DeploymentsClient
 interfacesClient                InterfacesClient
 disksClient                     DisksClient
 storageAccountsClient           AccountsClient
 containerServicesClient         containerservice.ContainerServicesClient
 managedContainerServicesClient  containerservice.ManagedClustersClient
}

func newServicePrincipalTokenFromCredentials(config *Config, env *azure.Environment) (*adal.ServicePrincipalToken, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, config.TenantID)
 if err != nil {
  return nil, fmt.Errorf("creating the OAuth config: %v", err)
 }
 if config.UseManagedIdentityExtension {
  klog.V(2).Infoln("azure: using managed identity extension to retrieve access token")
  msiEndpoint, err := adal.GetMSIVMEndpoint()
  if err != nil {
   return nil, fmt.Errorf("Getting the managed service identity endpoint: %v", err)
  }
  return adal.NewServicePrincipalTokenFromMSI(msiEndpoint, env.ServiceManagementEndpoint)
 }
 if len(config.AADClientSecret) > 0 {
  klog.V(2).Infoln("azure: using client_id+client_secret to retrieve access token")
  return adal.NewServicePrincipalToken(*oauthConfig, config.AADClientID, config.AADClientSecret, env.ServiceManagementEndpoint)
 }
 if len(config.AADClientCertPath) > 0 && len(config.AADClientCertPassword) > 0 {
  klog.V(2).Infoln("azure: using jwt client_assertion (client_cert+client_private_key) to retrieve access token")
  certData, err := ioutil.ReadFile(config.AADClientCertPath)
  if err != nil {
   return nil, fmt.Errorf("reading the client certificate from file %s: %v", config.AADClientCertPath, err)
  }
  certificate, privateKey, err := decodePkcs12(certData, config.AADClientCertPassword)
  if err != nil {
   return nil, fmt.Errorf("decoding the client certificate: %v", err)
  }
  return adal.NewServicePrincipalTokenFromCertificate(*oauthConfig, config.AADClientID, certificate, privateKey, env.ServiceManagementEndpoint)
 }
 return nil, fmt.Errorf("No credentials provided for AAD application %s", config.AADClientID)
}
func newAzClient(cfg *Config, env *azure.Environment) (*azClient, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 spt, err := newServicePrincipalTokenFromCredentials(cfg, env)
 if err != nil {
  return nil, err
 }
 scaleSetsClient := newAzVirtualMachineScaleSetsClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created scale set client with authorizer: %v", scaleSetsClient)
 scaleSetVMsClient := newAzVirtualMachineScaleSetVMsClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created scale set vm client with authorizer: %v", scaleSetVMsClient)
 virtualMachinesClient := newAzVirtualMachinesClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created vm client with authorizer: %v", virtualMachinesClient)
 deploymentsClient := newAzDeploymentsClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created deployments client with authorizer: %v", deploymentsClient)
 interfacesClient := newAzInterfacesClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created interfaces client with authorizer: %v", interfacesClient)
 storageAccountsClient := newAzAccountsClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created storage accounts client with authorizer: %v", storageAccountsClient)
 disksClient := newAzDisksClient(cfg.SubscriptionID, env.ResourceManagerEndpoint, spt)
 klog.V(5).Infof("Created disks client with authorizer: %v", disksClient)
 containerServicesClient := containerservice.NewContainerServicesClient(cfg.SubscriptionID)
 containerServicesClient.BaseURI = env.ResourceManagerEndpoint
 containerServicesClient.Authorizer = autorest.NewBearerAuthorizer(spt)
 containerServicesClient.PollingDelay = 5 * time.Second
 containerServicesClient.Sender = autorest.CreateSender()
 klog.V(5).Infof("Created Container services client with authorizer: %v", containerServicesClient)
 managedContainerServicesClient := containerservice.NewManagedClustersClient(cfg.SubscriptionID)
 managedContainerServicesClient.BaseURI = env.ResourceManagerEndpoint
 managedContainerServicesClient.Authorizer = autorest.NewBearerAuthorizer(spt)
 managedContainerServicesClient.PollingDelay = 5 * time.Second
 managedContainerServicesClient.Sender = autorest.CreateSender()
 klog.V(5).Infof("Created Managed Container services client with authorizer: %v", managedContainerServicesClient)
 return &azClient{disksClient: disksClient, interfacesClient: interfacesClient, virtualMachineScaleSetsClient: scaleSetsClient, virtualMachineScaleSetVMsClient: scaleSetVMsClient, deploymentsClient: deploymentsClient, virtualMachinesClient: virtualMachinesClient, storageAccountsClient: storageAccountsClient, containerServicesClient: containerServicesClient, managedContainerServicesClient: managedContainerServicesClient}, nil
}
