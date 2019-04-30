package containerservice

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

type ContainerServicesClient struct{ BaseClient }

func NewContainerServicesClient(subscriptionID string) ContainerServicesClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewContainerServicesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}
func NewContainerServicesClientWithBaseURI(baseURI string, subscriptionID string) ContainerServicesClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ContainerServicesClient{NewWithBaseURI(baseURI, subscriptionID)}
}
func (client ContainerServicesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, containerServiceName string, parameters ContainerService) (result ContainerServicesCreateOrUpdateFutureType, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validation.Validate([]validation.Validation{{TargetValue: parameters, Constraints: []validation.Constraint{{Target: "parameters.Properties", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.Properties.OrchestratorProfile", Name: validation.Null, Rule: true, Chain: nil}, {Target: "parameters.Properties.CustomProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.Properties.CustomProfile.Orchestrator", Name: validation.Null, Rule: true, Chain: nil}}}, {Target: "parameters.Properties.ServicePrincipalProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.Properties.ServicePrincipalProfile.ClientID", Name: validation.Null, Rule: true, Chain: nil}, {Target: "parameters.Properties.ServicePrincipalProfile.KeyVaultSecretRef", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.Properties.ServicePrincipalProfile.KeyVaultSecretRef.VaultID", Name: validation.Null, Rule: true, Chain: nil}, {Target: "parameters.Properties.ServicePrincipalProfile.KeyVaultSecretRef.SecretName", Name: validation.Null, Rule: true, Chain: nil}}}}}, {Target: "parameters.Properties.MasterProfile", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.Properties.MasterProfile.DNSPrefix", Name: validation.Null, Rule: true, Chain: nil}}}, {Target: "parameters.Properties.WindowsProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.Properties.WindowsProfile.AdminUsername", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.Properties.WindowsProfile.AdminUsername", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]+([._]?[a-zA-Z0-9]+)*$`, Chain: nil}}}, {Target: "parameters.Properties.WindowsProfile.AdminPassword", Name: validation.Null, Rule: true, Chain: nil}}}, {Target: "parameters.Properties.LinuxProfile", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.Properties.LinuxProfile.AdminUsername", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.Properties.LinuxProfile.AdminUsername", Name: validation.Pattern, Rule: `^[a-z][a-z0-9_-]*$`, Chain: nil}}}, {Target: "parameters.Properties.LinuxProfile.SSH", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.Properties.LinuxProfile.SSH.PublicKeys", Name: validation.Null, Rule: true, Chain: nil}}}}}, {Target: "parameters.Properties.DiagnosticsProfile", Name: validation.Null, Rule: false, Chain: []validation.Constraint{{Target: "parameters.Properties.DiagnosticsProfile.VMDiagnostics", Name: validation.Null, Rule: true, Chain: []validation.Constraint{{Target: "parameters.Properties.DiagnosticsProfile.VMDiagnostics.Enabled", Name: validation.Null, Rule: true, Chain: nil}}}}}}}}}}); err != nil {
		return result, validation.NewError("containerservice.ContainerServicesClient", "CreateOrUpdate", err.Error())
	}
	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, containerServiceName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}
	result, err = client.CreateOrUpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "CreateOrUpdate", result.Response(), "Failure sending request")
		return
	}
	return
}
func (client ContainerServicesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, containerServiceName string, parameters ContainerService) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"containerServiceName": autorest.Encode("path", containerServiceName), "resourceGroupName": autorest.Encode("path", resourceGroupName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2017-07-01"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsContentType("application/json; charset=utf-8"), autorest.AsPut(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices/{containerServiceName}", pathParameters), autorest.WithJSON(parameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ContainerServicesClient) CreateOrUpdateSender(req *http.Request) (future ContainerServicesCreateOrUpdateFutureType, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	err = autorest.Respond(resp, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted))
	if err != nil {
		return
	}
	future.Future, err = azure.NewFutureFromResponse(resp)
	return
}
func (client ContainerServicesClient) CreateOrUpdateResponder(resp *http.Response) (result ContainerService, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated, http.StatusAccepted), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ContainerServicesClient) Delete(ctx context.Context, resourceGroupName string, containerServiceName string) (result ContainerServicesDeleteFutureType, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.DeletePreparer(ctx, resourceGroupName, containerServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Delete", nil, "Failure preparing request")
		return
	}
	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Delete", result.Response(), "Failure sending request")
		return
	}
	return
}
func (client ContainerServicesClient) DeletePreparer(ctx context.Context, resourceGroupName string, containerServiceName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"containerServiceName": autorest.Encode("path", containerServiceName), "resourceGroupName": autorest.Encode("path", resourceGroupName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2017-07-01"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsDelete(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices/{containerServiceName}", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ContainerServicesClient) DeleteSender(req *http.Request) (future ContainerServicesDeleteFutureType, err error) {
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
func (client ContainerServicesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent), autorest.ByClosing())
	result.Response = resp
	return
}
func (client ContainerServicesClient) Get(ctx context.Context, resourceGroupName string, containerServiceName string) (result ContainerService, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.GetPreparer(ctx, resourceGroupName, containerServiceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Get", nil, "Failure preparing request")
		return
	}
	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Get", resp, "Failure sending request")
		return
	}
	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "Get", resp, "Failure responding to request")
	}
	return
}
func (client ContainerServicesClient) GetPreparer(ctx context.Context, resourceGroupName string, containerServiceName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"containerServiceName": autorest.Encode("path", containerServiceName), "resourceGroupName": autorest.Encode("path", resourceGroupName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2017-07-01"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices/{containerServiceName}", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ContainerServicesClient) GetSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ContainerServicesClient) GetResponder(resp *http.Response) (result ContainerService, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ContainerServicesClient) List(ctx context.Context) (result ListResultPage, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", nil, "Failure preparing request")
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.lr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", resp, "Failure sending request")
		return
	}
	result.lr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "List", resp, "Failure responding to request")
	}
	return
}
func (client ContainerServicesClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2017-07-01"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.ContainerService/containerServices", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ContainerServicesClient) ListSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ContainerServicesClient) ListResponder(resp *http.Response) (result ListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ContainerServicesClient) listNextResults(lastResults ListResult) (result ListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := lastResults.listResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}
func (client ContainerServicesClient) ListComplete(ctx context.Context) (result ListResultIterator, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.page, err = client.List(ctx)
	return
}
func (client ContainerServicesClient) ListByResourceGroup(ctx context.Context, resourceGroupName string) (result ListResultPage, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.fn = client.listByResourceGroupNextResults
	req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.lr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", resp, "Failure sending request")
		return
	}
	result.lr, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListByResourceGroup", resp, "Failure responding to request")
	}
	return
}
func (client ContainerServicesClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"resourceGroupName": autorest.Encode("path", resourceGroupName), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2017-07-01"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/containerServices", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ContainerServicesClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ContainerServicesClient) ListByResourceGroupResponder(resp *http.Response) (result ListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
func (client ContainerServicesClient) listByResourceGroupNextResults(lastResults ListResult) (result ListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := lastResults.listResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "listByResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
	}
	return
}
func (client ContainerServicesClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result ListResultIterator, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	result.page, err = client.ListByResourceGroup(ctx, resourceGroupName)
	return
}
func (client ContainerServicesClient) ListOrchestrators(ctx context.Context, location string, resourceType string) (result OrchestratorVersionProfileListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.ListOrchestratorsPreparer(ctx, location, resourceType)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListOrchestrators", nil, "Failure preparing request")
		return
	}
	resp, err := client.ListOrchestratorsSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListOrchestrators", resp, "Failure sending request")
		return
	}
	result, err = client.ListOrchestratorsResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.ContainerServicesClient", "ListOrchestrators", resp, "Failure responding to request")
	}
	return
}
func (client ContainerServicesClient) ListOrchestratorsPreparer(ctx context.Context, location string, resourceType string) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pathParameters := map[string]interface{}{"location": autorest.Encode("path", location), "subscriptionId": autorest.Encode("path", client.SubscriptionID)}
	const APIVersion = "2017-09-30"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	if len(resourceType) > 0 {
		queryParameters["resource-type"] = autorest.Encode("query", resourceType)
	}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.ContainerService/locations/{location}/orchestrators", pathParameters), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client ContainerServicesClient) ListOrchestratorsSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
}
func (client ContainerServicesClient) ListOrchestratorsResponder(resp *http.Response) (result OrchestratorVersionProfileListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
