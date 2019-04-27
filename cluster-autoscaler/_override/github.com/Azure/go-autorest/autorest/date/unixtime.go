package date

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"
)

var unixEpoch = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

type UnixTime time.Time

func (t UnixTime) Duration() time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return time.Time(t).Sub(unixEpoch)
}
func NewUnixTimeFromSeconds(seconds float64) UnixTime {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewUnixTimeFromDuration(time.Duration(seconds * float64(time.Second)))
}
func NewUnixTimeFromNanoseconds(nanoseconds int64) UnixTime {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return NewUnixTimeFromDuration(time.Duration(nanoseconds))
}
func NewUnixTimeFromDuration(dur time.Duration) UnixTime {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return UnixTime(unixEpoch.Add(dur))
}
func UnixEpoch() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return unixEpoch
}
func (t UnixTime) MarshalJSON() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	err := enc.Encode(float64(time.Time(t).UnixNano()) / 1e9)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func (t *UnixTime) UnmarshalJSON(text []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dec := json.NewDecoder(bytes.NewReader(text))
	var secondsSinceEpoch float64
	if err := dec.Decode(&secondsSinceEpoch); err != nil {
		return err
	}
	*t = NewUnixTimeFromSeconds(secondsSinceEpoch)
	return nil
}
func (t UnixTime) MarshalText() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cast := time.Time(t)
	return cast.MarshalText()
}
func (t *UnixTime) UnmarshalText(raw []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var unmarshaled time.Time
	if err := unmarshaled.UnmarshalText(raw); err != nil {
		return err
	}
	*t = UnixTime(unmarshaled)
	return nil
}
func (t UnixTime) MarshalBinary() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	buf := &bytes.Buffer{}
	payload := int64(t.Duration())
	if err := binary.Write(buf, binary.LittleEndian, &payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (t *UnixTime) UnmarshalBinary(raw []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var nanosecondsSinceEpoch int64
	if err := binary.Read(bytes.NewReader(raw), binary.LittleEndian, &nanosecondsSinceEpoch); err != nil {
		return err
	}
	*t = NewUnixTimeFromNanoseconds(nanosecondsSinceEpoch)
	return nil
}
