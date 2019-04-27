package history

import (
	"time"
)

type Sample struct {
	Value		float64
	Timestamp	time.Time
}
type Timeseries struct {
	Labels	map[string]string
	Samples	[]Sample
}
