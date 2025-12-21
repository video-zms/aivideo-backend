package ginutil

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"wespy-http-go/helper/timeutil"
	"wespy-http-go/middleware/encrypt"

	"github.com/gin-gonic/gin"
	"wespy-http-go/middleware/compress"
)

const (
	httpHeaderContentType = "Content-Type"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	ctx  *gin.Context
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	if w.ctx != nil {
		w.ctx.Set("write_body_step", timeutil.TrackPoint())
	}
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type AccessParams struct {
	StartTime    int64
	Method       string
	URI          string
	Protocol     string
	UA           string
	Referer      string
	PostData     string
	ClientIP     string
	ContentType  string
	EndTime      int64
	ResponseData []byte
	HttpCode     int
	PanicErr     string
	PanicStack   string
	ExtraHeader  string
}

type AccessLogF func(params *AccessParams)

func AccessLoggerF(logF AccessLogF, logResponse bool, ignoreKeyMap map[string]struct{}, c *gin.Context) {
	params := &AccessParams{
		StartTime:   timeutil.TrackPoint(),
		Method:      c.Request.Method,
		URI:         c.Request.RequestURI,
		Protocol:    c.Request.Proto,
		UA:          c.Request.UserAgent(),
		Referer:     c.Request.Referer(),
		ClientIP:    c.ClientIP(),
		ContentType: c.Request.Header.Get(httpHeaderContentType),
	}

	var w *bodyLogWriter
	if logResponse {
		w = &bodyLogWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer, ctx: c}
		c.Writer = w
	}

	defer func() {
		if err := recover(); err != nil {
			stack := GetStack()
			c.AbortWithStatus(http.StatusInternalServerError)
			params.PanicErr = fmt.Sprintf("%v", err)
			params.PanicStack = stack
		}

		_ = c.Request.ParseForm()
		var buf strings.Builder
		for k, vs := range c.Request.PostForm {
			if _, ok := ignoreKeyMap[k]; ok {
				continue
			}
			for _, v := range vs {
				buf.WriteString(k)
				buf.WriteByte('^')
				buf.WriteString(v)
				buf.WriteByte(' ')
				break
			}
		}
		params.PostData = buf.String()

		params.EndTime = timeutil.TrackPoint()
		params.HttpCode = c.Writer.Status()
		if logResponse && w != nil {
			params.ResponseData = w.body.Bytes()
		}
		params.ExtraHeader = c.Writer.Header().Get(compress.CompressHeaderKey) + "_" + c.Writer.Header().Get(encrypt.EncryptVersionHeader)
		logF(params)
	}()

	c.Next()
}

func GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
