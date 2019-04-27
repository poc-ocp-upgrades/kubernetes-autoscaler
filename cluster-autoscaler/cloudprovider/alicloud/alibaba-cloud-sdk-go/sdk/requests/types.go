package requests

import "strconv"

type Integer string

func NewInteger(integer int) Integer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Integer(strconv.Itoa(integer))
}
func (integer Integer) HasValue() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return integer != ""
}
func (integer Integer) GetValue() (int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strconv.Atoi(string(integer))
}
func NewInteger64(integer int64) Integer {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Integer(strconv.FormatInt(integer, 10))
}
func (integer Integer) GetValue64() (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strconv.ParseInt(string(integer), 10, 0)
}

type Boolean string

func NewBoolean(bool bool) Boolean {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Boolean(strconv.FormatBool(bool))
}
func (boolean Boolean) HasValue() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return boolean != ""
}
func (boolean Boolean) GetValue() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strconv.ParseBool(string(boolean))
}

type Float string

func NewFloat(f float64) Float {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Float(strconv.FormatFloat(f, 'f', 6, 64))
}
func (float Float) HasValue() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return float != ""
}
func (float Float) GetValue() (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strconv.ParseFloat(string(float), 64)
}
