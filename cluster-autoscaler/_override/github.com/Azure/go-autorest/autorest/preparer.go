package autorest

import (
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "io/ioutil"
 "mime/multipart"
 "net/http"
 "net/url"
 "strings"
)

const (
 mimeTypeJSON        = "application/json"
 mimeTypeOctetStream = "application/octet-stream"
 mimeTypeFormPost    = "application/x-www-form-urlencoded"
 headerAuthorization = "Authorization"
 headerContentType   = "Content-Type"
 headerUserAgent     = "User-Agent"
)

type Preparer interface {
 Prepare(*http.Request) (*http.Request, error)
}
type PreparerFunc func(*http.Request) (*http.Request, error)

func (pf PreparerFunc) Prepare(r *http.Request) (*http.Request, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return pf(r)
}

type PrepareDecorator func(Preparer) Preparer

func CreatePreparer(decorators ...PrepareDecorator) Preparer {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return DecoratePreparer(Preparer(PreparerFunc(func(r *http.Request) (*http.Request, error) {
  return r, nil
 })), decorators...)
}
func DecoratePreparer(p Preparer, decorators ...PrepareDecorator) Preparer {
 _logClusterCodePath()
 defer _logClusterCodePath()
 for _, decorate := range decorators {
  p = decorate(p)
 }
 return p
}
func Prepare(r *http.Request, decorators ...PrepareDecorator) (*http.Request, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 if r == nil {
  return nil, NewError("autorest", "Prepare", "Invoked without an http.Request")
 }
 return CreatePreparer(decorators...).Prepare(r)
}
func WithNothing() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   return p.Prepare(r)
  })
 }
}
func WithHeader(header string, value string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    if r.Header == nil {
     r.Header = make(http.Header)
    }
    r.Header.Set(http.CanonicalHeaderKey(header), value)
   }
   return r, err
  })
 }
}
func WithHeaders(headers map[string]interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 h := ensureValueStrings(headers)
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    if r.Header == nil {
     r.Header = make(http.Header)
    }
    for name, value := range h {
     r.Header.Set(http.CanonicalHeaderKey(name), value)
    }
   }
   return r, err
  })
 }
}
func WithBearerAuthorization(token string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithHeader(headerAuthorization, fmt.Sprintf("Bearer %s", token))
}
func AsContentType(contentType string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithHeader(headerContentType, contentType)
}
func WithUserAgent(ua string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithHeader(headerUserAgent, ua)
}
func AsFormURLEncoded() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return AsContentType(mimeTypeFormPost)
}
func AsJSON() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return AsContentType(mimeTypeJSON)
}
func AsOctetStream() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return AsContentType(mimeTypeOctetStream)
}
func WithMethod(method string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r.Method = method
   return p.Prepare(r)
  })
 }
}
func AsDelete() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("DELETE")
}
func AsGet() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("GET")
}
func AsHead() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("HEAD")
}
func AsOptions() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("OPTIONS")
}
func AsPatch() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("PATCH")
}
func AsPost() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("POST")
}
func AsPut() PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithMethod("PUT")
}
func WithBaseURL(baseURL string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    var u *url.URL
    if u, err = url.Parse(baseURL); err != nil {
     return r, err
    }
    if u.Scheme == "" {
     err = fmt.Errorf("autorest: No scheme detected in URL %s", baseURL)
    }
    if err == nil {
     r.URL = u
    }
   }
   return r, err
  })
 }
}
func WithCustomBaseURL(baseURL string, urlParameters map[string]interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 parameters := ensureValueStrings(urlParameters)
 for key, value := range parameters {
  baseURL = strings.Replace(baseURL, "{"+key+"}", value, -1)
 }
 return WithBaseURL(baseURL)
}
func WithFormData(v url.Values) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    s := v.Encode()
    if r.Header == nil {
     r.Header = make(http.Header)
    }
    r.Header.Set(http.CanonicalHeaderKey(headerContentType), mimeTypeFormPost)
    r.ContentLength = int64(len(s))
    r.Body = ioutil.NopCloser(strings.NewReader(s))
   }
   return r, err
  })
 }
}
func WithMultiPartFormData(formDataParameters map[string]interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    var body bytes.Buffer
    writer := multipart.NewWriter(&body)
    for key, value := range formDataParameters {
     if rc, ok := value.(io.ReadCloser); ok {
      var fd io.Writer
      if fd, err = writer.CreateFormFile(key, key); err != nil {
       return r, err
      }
      if _, err = io.Copy(fd, rc); err != nil {
       return r, err
      }
     } else {
      if err = writer.WriteField(key, ensureValueString(value)); err != nil {
       return r, err
      }
     }
    }
    if err = writer.Close(); err != nil {
     return r, err
    }
    if r.Header == nil {
     r.Header = make(http.Header)
    }
    r.Header.Set(http.CanonicalHeaderKey(headerContentType), writer.FormDataContentType())
    r.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
    r.ContentLength = int64(body.Len())
    return r, err
   }
   return r, err
  })
 }
}
func WithFile(f io.ReadCloser) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    b, err := ioutil.ReadAll(f)
    if err != nil {
     return r, err
    }
    r.Body = ioutil.NopCloser(bytes.NewReader(b))
    r.ContentLength = int64(len(b))
   }
   return r, err
  })
 }
}
func WithBool(v bool) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithString(fmt.Sprintf("%v", v))
}
func WithFloat32(v float32) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithString(fmt.Sprintf("%v", v))
}
func WithFloat64(v float64) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithString(fmt.Sprintf("%v", v))
}
func WithInt32(v int32) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithString(fmt.Sprintf("%v", v))
}
func WithInt64(v int64) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return WithString(fmt.Sprintf("%v", v))
}
func WithString(v string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    r.ContentLength = int64(len(v))
    r.Body = ioutil.NopCloser(strings.NewReader(v))
   }
   return r, err
  })
 }
}
func WithJSON(v interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    b, err := json.Marshal(v)
    if err == nil {
     r.ContentLength = int64(len(b))
     r.Body = ioutil.NopCloser(bytes.NewReader(b))
    }
   }
   return r, err
  })
 }
}
func WithPath(path string) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    if r.URL == nil {
     return r, NewError("autorest", "WithPath", "Invoked with a nil URL")
    }
    if r.URL, err = parseURL(r.URL, path); err != nil {
     return r, err
    }
   }
   return r, err
  })
 }
}
func WithEscapedPathParameters(path string, pathParameters map[string]interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 parameters := escapeValueStrings(ensureValueStrings(pathParameters))
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    if r.URL == nil {
     return r, NewError("autorest", "WithEscapedPathParameters", "Invoked with a nil URL")
    }
    for key, value := range parameters {
     path = strings.Replace(path, "{"+key+"}", value, -1)
    }
    if r.URL, err = parseURL(r.URL, path); err != nil {
     return r, err
    }
   }
   return r, err
  })
 }
}
func WithPathParameters(path string, pathParameters map[string]interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 parameters := ensureValueStrings(pathParameters)
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    if r.URL == nil {
     return r, NewError("autorest", "WithPathParameters", "Invoked with a nil URL")
    }
    for key, value := range parameters {
     path = strings.Replace(path, "{"+key+"}", value, -1)
    }
    if r.URL, err = parseURL(r.URL, path); err != nil {
     return r, err
    }
   }
   return r, err
  })
 }
}
func parseURL(u *url.URL, path string) (*url.URL, error) {
 _logClusterCodePath()
 defer _logClusterCodePath()
 p := strings.TrimRight(u.String(), "/")
 if !strings.HasPrefix(path, "/") {
  path = "/" + path
 }
 return url.Parse(p + path)
}
func WithQueryParameters(queryParameters map[string]interface{}) PrepareDecorator {
 _logClusterCodePath()
 defer _logClusterCodePath()
 parameters := ensureValueStrings(queryParameters)
 return func(p Preparer) Preparer {
  return PreparerFunc(func(r *http.Request) (*http.Request, error) {
   r, err := p.Prepare(r)
   if err == nil {
    if r.URL == nil {
     return r, NewError("autorest", "WithQueryParameters", "Invoked with a nil URL")
    }
    v := r.URL.Query()
    for key, value := range parameters {
     d, err := url.QueryUnescape(value)
     if err != nil {
      return r, err
     }
     v.Add(key, d)
    }
    r.URL.RawQuery = v.Encode()
   }
   return r, err
  })
 }
}
