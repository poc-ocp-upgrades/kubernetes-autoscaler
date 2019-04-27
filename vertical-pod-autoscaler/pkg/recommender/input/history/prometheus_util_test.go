package history

import (
	"strings"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestSingleTimeseries(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := `{"status":"success",
	       "data":{
		 "resultType":"matrix",
		 "result":[{
		   "metric":{
	             "__name__":"up",
		     "x":"y"},
		   "values":[[1515422500.45,"0."],
			    [1515422560.453,"1."],
			    [1515422620.45,"0"]]}]}}"`
	res, err := decodeTimeseriesFromResponse(strings.NewReader(s))
	assert.Nil(t, err)
	assert.Equal(t, res, []Timeseries{{Labels: map[string]string{"__name__": "up", "x": "y"}, Samples: []Sample{{Value: 0, Timestamp: time.Unix(1515422500, 0)}, {Value: 1, Timestamp: time.Unix(1515422560, 0)}, {Value: 0, Timestamp: time.Unix(1515422620, 0)}}}})
}
func TestEmptyResponse(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := `{"status":"success", "data":{"resultType":"matrix", "result":[]}}`
	res, err := decodeTimeseriesFromResponse(strings.NewReader(s))
	assert.Nil(t, err)
	assert.Equal(t, res, []Timeseries{})
}
func TestResponseError(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := `{"status":"error", "error":"my bad", "errorType":"some"}`
	res, err := decodeTimeseriesFromResponse(strings.NewReader(s))
	assert.Nil(t, res)
	assert.NotNil(t, err)
}
func TestParseError(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := `{"status":"success", "other-key":[unparsable], "data":{"resultType":"matrix", "result":[]}}`
	res, err := decodeTimeseriesFromResponse(strings.NewReader(s))
	assert.Nil(t, res)
	assert.NotNil(t, err)
}
func TestTwoTimeseries(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s := `{"status":"success",
	       "data":{
		 "resultType":"matrix",
		 "result":[
		   {"metric":{"x":"y"},
		    "values":[[1515422620,"15"]]},
		   {"metric":{"x":"z"},
		    "values":[]}]}}"`
	res, err := decodeTimeseriesFromResponse(strings.NewReader(s))
	assert.Nil(t, err)
	assert.Equal(t, res, []Timeseries{{Labels: map[string]string{"x": "y"}, Samples: []Sample{{Value: 15, Timestamp: time.Unix(1515422620, 0)}}}, {Labels: map[string]string{"x": "z"}, Samples: []Sample{}}})
}
