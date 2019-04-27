package autorest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
	"github.com/Azure/go-autorest/logger"
	"github.com/Azure/go-autorest/version"
)

const (
	DefaultPollingDelay	= 60 * time.Second
	DefaultPollingDuration	= 15 * time.Minute
	DefaultRetryAttempts	= 3
	DefaultRetryDuration	= 30 * time.Second
)

var (
	StatusCodesForRetry = []int{http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}
)

const (
	requestFormat	= `HTTP Request Begin ===================================================
%s
===================================================== HTTP Request End
`
	responseFormat	= `HTTP Response Begin ===================================================
%s
===================================================== HTTP Response End
`
)

type Response struct {
	*http.Response `json:"-"`
}
type LoggingInspector struct{ Logger *log.Logger }

func (li LoggingInspector) WithInspection() PrepareDecorator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			var body, b bytes.Buffer
			defer r.Body.Close()
			r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &body))
			if err := r.Write(&b); err != nil {
				return nil, fmt.Errorf("Failed to write response: %v", err)
			}
			li.Logger.Printf(requestFormat, b.String())
			r.Body = ioutil.NopCloser(&body)
			return p.Prepare(r)
		})
	}
}
func (li LoggingInspector) ByInspecting() RespondDecorator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			var body, b bytes.Buffer
			defer resp.Body.Close()
			resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, &body))
			if err := resp.Write(&b); err != nil {
				return fmt.Errorf("Failed to write response: %v", err)
			}
			li.Logger.Printf(responseFormat, b.String())
			resp.Body = ioutil.NopCloser(&body)
			return r.Respond(resp)
		})
	}
}

type Client struct {
	Authorizer				Authorizer
	Sender					Sender
	RequestInspector			PrepareDecorator
	ResponseInspector			RespondDecorator
	PollingDelay				time.Duration
	PollingDuration				time.Duration
	RetryAttempts				int
	RetryDuration				time.Duration
	UserAgent				string
	Jar					http.CookieJar
	SkipResourceProviderRegistration	bool
}

func NewClientWithUserAgent(ua string) Client {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := Client{PollingDelay: DefaultPollingDelay, PollingDuration: DefaultPollingDuration, RetryAttempts: DefaultRetryAttempts, RetryDuration: DefaultRetryDuration, UserAgent: version.UserAgent()}
	c.Sender = c.sender()
	c.AddToUserAgent(ua)
	return c
}
func (c *Client) AddToUserAgent(extension string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if extension != "" {
		c.UserAgent = fmt.Sprintf("%s %s", c.UserAgent, extension)
		return nil
	}
	return fmt.Errorf("Extension was empty, User Agent stayed as %s", c.UserAgent)
}
func (c Client) Do(r *http.Request) (*http.Response, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if r.UserAgent() == "" {
		r, _ = Prepare(r, WithUserAgent(c.UserAgent))
	}
	r, err := Prepare(r, c.WithAuthorization(), c.WithInspection())
	if err != nil {
		var resp *http.Response
		if detErr, ok := err.(DetailedError); ok {
			resp = detErr.Response
		}
		return resp, NewErrorWithError(err, "autorest/Client", "Do", nil, "Preparing request failed")
	}
	logger.Instance.WriteRequest(r, logger.Filter{Header: func(k string, v []string) (bool, []string) {
		if strings.EqualFold(k, "Authorization") || strings.EqualFold(k, "Ocp-Apim-Subscription-Key") {
			v = []string{"**REDACTED**"}
		}
		return true, v
	}})
	resp, err := SendWithSender(c.sender(), r)
	logger.Instance.WriteResponse(resp, logger.Filter{})
	Respond(resp, c.ByInspecting())
	return resp, err
}
func (c Client) sender() Sender {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.Sender == nil {
		j, _ := cookiejar.New(nil)
		return &http.Client{Jar: j}
	}
	return c.Sender
}
func (c Client) WithAuthorization() PrepareDecorator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return c.authorizer().WithAuthorization()
}
func (c Client) authorizer() Authorizer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.Authorizer == nil {
		return NullAuthorizer{}
	}
	return c.Authorizer
}
func (c Client) WithInspection() PrepareDecorator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.RequestInspector == nil {
		return WithNothing()
	}
	return c.RequestInspector
}
func (c Client) ByInspecting() RespondDecorator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c.ResponseInspector == nil {
		return ByIgnoring()
	}
	return c.ResponseInspector
}
