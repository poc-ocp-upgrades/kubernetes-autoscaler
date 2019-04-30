package date

import (
	"strings"
	"time"
)

func ParseTime(format string, t string) (d time.Time, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return time.Parse(format, strings.ToUpper(t))
}
