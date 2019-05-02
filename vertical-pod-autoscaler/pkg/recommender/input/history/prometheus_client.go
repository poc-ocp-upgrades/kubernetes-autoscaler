package history

import (
 "fmt"
 "net/http"
 "net/url"
 "time"
)

var (
 numRetries = 10
 retryDelay = 3 * time.Second
)

type PrometheusClient interface {
 GetTimeseries(query string) ([]Timeseries, error)
}
type httpGetter interface {
 Get(url string) (*http.Response, error)
}
type prometheusClient struct {
 httpClient httpGetter
 address    string
}

func NewPrometheusClient(httpClient httpGetter, address string) PrometheusClient {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &prometheusClient{httpClient: httpClient, address: address}
}
func getUrlWithQuery(address, query string) (string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 url, err := url.Parse(address)
 if err != nil {
  return "", err
 }
 url.Path = "api/v1/query"
 queryValues := url.Query()
 queryValues.Set("query", query)
 url.RawQuery = queryValues.Encode()
 return url.String(), nil
}
func retry(callback func() error, attempts int, delay time.Duration) error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for i := 1; ; i++ {
  err := callback()
  if err == nil {
   return nil
  }
  if i >= attempts {
   return fmt.Errorf("tried %d times, last error: %v", attempts, err)
  }
  time.Sleep(delay)
 }
}
func (c *prometheusClient) GetTimeseries(query string) ([]Timeseries, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 url, err := getUrlWithQuery(c.address, query)
 if err != nil {
  return nil, fmt.Errorf("couldn't construct url to Prometheus: %v", err)
 }
 var resp *http.Response
 err = retry(func() error {
  resp, err = c.httpClient.Get(url)
  if err != nil {
   return fmt.Errorf("error getting data from Prometheus: %v", err)
  }
  if resp.StatusCode != http.StatusOK {
   return fmt.Errorf("bad HTTP status: %v %s", resp.StatusCode, http.StatusText(resp.StatusCode))
  }
  return nil
 }, numRetries, retryDelay)
 if err != nil {
  return nil, fmt.Errorf("Retrying GetTimeseries unsuccessful: %v", err)
 }
 return decodeTimeseriesFromResponse(resp.Body)
}
