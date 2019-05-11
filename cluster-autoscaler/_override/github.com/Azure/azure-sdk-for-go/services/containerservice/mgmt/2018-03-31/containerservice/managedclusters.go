package containerservice

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

type ManagedClustersClient struct{ BaseClient }

func NewManagedClustersClient(subscriptionID string) ManagedClustersClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewManagedClustersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}
func NewManagedClustersClientWithBaseURI(baseURI string, subscriptionID string) ManagedClustersClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ManagedClustersClient{NewWithBaseURI(baseURI, subscriptionID)}
}
func (client ManagedClustersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, resourceName string, parameters ManagedCluster) (result ManagedClustersCreateOrUpdateFuture, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validation.Validate([]validation.Validation{{TargetValue: parameters, Constraints: []validation.Constraint{{Target: "parameters.ManagedClusterProperties", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.LinuxProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.LinuxProfile.AdminUsername", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.LinuxProfile.AdminUsername", Name: validation.Pattern, Rule: `^[a-z][a-z0-9_-]*$`, Chain: nil}}}, {Target: "parameters.ManagedClusterProperties.LinuxProfile.SSH", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.LinuxProfile.SSH.PublicKeys", Name: validation.Null, Rule: true, Chain: nil}}}}}, {Target: "parameters.ManagedClusterProperties.ServicePrincipalProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.ServicePrincipalProfile.ClientID", Name: validation.Null, Rule: true, Chain: nil}, {Target: "parameters.ManagedClusterProperties.ServicePrincipalProfile.KeyVaultSecretRef", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.ServicePrincipalProfile.KeyVaultSecretRef.VaultID", Name: validation.Null, Rule: true, Chain: nil}, {Target: "parameters.ManagedClusterProperties.ServicePrincipalProfile.KeyVaultSecretRef.SecretName", Name: validation.Null, Rule: true, Chain: nil}}}}}, {Target: "parameters.ManagedClusterProperties.NetworkProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.NetworkProfile.PodCidr", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.NetworkProfile.PodCidr", Name: validation.Pattern, Rule: `^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`, Chain: nil}}}, {Target: "parameters.ManagedClusterProperties.NetworkProfile.ServiceCidr", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.NetworkProfile.ServiceCidr", Name: validation.Pattern, Rule: `^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`, Chain: nil}}}, {Target: "parameters.ManagedClusterProperties.NetworkProfile.DNSServiceIP", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.NetworkProfile.DNSServiceIP", Name: validation.Pattern, Rule: `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, Chain: nil}}}, {Target: "parameters.ManagedClusterProperties.NetworkProfile.DockerBridgeCidr", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.NetworkProfile.DockerBridgeCidr", Name: validation.Pattern, Rule: `^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))?$`, Chain: nil}}}}}, {Target: "parameters.ManagedClusterProperties.AadProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.ManagedClusterProperties.AadProfile.ClientAppID", Name: validation.Null, Rule: true, Chain: nil}, {Target: "parameters.ManagedClusterProperties.AadProfile.ServerAppID", Name: validation.Null, Rule: true, Chain: nil}}}}}}}}); err != nil {
		return result, validation.NewError("containerservice.ManagedClustersClient", "CreateOrUpdate", err.Error())
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, resourceName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}
	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}
	return
}
func (client ManagedClustersClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, resourceName string, parameters ManagedCluster) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "resourceName": autorest.Encode("path", resourceName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsContentType("application/json; charset=utf-8"), autorest.AsPut(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParameters), autorest.WithJSON(parameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) CreateOrUpdateSender(req *http.Request) (future ManagedClustersCreateOrUpdateFuture, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	err = autorest.Respond(resp, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}
func (client ManagedClustersClient) CreateOrUpdateResponder(resp *http.Response) (result ManagedCluster, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ManagedClustersClient) Delete(ctx context.Context, resourceGroupName string, resourceName string) (result ManagedClustersDeleteFuture, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.DeletePreparer(ctx, resourceGroupName, resourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "Delete", nil, "Failure preparing request")
		return
	}
	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "Delete", result.Response(), "Failure sending request")
		return
	}
	return
}
func (client ManagedClustersClient) DeletePreparer(ctx context.Context, resourceGroupName string, resourceName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "resourceName": autorest.Encode("path", resourceName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsDelete(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) DeleteSender(req *http.Request) (future ManagedClustersDeleteFuture, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	err = autorest.Respond(resp, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}
func (client ManagedClustersClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent), autorest.ByClosing())
	result.Response = resp
	return
}
func (client ManagedClustersClient) Get(ctx context.Context, resourceGroupName string, resourceName string) (result ManagedCluster, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.GetPreparer(ctx, resourceGroupName, resourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "Get", nil, "Failure preparing request")
		return
	}
	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "Get", resp, "Failure sending request")
		return
	}
	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "Get", resp, "Failure responding to request")
	}
	return
}
func (client ManagedClustersClient) GetPreparer(ctx context.Context, resourceGroupName string, resourceName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "resourceName": autorest.Encode("path", resourceName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) GetSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ManagedClustersClient) GetResponder(resp *http.Response) (result ManagedCluster, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ManagedClustersClient) GetAccessProfile(ctx context.Context, resourceGroupName string, resourceName string, roleName string) (result ManagedClusterAccessProfile, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.GetAccessProfilePreparer(ctx, resourceGroupName, resourceName, roleName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "GetAccessProfile", nil, "Failure preparing request")
		return
	}
	resp, err := client.GetAccessProfileSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "GetAccessProfile", resp, "Failure sending request")
		return
	}
	result, err = client.GetAccessProfileResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "GetAccessProfile", resp, "Failure responding to request")
	}
	return
}
func (client ManagedClustersClient) GetAccessProfilePreparer(ctx context.Context, resourceGroupName string, resourceName string, roleName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "resourceName": autorest.Encode("path", resourceName), "roleName": autorest.Encode("path", roleName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsPost(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/accessProfiles/{roleName}/listCredential", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) GetAccessProfileSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ManagedClustersClient) GetAccessProfileResponder(resp *http.Response) (result ManagedClusterAccessProfile, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ManagedClustersClient) GetUpgradeProfile(ctx context.Context, resourceGroupName string, resourceName string) (result ManagedClusterUpgradeProfile, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.GetUpgradeProfilePreparer(ctx, resourceGroupName, resourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "GetUpgradeProfile", nil, "Failure preparing request")
		return
	}
	resp, err := client.GetUpgradeProfileSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "GetUpgradeProfile", resp, "Failure sending request")
		return
	}
	result, err = client.GetUpgradeProfileResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "GetUpgradeProfile", resp, "Failure responding to request")
	}
	return
}
func (client ManagedClustersClient) GetUpgradeProfilePreparer(ctx context.Context, resourceGroupName string, resourceName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "resourceName": autorest.Encode("path", resourceName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/upgradeProfiles/default", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) GetUpgradeProfileSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ManagedClustersClient) GetUpgradeProfileResponder(resp *http.Response) (result ManagedClusterUpgradeProfile, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ManagedClustersClient) List(ctx context.Context) (result ManagedClusterListResultPage, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "List", nil, "Failure preparing request")
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.mclr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "List", resp, "Failure sending request")
		return
	}
	result.mclr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "List", resp, "Failure responding to request")
	}
	return
}
func (client ManagedClustersClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.ContainerService/managedClusters", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) ListSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ManagedClustersClient) ListResponder(resp *http.Response) (result ManagedClusterListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ManagedClustersClient) listNextResults(lastResults ManagedClusterListResult) (result ManagedClusterListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := lastResults.managedClusterListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}
func (client ManagedClustersClient) ListComplete(ctx context.Context) (result ManagedClusterListResultIterator, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.page, err = client.List(ctx)
	return
}
func (client ManagedClustersClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result ManagedClusterListResultPage, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.fn = client.listByResourceGroupNextResults
	req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.mclr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "ListByResourceGroup", resp, "Failure sending request")
		return
	}
	result.mclr, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "ListByResourceGroup", resp, "Failure responding to request")
	}
	return
}
func (client ManagedClustersClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ManagedClustersClient) ListByResourceGroupResponder(resp *http.Response) (result ManagedClusterListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ManagedClustersClient) listByResourceGroupNextResults(lastResults ManagedClusterListResult) (result ManagedClusterListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := lastResults.managedClusterListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "listByResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
	}
	return
}
func (client ManagedClustersClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result ManagedClusterListResultIterator, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.page, err = client.ListByResourceGroup(ctx, resourceGroupName)
	return
}
func (client ManagedClustersClient) UpdateTags(ctx context.Context, resourceGroupName string, resourceName string, parameters TagsObject) (result ManagedClustersUpdateTagsFuture, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.UpdateTagsPreparer(ctx, resourceGroupName, resourceName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "UpdateTags", nil, "Failure preparing request")
		return
	}
	result, err = client.UpdateTagsSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ManagedClustersClient", "UpdateTags", result.Response(), "Failure sending request")
		return
	}
	return
}
func (client ManagedClustersClient) UpdateTagsPreparer(ctx context.Context, resourceGroupName string, resourceName string, parameters TagsObject) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "resourceName": autorest.Encode("path", resourceName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsContentType("application/json; charset=utf-8"), autorest.AsPatch(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParameters), autorest.WithJSON(parameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ManagedClustersClient) UpdateTagsSender(req *http.Request) (future ManagedClustersUpdateTagsFuture, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	err = autorest.Respond(resp, azure.WithErrorUnlessStatusCode(http.StatusOK))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}
func (client ManagedClustersClient) UpdateTagsResponder(resp *http.Response) (result ManagedCluster, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
