package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rifflock/lfshook"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	LogContextKey       = "log_entry"
	ContextKeyRequestID = "requestId"
	ContextKeyNoLog     = "_no_log"
	SLoggerKey          = "slogger"
)

// header key
const (
	RequestIDHeaderKey = "X_Launchpad_RequestId"
)

const (
	printRequestLen  = 10240
	printResponseLen = 10240
)

type Conf struct {
	LogLevel string `json:"log_level"`
	LogPath  string `json:"log_path"`
}

type Formatter struct {
	logrus.Formatter
}

type Logger struct {
	*logrus.Logger
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type Response struct {
	Data any    `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

var loc, _ = time.LoadLocation("Asia/Shanghai")

var (
	Entry   *logrus.Entry
	SLogger *Logger
)

func newFormatter(formatter logrus.Formatter) *Formatter {
	return &Formatter{formatter}
}

func (u Formatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.In(loc)
	return u.Formatter.Format(e)
}

func InitLog(conf Conf) (*Logger, error) {
	var l logrus.Level
	if conf.LogLevel == "" {
		l = logrus.InfoLevel
	} else {
		var err error
		l, err = logrus.ParseLevel(conf.LogLevel)
		if err != nil {
			return &Logger{logrus.StandardLogger()}, err
		}
	}
	logrus.SetLevel(l)
	if conf.LogPath != "" {
		hook := &lumberjack.Logger{
			Filename:   conf.LogPath + "." + time.Now().Format("2006-01-02"),
			MaxAge:     7, // 最大保存天数（day）
			MaxBackups: 3,
			Compress:   true, // 是否压缩
			LocalTime:  false,
		}

		logrus.SetOutput(hook)
		logrus.AddHook(newRotateHook(conf.LogPath, 7*24*time.Hour, 24*time.Hour))
	}
	logrus.SetFormatter(newFormatter(&logrus.JSONFormatter{}))
	SLogger = &Logger{logrus.StandardLogger()}
	return &Logger{logrus.StandardLogger()}, nil
}

func newRotateHook(logPath string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	writer, err := rotatelogs.New(
		logPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		logrus.Errorf("config local file system logger serror. %+v", errors.WithStack(err))
	}
	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
}

func (log *Logger) LogEntryWithContext(ctx *gin.Context, fieldsList ...logrus.Fields) *logrus.Entry {
	var fields logrus.Fields
	if len(fieldsList) > 0 {
		fields = fieldsList[0]
	}

	fields["requestId"] = GetRequestID(ctx)
	var ok bool
	Entry, ok = ctx.Value(LogContextKey).(*logrus.Entry)
	if !ok {
		Entry = logrus.NewEntry(log.Logger).WithFields(fields)
		ctx.Set(LogContextKey, Entry)
		return Entry.WithFields(fields)
	}
	return Entry.WithFields(fields)
}

func AddFields(c *gin.Context, fieldsList ...logrus.Fields) *logrus.Entry {
	var fields logrus.Fields
	if len(fieldsList) > 0 {
		fields = fieldsList[0]
	}

	return SLogger.LogEntryWithContext(c, fields)
}

func GetRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return genRequestId()
	}

	// 从ctx中获取
	if r := ctx.GetString(ContextKeyRequestID); r != "" {
		return r
	}

	// 优先从header中获取
	var requestId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		requestId = ctx.Request.Header.Get(RequestIDHeaderKey)
	}

	// 新生成
	if requestId == "" {
		requestId = genRequestId()
	}

	ctx.Set(ContextKeyRequestID, requestId)
	ctx.Request.Header.Set(RequestIDHeaderKey, requestId)
	ctx.Writer.Header().Set(RequestIDHeaderKey, requestId)
	return requestId
}

func getLogger() *Logger {
	if SLogger == nil {
		SLogger = &Logger{logrus.StandardLogger()}
	}
	return SLogger
}

func GetEntry() *logrus.Entry {
	if Entry == nil {
		Entry = logrus.NewEntry(getLogger().Logger)
	}
	return Entry
}

// 通用字段封装
func commonLogger(ctx *gin.Context) *logrus.Entry {
	if ctx == nil {
		return GetEntry()
	}
	return SLogger.LogEntryWithContext(ctx, logrus.Fields{
		"module":  GetAppName(ctx),
		"localIp": GetLocalIp(),
		"uri":     ctx.Request.URL.Path,
	})
}

// Debug 提供给业务使用的server log 日志打印方法
func Debug(ctx *gin.Context, args ...interface{}) {
	commonLogger(ctx).Debug(args...)
}

func Debugf(ctx *gin.Context, format string, args ...interface{}) {
	commonLogger(ctx).Debugf(format, args...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	commonLogger(ctx).Info(args...)
}

func Infof(ctx *gin.Context, format string, args ...interface{}) {
	commonLogger(ctx).Infof(format, args...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	commonLogger(ctx).Warn(args...)
}

func Warnf(ctx *gin.Context, format string, args ...interface{}) {
	commonLogger(ctx).Warnf(format, args...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	commonLogger(ctx).Error(args...)
}

func Errorf(ctx *gin.Context, format string, args ...interface{}) {
	commonLogger(ctx).Errorf(format, args...)
}

func SetNoLogFlag(ctx *gin.Context) {
	ctx.Set(ContextKeyNoLog, true)
}

func NoLog(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	flag, ok := ctx.Get(ContextKeyNoLog)
	if ok && flag == true {
		return true
	}
	return false
}

// AddNotice 用户自定义Notice
func AddNotice(ctx *gin.Context, key string, val interface{}) {
	if meta, ok := CtxFromGinContext(ctx); ok {
		if n := Value(meta, Notice); n != nil {
			if _, ok = n.(map[string]interface{}); ok {
				notices := n.(map[string]interface{})
				notices[key] = val
			}
		}
	}
}

// GetCustomerKeyValue 获得所有用户自定义的Notice
func GetCustomerKeyValue(ctx *gin.Context) map[string]interface{} {
	meta, ok := CtxFromGinContext(ctx)
	if !ok {
		return nil
	}

	n := Value(meta, Notice)
	if n == nil {
		return nil
	}
	if notices, ok := n.(map[string]interface{}); ok {
		return notices
	}

	return nil
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	s = strings.Replace(s, "\n", "", -1)
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if SLogger != nil {
			c.Set(SLoggerKey, SLogger)
		}
		start := time.Now()

		// 请求报文
		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = c.GetRawData()
			if err != nil {
				SLogger.Warnf("get http request body serror: %s", err.Error())
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		param := make(map[string]interface{})
		_ = json.Unmarshal(requestBody, &param)

		requestParam := ""
		if c.Request.URL.RawQuery != "" {
			requestParam += "&" + c.Request.URL.RawQuery
		}

		if len(requestParam) > printRequestLen {
			requestParam = requestParam[:printRequestLen]
		}

		// 请求url
		originUri := c.Request.URL.Path
		path := c.FullPath()

		fields := map[string]interface{}{
			"logType":          "request",
			"originUri":        originUri,
			"uri":              path,
			"host":             c.Request.Host,
			"httpProto":        c.Request.Proto,
			"method":           c.Request.Method,
			"clientIp":         GetLocalIp(),
			"refer":            c.Request.Referer(),
			"requestId":        GetRequestID(c),
			"requestStartTime": start,
			"requestParam":     requestParam,
			"requestBody":      param,
			"cookie":           getCookie(c),
			"module":           GetAppName(c),
			"timestamp":        start.Unix(),
			"uniqUri":          c.Request.Method + "_" + path,
		}

		logg := SLogger.WithFields(fields)
		logg.Info()
		UseMetadata(c)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		c.Writer = blw

		defer func() {
			msg := ""
			if err := recover(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				stack := getStack(3)

				msg = fmt.Sprintf("panic recovered:[%s] - stack info: [%s]", err, stack)
				end := time.Now()
				responseFields := map[string]interface{}{
					"logType":        "response",
					"status":         c.Writer.Status(),
					"requestEndTime": end,
					"timestamp":      end.Unix(),
				}

				for k, v := range responseFields {
					fields[k] = v
				}
				logg = SLogger.WithFields(fields)
				logg.Error(msg)
				return
			}

			switch {
			case c.Writer.Status() >= http.StatusInternalServerError:
				msg += fmt.Sprintf("[%s]", c.Errors.String())
				logg.Error(msg)
			case len(c.Errors) > 0:
				msg += fmt.Sprintf("[%s]", c.Errors.String())
				logg.Info(msg)
			default:
				logg.Info(msg)
			}
		}()
		c.Next()

		responseStr := ""
		if blw.body != nil {
			if len(blw.body.String()) <= printResponseLen {
				responseStr = blw.body.String()
			} else {
				responseStr = blw.body.String()[:printResponseLen]
			}
		}
		response := Response{}
		_ = json.Unmarshal([]byte(responseStr), &response)

		cost := uint(time.Since(start) / time.Millisecond)
		// 用户自定义notice
		customerFields := GetCustomerKeyValue(c)

		end := time.Now()

		responseFields := map[string]interface{}{
			"logType":        "response",
			"responseCode":   response.Code,
			"cost":           cost,
			"status":         c.Writer.Status(),
			"refer":          c.Request.Referer(),
			"requestEndTime": end,
			"response":       response,
		}

		for k, v := range responseFields {
			fields[k] = v
		}

		for k, v := range customerFields {
			fields[k] = v
		}

		logg = SLogger.WithFields(fields)

	}
}

func getStack(skip int) string {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.String()
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.ReplaceAll(name, centerDot, dot)
	return name
}

func genRequestId() (requestId string) {
	// 随机生成
	usec := uint64(time.Now().UnixNano())
	requestId = strconv.FormatUint(usec&0x7FFFFFFF|0x80000000, 10)
	return requestId
}

// // access 添加kv打印
// func AddNotice(k string, v interface{}) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		log.AddNotice(c, k, v)
// 		c.Next()
// 	}
// }

func getCookie(ctx *gin.Context) string {
	cStr := ""
	for _, c := range ctx.Request.Cookies() {
		cStr += fmt.Sprintf("%s=%s&", c.Name, c.Value)
	}
	return strings.TrimRight(cStr, "&")
}

func SetModule(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("module", module)
	}
}

func GetAppName(ctx *gin.Context) string {
	name, _ := ctx.Get("module")
	if name == nil {
		return "launchpad"
	}
	return name.(string)
}

// GetLocalIp 获取本机ip
func GetLocalIp() string {
	addrs, _ := net.InterfaceAddrs()
	var ip string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				if ip != "127.0.0.1" {
					return ip
				}
			}
		}
	}
	return "127.0.0.1"
}
