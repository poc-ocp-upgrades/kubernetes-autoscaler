package containerservice

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

type OperationsClient struct{ BaseClient }

func NewOperationsClient(subscriptionID string) OperationsClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewOperationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return OperationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}
func (client OperationsClient) List(ctx context.Context) (result OperationListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := client.ListPreparer(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.OperationsClient", "List", nil, "Failure preparing request")
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerservice.OperationsClient", "List", resp, "Failure sending request")
		return
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservice.OperationsClient", "List", resp, "Failure responding to request")
	}
	return
}
func (client OperationsClient) ListPreparer(ctx context.Context) (*http.Request, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const APIVersion = "2018-03-31"
	queryParameters := map[string]interface{}{"api-version": APIVersion}
	preparer := autorest.CreatePreparer(autorest.AsGet(), autorest.WithBaseURL(client.BaseURI), autorest.WithPath("/providers/Microsoft.ContainerService/operations"), autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}
func (client OperationsClient) ListSender(req *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return autorest.SendWithSender(client, req, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
}
func (client OperationsClient) ListResponder(resp *http.Response) (result OperationListResult, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = autorest.Respond(resp, client.ByInspecting(), azure.WithErrorUnlessStatusCode(http.StatusOK), autorest.ByUnmarshallingJSON(&result), autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
