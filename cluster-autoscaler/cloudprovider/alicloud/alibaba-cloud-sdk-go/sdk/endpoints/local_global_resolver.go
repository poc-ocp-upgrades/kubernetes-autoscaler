package endpoints

import (
	"fmt"
	"github.com/jmespath/go-jmespath"
	"strings"
)

type LocalGlobalResolver struct{}

func (resolver *LocalGlobalResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	endpointExpression := fmt.Sprintf("products[?code=='%s'].global_endpoint", strings.ToLower(param.Product))
	endpointData, err := jmespath.Search(endpointExpression, getEndpointConfigData())
	if err == nil && endpointData != nil && len(endpointData.([]interface{})) > 0 {
		endpoint = endpointData.([]interface{})[0].(string)
		support = len(endpoint) > 0
		return endpoint, support, nil
	}
	support = false
	return
}
