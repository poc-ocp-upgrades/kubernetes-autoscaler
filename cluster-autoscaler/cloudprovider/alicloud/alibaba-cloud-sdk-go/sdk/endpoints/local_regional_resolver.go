package endpoints

import (
 "fmt"
 "github.com/jmespath/go-jmespath"
 "strings"
)

type LocalRegionalResolver struct{}

func (resolver *LocalRegionalResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 regionalExpression := fmt.Sprintf("products[?code=='%s'].regional_endpoints", strings.ToLower(param.Product))
 regionalData, err := jmespath.Search(regionalExpression, getEndpointConfigData())
 if err == nil && regionalData != nil && len(regionalData.([]interface{})) > 0 {
  endpointExpression := fmt.Sprintf("[0][?region=='%s'].endpoint", strings.ToLower(param.RegionId))
  endpointData, err := jmespath.Search(endpointExpression, regionalData)
  if err == nil && endpointData != nil && len(endpointData.([]interface{})) > 0 {
   endpoint = endpointData.([]interface{})[0].(string)
   support = len(endpoint) > 0
   return endpoint, support, nil
  }
 }
 support = false
 return
}
