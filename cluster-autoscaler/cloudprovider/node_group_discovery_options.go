package cloudprovider

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	autoDiscovererTypeMIG		= "mig"
	autoDiscovererTypeASG		= "asg"
	autoDiscovererTypeLabel		= "label"
	migAutoDiscovererKeyPrefix	= "namePrefix"
	migAutoDiscovererKeyMinNodes	= "min"
	migAutoDiscovererKeyMaxNodes	= "max"
	asgAutoDiscovererKeyTag		= "tag"
)

var validMIGAutoDiscovererKeys = strings.Join([]string{migAutoDiscovererKeyPrefix, migAutoDiscovererKeyMinNodes, migAutoDiscovererKeyMaxNodes}, ", ")

type NodeGroupDiscoveryOptions struct {
	NodeGroupSpecs			[]string
	NodeGroupAutoDiscoverySpecs	[]string
}

func (o NodeGroupDiscoveryOptions) StaticDiscoverySpecified() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(o.NodeGroupSpecs) > 0
}
func (o NodeGroupDiscoveryOptions) AutoDiscoverySpecified() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return len(o.NodeGroupAutoDiscoverySpecs) > 0
}
func (o NodeGroupDiscoveryOptions) DiscoverySpecified() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return o.StaticDiscoverySpecified() || o.AutoDiscoverySpecified()
}
func (o NodeGroupDiscoveryOptions) ParseMIGAutoDiscoverySpecs() ([]MIGAutoDiscoveryConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfgs := make([]MIGAutoDiscoveryConfig, len(o.NodeGroupAutoDiscoverySpecs))
	var err error
	for i, spec := range o.NodeGroupAutoDiscoverySpecs {
		cfgs[i], err = parseMIGAutoDiscoverySpec(spec)
		if err != nil {
			return nil, err
		}
	}
	return cfgs, nil
}
func (o NodeGroupDiscoveryOptions) ParseASGAutoDiscoverySpecs() ([]ASGAutoDiscoveryConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfgs := make([]ASGAutoDiscoveryConfig, len(o.NodeGroupAutoDiscoverySpecs))
	var err error
	for i, spec := range o.NodeGroupAutoDiscoverySpecs {
		cfgs[i], err = parseASGAutoDiscoverySpec(spec)
		if err != nil {
			return nil, err
		}
	}
	return cfgs, nil
}
func (o NodeGroupDiscoveryOptions) ParseLabelAutoDiscoverySpecs() ([]LabelAutoDiscoveryConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfgs := make([]LabelAutoDiscoveryConfig, len(o.NodeGroupAutoDiscoverySpecs))
	var err error
	for i, spec := range o.NodeGroupAutoDiscoverySpecs {
		cfgs[i], err = parseLabelAutoDiscoverySpec(spec)
		if err != nil {
			return nil, err
		}
	}
	return cfgs, nil
}

type MIGAutoDiscoveryConfig struct {
	Re	*regexp.Regexp
	MinSize	int
	MaxSize	int
}

func parseMIGAutoDiscoverySpec(spec string) (MIGAutoDiscoveryConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfg := MIGAutoDiscoveryConfig{}
	tokens := strings.Split(spec, ":")
	if len(tokens) != 2 {
		return cfg, fmt.Errorf("spec \"%s\" should be discoverer:key=value,key=value", spec)
	}
	discoverer := tokens[0]
	if discoverer != autoDiscovererTypeMIG {
		return cfg, fmt.Errorf("unsupported discoverer specified: %s", discoverer)
	}
	for _, arg := range strings.Split(tokens[1], ",") {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			return cfg, fmt.Errorf("invalid key=value pair %s", kv)
		}
		k, v := kv[0], kv[1]
		var err error
		switch k {
		case migAutoDiscovererKeyPrefix:
			if cfg.Re, err = regexp.Compile(fmt.Sprintf("^%s.+", v)); err != nil {
				return cfg, fmt.Errorf("invalid instance group name prefix \"%s\" - \"^%s.+\" must be a valid RE2 regexp", v, v)
			}
		case migAutoDiscovererKeyMinNodes:
			if cfg.MinSize, err = strconv.Atoi(v); err != nil {
				return cfg, fmt.Errorf("invalid minimum nodes: %s", v)
			}
		case migAutoDiscovererKeyMaxNodes:
			if cfg.MaxSize, err = strconv.Atoi(v); err != nil {
				return cfg, fmt.Errorf("invalid maximum nodes: %s", v)
			}
		default:
			return cfg, fmt.Errorf("unsupported key \"%s\" is specified for discoverer \"%s\". Supported keys are \"%s\"", k, discoverer, validMIGAutoDiscovererKeys)
		}
	}
	if cfg.Re == nil || cfg.Re.String() == "^.+" {
		return cfg, errors.New("empty instance group name prefix supplied")
	}
	if cfg.MinSize > cfg.MaxSize {
		return cfg, fmt.Errorf("minimum size %d is greater than maximum size %d", cfg.MinSize, cfg.MaxSize)
	}
	if cfg.MaxSize < 1 {
		return cfg, fmt.Errorf("maximum size %d must be at least 1", cfg.MaxSize)
	}
	return cfg, nil
}

type ASGAutoDiscoveryConfig struct{ Tags map[string]string }

func parseASGAutoDiscoverySpec(spec string) (ASGAutoDiscoveryConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfg := ASGAutoDiscoveryConfig{}
	tokens := strings.Split(spec, ":")
	if len(tokens) != 2 {
		return cfg, fmt.Errorf("Invalid node group auto discovery spec specified via --node-group-auto-discovery: %s", spec)
	}
	discoverer := tokens[0]
	if discoverer != autoDiscovererTypeASG {
		return cfg, fmt.Errorf("Unsupported discoverer specified: %s", discoverer)
	}
	param := tokens[1]
	kv := strings.SplitN(param, "=", 2)
	if len(kv) != 2 {
		return cfg, fmt.Errorf("invalid key=value pair %s", kv)
	}
	k, v := kv[0], kv[1]
	if k != asgAutoDiscovererKeyTag {
		return cfg, fmt.Errorf("Unsupported parameter key \"%s\" is specified for discoverer \"%s\". The only supported key is \"%s\"", k, discoverer, asgAutoDiscovererKeyTag)
	}
	if v == "" {
		return cfg, errors.New("tag value not supplied")
	}
	p := strings.Split(v, ",")
	if len(p) == 0 {
		return cfg, fmt.Errorf("Invalid ASG tag for auto discovery specified: ASG tag must not be empty")
	}
	cfg.Tags = make(map[string]string, len(p))
	for _, label := range p {
		lp := strings.SplitN(label, "=", 2)
		if len(lp) > 1 {
			cfg.Tags[lp[0]] = lp[1]
			continue
		}
		cfg.Tags[lp[0]] = ""
	}
	return cfg, nil
}

type LabelAutoDiscoveryConfig struct{ Selector map[string]string }

func parseLabelAutoDiscoverySpec(spec string) (LabelAutoDiscoveryConfig, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cfg := LabelAutoDiscoveryConfig{Selector: make(map[string]string)}
	tokens := strings.Split(spec, ":")
	if len(tokens) != 2 {
		return cfg, fmt.Errorf("spec \"%s\" should be discoverer:key=value,key=value", spec)
	}
	discoverer := tokens[0]
	if discoverer != autoDiscovererTypeLabel {
		return cfg, fmt.Errorf("unsupported discoverer specified: %s", discoverer)
	}
	for _, arg := range strings.Split(tokens[1], ",") {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			return cfg, fmt.Errorf("invalid key=value pair %s", kv)
		}
		k, v := kv[0], kv[1]
		cfg.Selector[k] = v
	}
	return cfg, nil
}
