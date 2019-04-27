package adal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	logPrefix = "autorest/adal/devicetoken:"
)

var (
	ErrDeviceGeneric		= fmt.Errorf("%s Error while retrieving OAuth token: Unknown Error", logPrefix)
	ErrDeviceAccessDenied		= fmt.Errorf("%s Error while retrieving OAuth token: Access Denied", logPrefix)
	ErrDeviceAuthorizationPending	= fmt.Errorf("%s Error while retrieving OAuth token: Authorization Pending", logPrefix)
	ErrDeviceCodeExpired		= fmt.Errorf("%s Error while retrieving OAuth token: Code Expired", logPrefix)
	ErrDeviceSlowDown		= fmt.Errorf("%s Error while retrieving OAuth token: Slow Down", logPrefix)
	ErrDeviceCodeEmpty		= fmt.Errorf("%s Error while retrieving device code: Device Code Empty", logPrefix)
	ErrOAuthTokenEmpty		= fmt.Errorf("%s Error while retrieving OAuth token: Token Empty", logPrefix)
	errCodeSendingFails		= "Error occurred while sending request for Device Authorization Code"
	errCodeHandlingFails		= "Error occurred while handling response from the Device Endpoint"
	errTokenSendingFails		= "Error occurred while sending request with device code for a token"
	errTokenHandlingFails		= "Error occurred while handling response from the Token Endpoint (during device flow)"
	errStatusNotOK			= "Error HTTP status != 200"
)

type DeviceCode struct {
	DeviceCode	*string	`json:"device_code,omitempty"`
	UserCode	*string	`json:"user_code,omitempty"`
	VerificationURL	*string	`json:"verification_url,omitempty"`
	ExpiresIn	*int64	`json:"expires_in,string,omitempty"`
	Interval	*int64	`json:"interval,string,omitempty"`
	Message		*string	`json:"message"`
	Resource	string
	OAuthConfig	OAuthConfig
	ClientID	string
}
type TokenError struct {
	Error			*string	`json:"error,omitempty"`
	ErrorCodes		[]int	`json:"error_codes,omitempty"`
	ErrorDescription	*string	`json:"error_description,omitempty"`
	Timestamp		*string	`json:"timestamp,omitempty"`
	TraceID			*string	`json:"trace_id,omitempty"`
}
type deviceToken struct {
	Token
	TokenError
}

func InitiateDeviceAuth(sender Sender, oauthConfig OAuthConfig, clientID, resource string) (*DeviceCode, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := url.Values{"client_id": []string{clientID}, "resource": []string{resource}}
	s := v.Encode()
	body := ioutil.NopCloser(strings.NewReader(s))
	req, err := http.NewRequest(http.MethodPost, oauthConfig.DeviceCodeEndpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errCodeSendingFails, err.Error())
	}
	req.ContentLength = int64(len(s))
	req.Header.Set(contentType, mimeTypeFormPost)
	resp, err := sender.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errCodeSendingFails, err.Error())
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errCodeHandlingFails, err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errCodeHandlingFails, errStatusNotOK)
	}
	if len(strings.Trim(string(rb), " ")) == 0 {
		return nil, ErrDeviceCodeEmpty
	}
	var code DeviceCode
	err = json.Unmarshal(rb, &code)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errCodeHandlingFails, err.Error())
	}
	code.ClientID = clientID
	code.Resource = resource
	code.OAuthConfig = oauthConfig
	return &code, nil
}
func CheckForUserCompletion(sender Sender, code *DeviceCode) (*Token, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v := url.Values{"client_id": []string{code.ClientID}, "code": []string{*code.DeviceCode}, "grant_type": []string{OAuthGrantTypeDeviceCode}, "resource": []string{code.Resource}}
	s := v.Encode()
	body := ioutil.NopCloser(strings.NewReader(s))
	req, err := http.NewRequest(http.MethodPost, code.OAuthConfig.TokenEndpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errTokenSendingFails, err.Error())
	}
	req.ContentLength = int64(len(s))
	req.Header.Set(contentType, mimeTypeFormPost)
	resp, err := sender.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errTokenSendingFails, err.Error())
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errTokenHandlingFails, err.Error())
	}
	if resp.StatusCode != http.StatusOK && len(strings.Trim(string(rb), " ")) == 0 {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errTokenHandlingFails, errStatusNotOK)
	}
	if len(strings.Trim(string(rb), " ")) == 0 {
		return nil, ErrOAuthTokenEmpty
	}
	var token deviceToken
	err = json.Unmarshal(rb, &token)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %s", logPrefix, errTokenHandlingFails, err.Error())
	}
	if token.Error == nil {
		return &token.Token, nil
	}
	switch *token.Error {
	case "authorization_pending":
		return nil, ErrDeviceAuthorizationPending
	case "slow_down":
		return nil, ErrDeviceSlowDown
	case "access_denied":
		return nil, ErrDeviceAccessDenied
	case "code_expired":
		return nil, ErrDeviceCodeExpired
	default:
		return nil, ErrDeviceGeneric
	}
}
func WaitForUserCompletion(sender Sender, code *DeviceCode) (*Token, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	intervalDuration := time.Duration(*code.Interval) * time.Second
	waitDuration := intervalDuration
	for {
		token, err := CheckForUserCompletion(sender, code)
		if err == nil {
			return token, nil
		}
		switch err {
		case ErrDeviceSlowDown:
			waitDuration += waitDuration
		case ErrDeviceAuthorizationPending:
		default:
			return nil, err
		}
		if waitDuration > (intervalDuration * 3) {
			return nil, fmt.Errorf("%s Error waiting for user to complete device flow. Server told us to slow_down too much", logPrefix)
		}
		time.Sleep(waitDuration)
	}
}
