package date

import (
	"errors"
	"time"
)

const (
	rfc1123JSON	= `"` + time.RFC1123 + `"`
	rfc1123		= time.RFC1123
)

type TimeRFC1123 struct{ time.Time }

func (t *TimeRFC1123) UnmarshalJSON(data []byte) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Time, err = ParseTime(rfc1123JSON, string(data))
	if err != nil {
		return err
	}
	return nil
}
func (t TimeRFC1123) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	b := []byte(t.Format(rfc1123JSON))
	return b, nil
}
func (t TimeRFC1123) MarshalText() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalText: year outside of range [0,9999]")
	}
	b := []byte(t.Format(rfc1123))
	return b, nil
}
func (t *TimeRFC1123) UnmarshalText(data []byte) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Time, err = ParseTime(rfc1123, string(data))
	if err != nil {
		return err
	}
	return nil
}
func (t TimeRFC1123) MarshalBinary() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.MarshalText()
}
func (t *TimeRFC1123) UnmarshalBinary(data []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.UnmarshalText(data)
}
func (t TimeRFC1123) ToTime() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return t.Time
}
func (t TimeRFC1123) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b, err := t.MarshalText()
	if err != nil {
		return ""
	}
	return string(b)
}
