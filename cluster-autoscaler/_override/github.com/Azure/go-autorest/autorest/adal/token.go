package adal

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/version"
	"github.com/dgrijalva/jwt-go"
)

const (
	defaultRefresh					= 5 * time.Minute
	OAuthGrantTypeDeviceCode		= "device_code"
	OAuthGrantTypeClientCredentials	= "client_credentials"
	OAuthGrantTypeUserPass			= "password"
	OAuthGrantTypeRefreshToken		= "refresh_token"
	OAuthGrantTypeAuthorizationCode	= "authorization_code"
	metadataHeader					= "Metadata"
	msiEndpoint						= "http://169.254.169.254/metadata/identity/oauth2/token"
	defaultMaxMSIRefreshAttempts	= 5
)

type OAuthTokenProvider interface{ OAuthToken() string }
type TokenRefreshError interface {
	error
	Response() *http.Response
}
type Refresher interface {
	Refresh() error
	RefreshExchange(resource string) error
	EnsureFresh() error
}
type RefresherWithContext interface {
	RefreshWithContext(ctx context.Context) error
	RefreshExchangeWithContext(ctx context.Context, resource string) error
	EnsureFreshWithContext(ctx context.Context) error
}
type TokenRefreshCallback func(Token) error
type Token struct {
	AccessToken		string		`json:"access_token"`
	RefreshToken	string		`json:"refresh_token"`
	ExpiresIn		json.Number	`json:"expires_in"`
	ExpiresOn		json.Number	`json:"expires_on"`
	NotBefore		json.Number	`json:"not_before"`
	Resource		string		`json:"resource"`
	Type			string		`json:"token_type"`
}

func newToken() Token {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Token{ExpiresIn: "0", ExpiresOn: "0", NotBefore: "0"}
}
func (t Token) IsZero() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t == Token{}
}
func (t Token) Expires() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s, err := t.ExpiresOn.Float64()
	if err != nil {
		s = -3600
	}
	expiration := date.NewUnixTimeFromSeconds(s)
	return time.Time(expiration).UTC()
}
func (t Token) IsExpired() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.WillExpireIn(0)
}
func (t Token) WillExpireIn(d time.Duration) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return !t.Expires().After(time.Now().Add(d))
}
func (t *Token) OAuthToken() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.AccessToken
}

type ServicePrincipalSecret interface {
	SetAuthenticationValues(spt *ServicePrincipalToken, values *url.Values) error
}
type ServicePrincipalNoSecret struct{}

func (noSecret *ServicePrincipalNoSecret) SetAuthenticationValues(spt *ServicePrincipalToken, v *url.Values) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Errorf("Manually created ServicePrincipalToken does not contain secret material to retrieve a new access token")
}
func (noSecret ServicePrincipalNoSecret) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type tokenType struct {
		Type string `json:"type"`
	}
	return json.Marshal(tokenType{Type: "ServicePrincipalNoSecret"})
}

type ServicePrincipalTokenSecret struct {
	ClientSecret string `json:"value"`
}

func (tokenSecret *ServicePrincipalTokenSecret) SetAuthenticationValues(spt *ServicePrincipalToken, v *url.Values) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v.Set("client_secret", tokenSecret.ClientSecret)
	return nil
}
func (tokenSecret ServicePrincipalTokenSecret) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type tokenType struct {
		Type	string	`json:"type"`
		Value	string	`json:"value"`
	}
	return json.Marshal(tokenType{Type: "ServicePrincipalTokenSecret", Value: tokenSecret.ClientSecret})
}

type ServicePrincipalCertificateSecret struct {
	Certificate	*x509.Certificate
	PrivateKey	*rsa.PrivateKey
}

func (secret *ServicePrincipalCertificateSecret) SignJwt(spt *ServicePrincipalToken) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hasher := sha1.New()
	_, err := hasher.Write(secret.Certificate.Raw)
	if err != nil {
		return "", err
	}
	thumbprint := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	jti := make([]byte, 20)
	_, err = rand.Read(jti)
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodRS256)
	token.Header["x5t"] = thumbprint
	x5c := []string{base64.StdEncoding.EncodeToString(secret.Certificate.Raw)}
	token.Header["x5c"] = x5c
	token.Claims = jwt.MapClaims{"aud": spt.inner.OauthConfig.TokenEndpoint.String(), "iss": spt.inner.ClientID, "sub": spt.inner.ClientID, "jti": base64.URLEncoding.EncodeToString(jti), "nbf": time.Now().Unix(), "exp": time.Now().Add(time.Hour * 24).Unix()}
	signedString, err := token.SignedString(secret.PrivateKey)
	return signedString, err
}
func (secret *ServicePrincipalCertificateSecret) SetAuthenticationValues(spt *ServicePrincipalToken, v *url.Values) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	jwt, err := secret.SignJwt(spt)
	if err != nil {
		return err
	}
	v.Set("client_assertion", jwt)
	v.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	return nil
}
func (secret ServicePrincipalCertificateSecret) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, errors.New("marshalling ServicePrincipalCertificateSecret is not supported")
}

type ServicePrincipalMSISecret struct{}

func (msiSecret *ServicePrincipalMSISecret) SetAuthenticationValues(spt *ServicePrincipalToken, v *url.Values) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (msiSecret ServicePrincipalMSISecret) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, errors.New("marshalling ServicePrincipalMSISecret is not supported")
}

type ServicePrincipalUsernamePasswordSecret struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

func (secret *ServicePrincipalUsernamePasswordSecret) SetAuthenticationValues(spt *ServicePrincipalToken, v *url.Values) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v.Set("username", secret.Username)
	v.Set("password", secret.Password)
	return nil
}
func (secret ServicePrincipalUsernamePasswordSecret) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type tokenType struct {
		Type		string	`json:"type"`
		Username	string	`json:"username"`
		Password	string	`json:"password"`
	}
	return json.Marshal(tokenType{Type: "ServicePrincipalUsernamePasswordSecret", Username: secret.Username, Password: secret.Password})
}

type ServicePrincipalAuthorizationCodeSecret struct {
	ClientSecret		string	`json:"value"`
	AuthorizationCode	string	`json:"authCode"`
	RedirectURI			string	`json:"redirect"`
}

func (secret *ServicePrincipalAuthorizationCodeSecret) SetAuthenticationValues(spt *ServicePrincipalToken, v *url.Values) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	v.Set("code", secret.AuthorizationCode)
	v.Set("client_secret", secret.ClientSecret)
	v.Set("redirect_uri", secret.RedirectURI)
	return nil
}
func (secret ServicePrincipalAuthorizationCodeSecret) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	type tokenType struct {
		Type		string	`json:"type"`
		Value		string	`json:"value"`
		AuthCode	string	`json:"authCode"`
		Redirect	string	`json:"redirect"`
	}
	return json.Marshal(tokenType{Type: "ServicePrincipalAuthorizationCodeSecret", Value: secret.ClientSecret, AuthCode: secret.AuthorizationCode, Redirect: secret.RedirectURI})
}

type ServicePrincipalToken struct {
	inner					servicePrincipalToken
	refreshLock				*sync.RWMutex
	sender					Sender
	refreshCallbacks		[]TokenRefreshCallback
	MaxMSIRefreshAttempts	int
}

func (spt ServicePrincipalToken) MarshalTokenJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.Marshal(spt.inner.Token)
}
func (spt *ServicePrincipalToken) SetRefreshCallbacks(callbacks []TokenRefreshCallback) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.refreshCallbacks = callbacks
}
func (spt ServicePrincipalToken) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return json.Marshal(spt.inner)
}
func (spt *ServicePrincipalToken) UnmarshalJSON(data []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	raw := map[string]interface{}{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	secret := raw["secret"].(map[string]interface{})
	switch secret["type"] {
	case "ServicePrincipalNoSecret":
		spt.inner.Secret = &ServicePrincipalNoSecret{}
	case "ServicePrincipalTokenSecret":
		spt.inner.Secret = &ServicePrincipalTokenSecret{}
	case "ServicePrincipalCertificateSecret":
		return errors.New("unmarshalling ServicePrincipalCertificateSecret is not supported")
	case "ServicePrincipalMSISecret":
		return errors.New("unmarshalling ServicePrincipalMSISecret is not supported")
	case "ServicePrincipalUsernamePasswordSecret":
		spt.inner.Secret = &ServicePrincipalUsernamePasswordSecret{}
	case "ServicePrincipalAuthorizationCodeSecret":
		spt.inner.Secret = &ServicePrincipalAuthorizationCodeSecret{}
	default:
		return fmt.Errorf("unrecognized token type '%s'", secret["type"])
	}
	err = json.Unmarshal(data, &spt.inner)
	if err != nil {
		return err
	}
	spt.refreshLock = &sync.RWMutex{}
	spt.sender = &http.Client{}
	return nil
}

type servicePrincipalToken struct {
	Token			Token					`json:"token"`
	Secret			ServicePrincipalSecret	`json:"secret"`
	OauthConfig		OAuthConfig				`json:"oauth"`
	ClientID		string					`json:"clientID"`
	Resource		string					`json:"resource"`
	AutoRefresh		bool					`json:"autoRefresh"`
	RefreshWithin	time.Duration			`json:"refreshWithin"`
}

func validateOAuthConfig(oac OAuthConfig) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if oac.IsZero() {
		return fmt.Errorf("parameter 'oauthConfig' cannot be zero-initialized")
	}
	return nil
}
func NewServicePrincipalTokenWithSecret(oauthConfig OAuthConfig, id string, resource string, secret ServicePrincipalSecret, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(id, "id"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("parameter 'secret' cannot be nil")
	}
	spt := &ServicePrincipalToken{inner: servicePrincipalToken{Token: newToken(), OauthConfig: oauthConfig, Secret: secret, ClientID: id, Resource: resource, AutoRefresh: true, RefreshWithin: defaultRefresh}, refreshLock: &sync.RWMutex{}, sender: &http.Client{}, refreshCallbacks: callbacks}
	return spt, nil
}
func NewServicePrincipalTokenFromManualToken(oauthConfig OAuthConfig, clientID string, resource string, token Token, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientID, "clientID"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	if token.IsZero() {
		return nil, fmt.Errorf("parameter 'token' cannot be zero-initialized")
	}
	spt, err := NewServicePrincipalTokenWithSecret(oauthConfig, clientID, resource, &ServicePrincipalNoSecret{}, callbacks...)
	if err != nil {
		return nil, err
	}
	spt.inner.Token = token
	return spt, nil
}
func NewServicePrincipalTokenFromManualTokenSecret(oauthConfig OAuthConfig, clientID string, resource string, token Token, secret ServicePrincipalSecret, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientID, "clientID"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("parameter 'secret' cannot be nil")
	}
	if token.IsZero() {
		return nil, fmt.Errorf("parameter 'token' cannot be zero-initialized")
	}
	spt, err := NewServicePrincipalTokenWithSecret(oauthConfig, clientID, resource, secret, callbacks...)
	if err != nil {
		return nil, err
	}
	spt.inner.Token = token
	return spt, nil
}
func NewServicePrincipalToken(oauthConfig OAuthConfig, clientID string, secret string, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientID, "clientID"); err != nil {
		return nil, err
	}
	if err := validateStringParam(secret, "secret"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	return NewServicePrincipalTokenWithSecret(oauthConfig, clientID, resource, &ServicePrincipalTokenSecret{ClientSecret: secret}, callbacks...)
}
func NewServicePrincipalTokenFromCertificate(oauthConfig OAuthConfig, clientID string, certificate *x509.Certificate, privateKey *rsa.PrivateKey, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientID, "clientID"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	if certificate == nil {
		return nil, fmt.Errorf("parameter 'certificate' cannot be nil")
	}
	if privateKey == nil {
		return nil, fmt.Errorf("parameter 'privateKey' cannot be nil")
	}
	return NewServicePrincipalTokenWithSecret(oauthConfig, clientID, resource, &ServicePrincipalCertificateSecret{PrivateKey: privateKey, Certificate: certificate}, callbacks...)
}
func NewServicePrincipalTokenFromUsernamePassword(oauthConfig OAuthConfig, clientID string, username string, password string, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientID, "clientID"); err != nil {
		return nil, err
	}
	if err := validateStringParam(username, "username"); err != nil {
		return nil, err
	}
	if err := validateStringParam(password, "password"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	return NewServicePrincipalTokenWithSecret(oauthConfig, clientID, resource, &ServicePrincipalUsernamePasswordSecret{Username: username, Password: password}, callbacks...)
}
func NewServicePrincipalTokenFromAuthorizationCode(oauthConfig OAuthConfig, clientID string, clientSecret string, authorizationCode string, redirectURI string, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateOAuthConfig(oauthConfig); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientID, "clientID"); err != nil {
		return nil, err
	}
	if err := validateStringParam(clientSecret, "clientSecret"); err != nil {
		return nil, err
	}
	if err := validateStringParam(authorizationCode, "authorizationCode"); err != nil {
		return nil, err
	}
	if err := validateStringParam(redirectURI, "redirectURI"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	return NewServicePrincipalTokenWithSecret(oauthConfig, clientID, resource, &ServicePrincipalAuthorizationCodeSecret{ClientSecret: clientSecret, AuthorizationCode: authorizationCode, RedirectURI: redirectURI}, callbacks...)
}
func GetMSIVMEndpoint() (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return msiEndpoint, nil
}
func NewServicePrincipalTokenFromMSI(msiEndpoint, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newServicePrincipalTokenFromMSI(msiEndpoint, resource, nil, callbacks...)
}
func NewServicePrincipalTokenFromMSIWithUserAssignedID(msiEndpoint, resource string, userAssignedID string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return newServicePrincipalTokenFromMSI(msiEndpoint, resource, &userAssignedID, callbacks...)
}
func newServicePrincipalTokenFromMSI(msiEndpoint, resource string, userAssignedID *string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := validateStringParam(msiEndpoint, "msiEndpoint"); err != nil {
		return nil, err
	}
	if err := validateStringParam(resource, "resource"); err != nil {
		return nil, err
	}
	if userAssignedID != nil {
		if err := validateStringParam(*userAssignedID, "userAssignedID"); err != nil {
			return nil, err
		}
	}
	msiEndpointURL, err := url.Parse(msiEndpoint)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("resource", resource)
	v.Set("api-version", "2018-02-01")
	if userAssignedID != nil {
		v.Set("client_id", *userAssignedID)
	}
	msiEndpointURL.RawQuery = v.Encode()
	spt := &ServicePrincipalToken{inner: servicePrincipalToken{Token: newToken(), OauthConfig: OAuthConfig{TokenEndpoint: *msiEndpointURL}, Secret: &ServicePrincipalMSISecret{}, Resource: resource, AutoRefresh: true, RefreshWithin: defaultRefresh}, refreshLock: &sync.RWMutex{}, sender: &http.Client{}, refreshCallbacks: callbacks, MaxMSIRefreshAttempts: defaultMaxMSIRefreshAttempts}
	if userAssignedID != nil {
		spt.inner.ClientID = *userAssignedID
	}
	return spt, nil
}

type tokenRefreshError struct {
	message	string
	resp	*http.Response
}

func (tre tokenRefreshError) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tre.message
}
func (tre tokenRefreshError) Response() *http.Response {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tre.resp
}
func newTokenRefreshError(message string, resp *http.Response) TokenRefreshError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return tokenRefreshError{message: message, resp: resp}
}
func (spt *ServicePrincipalToken) EnsureFresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return spt.EnsureFreshWithContext(context.Background())
}
func (spt *ServicePrincipalToken) EnsureFreshWithContext(ctx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if spt.inner.AutoRefresh && spt.inner.Token.WillExpireIn(spt.inner.RefreshWithin) {
		spt.refreshLock.Lock()
		defer spt.refreshLock.Unlock()
		if spt.inner.Token.WillExpireIn(spt.inner.RefreshWithin) {
			return spt.refreshInternal(ctx, spt.inner.Resource)
		}
	}
	return nil
}
func (spt *ServicePrincipalToken) InvokeRefreshCallbacks(token Token) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if spt.refreshCallbacks != nil {
		for _, callback := range spt.refreshCallbacks {
			err := callback(spt.inner.Token)
			if err != nil {
				return fmt.Errorf("adal: TokenRefreshCallback handler failed. Error = '%v'", err)
			}
		}
	}
	return nil
}
func (spt *ServicePrincipalToken) Refresh() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return spt.RefreshWithContext(context.Background())
}
func (spt *ServicePrincipalToken) RefreshWithContext(ctx context.Context) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.refreshLock.Lock()
	defer spt.refreshLock.Unlock()
	return spt.refreshInternal(ctx, spt.inner.Resource)
}
func (spt *ServicePrincipalToken) RefreshExchange(resource string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return spt.RefreshExchangeWithContext(context.Background(), resource)
}
func (spt *ServicePrincipalToken) RefreshExchangeWithContext(ctx context.Context, resource string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.refreshLock.Lock()
	defer spt.refreshLock.Unlock()
	return spt.refreshInternal(ctx, resource)
}
func (spt *ServicePrincipalToken) getGrantType() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch spt.inner.Secret.(type) {
	case *ServicePrincipalUsernamePasswordSecret:
		return OAuthGrantTypeUserPass
	case *ServicePrincipalAuthorizationCodeSecret:
		return OAuthGrantTypeAuthorizationCode
	default:
		return OAuthGrantTypeClientCredentials
	}
}
func isIMDS(u url.URL) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	imds, err := url.Parse(msiEndpoint)
	if err != nil {
		return false
	}
	return u.Host == imds.Host && u.Path == imds.Path
}
func (spt *ServicePrincipalToken) refreshInternal(ctx context.Context, resource string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	req, err := http.NewRequest(http.MethodPost, spt.inner.OauthConfig.TokenEndpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("adal: Failed to build the refresh request. Error = '%v'", err)
	}
	req.Header.Add("User-Agent", version.UserAgent())
	req = req.WithContext(ctx)
	if !isIMDS(spt.inner.OauthConfig.TokenEndpoint) {
		v := url.Values{}
		v.Set("client_id", spt.inner.ClientID)
		v.Set("resource", resource)
		if spt.inner.Token.RefreshToken != "" {
			v.Set("grant_type", OAuthGrantTypeRefreshToken)
			v.Set("refresh_token", spt.inner.Token.RefreshToken)
			if spt.getGrantType() == OAuthGrantTypeAuthorizationCode {
				err := spt.inner.Secret.SetAuthenticationValues(spt, &v)
				if err != nil {
					return err
				}
			}
		} else {
			v.Set("grant_type", spt.getGrantType())
			err := spt.inner.Secret.SetAuthenticationValues(spt, &v)
			if err != nil {
				return err
			}
		}
		s := v.Encode()
		body := ioutil.NopCloser(strings.NewReader(s))
		req.ContentLength = int64(len(s))
		req.Header.Set(contentType, mimeTypeFormPost)
		req.Body = body
	}
	if _, ok := spt.inner.Secret.(*ServicePrincipalMSISecret); ok {
		req.Method = http.MethodGet
		req.Header.Set(metadataHeader, "true")
	}
	var resp *http.Response
	if isIMDS(spt.inner.OauthConfig.TokenEndpoint) {
		resp, err = retryForIMDS(spt.sender, req, spt.MaxMSIRefreshAttempts)
	} else {
		resp, err = spt.sender.Do(req)
	}
	if err != nil {
		return newTokenRefreshError(fmt.Sprintf("adal: Failed to execute the refresh request. Error = '%v'", err), nil)
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return newTokenRefreshError(fmt.Sprintf("adal: Refresh request failed. Status Code = '%d'. Failed reading response body: %v", resp.StatusCode, err), resp)
		}
		return newTokenRefreshError(fmt.Sprintf("adal: Refresh request failed. Status Code = '%d'. Response body: %s", resp.StatusCode, string(rb)), resp)
	}
	if err != nil {
		return fmt.Errorf("adal: Failed to read a new service principal token during refresh. Error = '%v'", err)
	}
	if len(strings.Trim(string(rb), " ")) == 0 {
		return fmt.Errorf("adal: Empty service principal token received during refresh")
	}
	var token Token
	err = json.Unmarshal(rb, &token)
	if err != nil {
		return fmt.Errorf("adal: Failed to unmarshal the service principal token during refresh. Error = '%v' JSON = '%s'", err, string(rb))
	}
	spt.inner.Token = token
	return spt.InvokeRefreshCallbacks(token)
}
func retryForIMDS(sender Sender, req *http.Request, maxAttempts int) (resp *http.Response, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	retries := []int{http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}
	retries = append(retries, http.StatusNotFound, http.StatusGone, http.StatusNotImplemented, http.StatusHTTPVersionNotSupported, http.StatusVariantAlsoNegotiates, http.StatusInsufficientStorage, http.StatusLoopDetected, http.StatusNotExtended, http.StatusNetworkAuthenticationRequired)
	const maxDelay time.Duration = 60 * time.Second
	attempt := 0
	delay := time.Duration(0)
	for attempt < maxAttempts {
		resp, err = sender.Do(req)
		if (err != nil && !isTemporaryNetworkError(err)) || resp == nil || resp.StatusCode == http.StatusOK || !containsInt(retries, resp.StatusCode) {
			return
		}
		attempt++
		delay += (time.Duration(math.Pow(2, float64(attempt))) * time.Second)
		if delay > maxDelay {
			delay = maxDelay
		}
		select {
		case <-time.After(delay):
		case <-req.Context().Done():
			err = req.Context().Err()
			return
		}
	}
	return
}
func isTemporaryNetworkError(err error) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if netErr, ok := err.(net.Error); !ok || (ok && netErr.Temporary()) {
		return true
	}
	return false
}
func containsInt(ints []int, n int) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, i := range ints {
		if i == n {
			return true
		}
	}
	return false
}
func (spt *ServicePrincipalToken) SetAutoRefresh(autoRefresh bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.inner.AutoRefresh = autoRefresh
}
func (spt *ServicePrincipalToken) SetRefreshWithin(d time.Duration) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.inner.RefreshWithin = d
	return
}
func (spt *ServicePrincipalToken) SetSender(s Sender) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.sender = s
}
func (spt *ServicePrincipalToken) OAuthToken() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.refreshLock.RLock()
	defer spt.refreshLock.RUnlock()
	return spt.inner.Token.OAuthToken()
}
func (spt *ServicePrincipalToken) Token() Token {
	_logClusterCodePath()
	defer _logClusterCodePath()
	spt.refreshLock.RLock()
	defer spt.refreshLock.RUnlock()
	return spt.inner.Token
}
