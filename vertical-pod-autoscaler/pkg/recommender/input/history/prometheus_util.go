package history

import (
 "encoding/json"
 "fmt"
 "io"
 "time"
)

type responseType struct {
 Status      string   `json:"status"`
 Data        dataType `json:"data"`
 ErrorType   string   `json:"errorType"`
 ErrorString string   `json:"error"`
}
type dataType struct {
 ResultType string          `json:"resultType"`
 Result     json.RawMessage `json:"result"`
}
type matrixType struct {
 Metric map[string]string `json:"metric"`
 Values [][]interface{}   `json:"values"`
}

func decodeSamples(input [][]interface{}) ([]Sample, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 res := make([]Sample, 0)
 for _, item := range input {
  if len(item) != 2 {
   return nil, fmt.Errorf("invalid length: %d", len(item))
  }
  ts, ok := item[0].(float64)
  if !ok {
   return nil, fmt.Errorf("invalid time: %v", item[0])
  }
  stringVal, ok := item[1].(string)
  if !ok {
   return nil, fmt.Errorf("invalid value: %v", item[1])
  }
  var val float64
  fmt.Sscan(stringVal, &val)
  res = append(res, Sample{Value: val, Timestamp: time.Unix(int64(ts), 0)})
 }
 return res, nil
}
func decodeTimeseriesFromResponse(input io.Reader) ([]Timeseries, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var resp responseType
 err := json.NewDecoder(input).Decode(&resp)
 if err != nil {
  return nil, fmt.Errorf("couldn't parse response: %v", err)
 }
 if resp.Status != "success" || resp.Data.ResultType != "matrix" {
  return nil, fmt.Errorf("invalid response status: %s or type: %s", resp.Status, resp.Data.ResultType)
 }
 var matrices []matrixType
 err = json.Unmarshal(resp.Data.Result, &matrices)
 if err != nil {
  return nil, fmt.Errorf("couldn't parse response matrix: %v", err)
 }
 res := make([]Timeseries, 0)
 for _, matrix := range matrices {
  samples, err := decodeSamples(matrix.Values)
  if err != nil {
   return []Timeseries{}, fmt.Errorf("error decoding samples: %v", err)
  }
  res = append(res, Timeseries{Labels: matrix.Metric, Samples: samples})
 }
 return res, nil
}
