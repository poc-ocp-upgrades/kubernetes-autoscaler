package sdk

import (
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider/alicloud/alibaba-cloud-sdk-go/sdk/utils"
	"net/http"
	"time"
)

type Config struct {
	AutoRetry			bool			`default:"true"`
	MaxRetryTime		int				`default:"3"`
	UserAgent			string			`default:""`
	Debug				bool			`default:"false"`
	Timeout				time.Duration	`default:"10000000000"`
	HttpTransport		*http.Transport	`default:""`
	EnableAsync			bool			`default:"false"`
	MaxTaskQueueSize	int				`default:"1000"`
	GoRoutinePoolSize	int				`default:"5"`
	Scheme				string			`default:"HTTP"`
}

func NewConfig() (config *Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	config = &Config{}
	utils.InitStructWithDefaultTag(config)
	return
}
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Timeout = timeout
	return c
}
func (c *Config) WithAutoRetry(isAutoRetry bool) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.AutoRetry = isAutoRetry
	return c
}
func (c *Config) WithMaxRetryTime(maxRetryTime int) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.MaxRetryTime = maxRetryTime
	return c
}
func (c *Config) WithUserAgent(userAgent string) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.UserAgent = userAgent
	return c
}
func (c *Config) WithHttpTransport(httpTransport *http.Transport) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.HttpTransport = httpTransport
	return c
}
func (c *Config) WithEnableAsync(isEnableAsync bool) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.EnableAsync = isEnableAsync
	return c
}
func (c *Config) WithMaxTaskQueueSize(maxTaskQueueSize int) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.MaxTaskQueueSize = maxTaskQueueSize
	return c
}
func (c *Config) WithGoRoutinePoolSize(goRoutinePoolSize int) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.GoRoutinePoolSize = goRoutinePoolSize
	return c
}
func (c *Config) WithDebug(isDebug bool) *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Debug = isDebug
	return c
}
