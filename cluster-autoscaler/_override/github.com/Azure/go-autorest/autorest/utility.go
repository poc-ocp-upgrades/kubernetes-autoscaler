package autorest

import (
 "bytes"
 "encoding/json"
 "encoding/xml"
 "fmt"
 "io"
 "net"
 "net/http"
 "net/url"
 "reflect"
 "strings"
 "github.com/Azure/go-autorest/autorest/adal"
)

type EncodedAs string

const (
 EncodedAsJSON EncodedAs = "JSON"
 EncodedAsXML  EncodedAs = "XML"
)

type Decoder interface{ Decode(v interface{}) error }

func NewDecoder(encodedAs EncodedAs, r io.Reader) Decoder {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if encodedAs == EncodedAsJSON {
  return json.NewDecoder(r)
 } else if encodedAs == EncodedAsXML {
  return xml.NewDecoder(r)
 }
 return nil
}
func CopyAndDecode(encodedAs EncodedAs, r io.Reader, v interface{}) (bytes.Buffer, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 b := bytes.Buffer{}
 return b, NewDecoder(encodedAs, io.TeeReader(r, &b)).Decode(v)
}
func TeeReadCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return &teeReadCloser{rc, io.TeeReader(rc, w)}
}

type teeReadCloser struct {
 rc io.ReadCloser
 r  io.Reader
}

func (t *teeReadCloser) Read(p []byte) (int, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return t.r.Read(p)
}
func (t *teeReadCloser) Close() error {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return t.rc.Close()
}
func containsInt(ints []int, n int) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, i := range ints {
  if i == n {
   return true
  }
 }
 return false
}
func escapeValueStrings(m map[string]string) map[string]string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for key, value := range m {
  m[key] = url.QueryEscape(value)
 }
 return m
}
func ensureValueStrings(mapOfInterface map[string]interface{}) map[string]string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 mapOfStrings := make(map[string]string)
 for key, value := range mapOfInterface {
  mapOfStrings[key] = ensureValueString(value)
 }
 return mapOfStrings
}
func ensureValueString(value interface{}) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if value == nil {
  return ""
 }
 switch v := value.(type) {
 case string:
  return v
 case []byte:
  return string(v)
 default:
  return fmt.Sprintf("%v", v)
 }
}
func MapToValues(m map[string]interface{}) url.Values {
 _logClusterCodePath()
 defer _logClusterCodePath()
 v := url.Values{}
 for key, value := range m {
  x := reflect.ValueOf(value)
  if x.Kind() == reflect.Array || x.Kind() == reflect.Slice {
   for i := 0; i < x.Len(); i++ {
    v.Add(key, ensureValueString(x.Index(i)))
   }
  } else {
   v.Add(key, ensureValueString(value))
  }
 }
 return v
}
func AsStringSlice(s interface{}) ([]string, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 v := reflect.ValueOf(s)
 if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
  return nil, NewError("autorest", "AsStringSlice", "the value's type is not an array.")
 }
 stringSlice := make([]string, 0, v.Len())
 for i := 0; i < v.Len(); i++ {
  stringSlice = append(stringSlice, v.Index(i).String())
 }
 return stringSlice, nil
}
func String(v interface{}, sep ...string) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if len(sep) == 0 {
  return ensureValueString(v)
 }
 stringSlice, ok := v.([]string)
 if ok == false {
  var err error
  stringSlice, err = AsStringSlice(v)
  if err != nil {
   panic(fmt.Sprintf("autorest: Couldn't convert value to a string %s.", err))
  }
 }
 return ensureValueString(strings.Join(stringSlice, sep[0]))
}
func Encode(location string, v interface{}, sep ...string) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 s := String(v, sep...)
 switch strings.ToLower(location) {
 case "path":
  return pathEscape(s)
 case "query":
  return queryEscape(s)
 default:
  return s
 }
}
func pathEscape(s string) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}
func queryEscape(s string) string {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return url.QueryEscape(s)
}
func ChangeToGet(req *http.Request) *http.Request {
 _logClusterCodePath()
 defer _logClusterCodePath()
 req.Method = "GET"
 req.Body = nil
 req.ContentLength = 0
 req.Header.Del("Content-Length")
 return req
}
func IsTokenRefreshError(err error) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if _, ok := err.(adal.TokenRefreshError); ok {
  return true
 }
 if de, ok := err.(DetailedError); ok {
  return IsTokenRefreshError(de.Original)
 }
 return false
}
func IsTemporaryNetworkError(err error) bool {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if netErr, ok := err.(net.Error); !ok || (ok && netErr.Temporary()) {
  return true
 }
 return false
}
