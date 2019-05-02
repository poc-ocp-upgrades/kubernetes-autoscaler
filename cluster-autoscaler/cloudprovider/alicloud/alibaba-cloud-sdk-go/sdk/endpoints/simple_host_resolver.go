package endpoints

type SimpleHostResolver struct{}

func (resolver *SimpleHostResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if support = len(param.Domain) > 0; support {
  endpoint = param.Domain
 }
 return
}
