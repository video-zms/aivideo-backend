package ginutil

import (
	"bytes"

	"github.com/bitly/go-simplejson"

	"github.com/gin-gonic/gin"
)

type kvGetF func(c *gin.Context, oldResp string) (k string, v interface{})

type noWriteWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w noWriteWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
	//return w.ResponseWriter.Write(b)
}
func (w noWriteWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
	//return w.ResponseWriter.WriteString(s)
}

func InjectRespF(kvGetter kvGetF, c *gin.Context) {
	injectField(kvGetter, c)
}

func injectField(kvGetter kvGetF, c *gin.Context) {
	//替换请求的writer
	origW := c.Writer
	w := &noWriteWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
	c.Writer = w

	//等待后续执行
	c.Next()

	//取出回包内容，准备注入
	k, v := kvGetter(c, w.body.String())
	if k == "" {
		//直接返回
		c.Writer = origW
		c.Writer.Write(w.body.Bytes())
		return
	}

	//在json result层注入
	j, err := simplejson.NewJson(w.body.Bytes())
	if err != nil {
		//直接返回
		c.Writer = origW
		c.Writer.Write(w.body.Bytes())
		return
	}
	j.Get("result").Set(k, v)
	bs, _ := j.MarshalJSON()

	c.Writer = origW
	c.Writer.Write(bs)
}
