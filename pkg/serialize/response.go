package serialize

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type RspItem struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func Response(code int, message string, data interface{}) *RspItem {
	return &RspItem{
		Code:   code,
		Msg:    message,
		Result: data,
	}
}

func GetResponseCode(c *gin.Context, responseData []byte) int {
	code := c.GetInt("response_code")
	if code != 0 {
		return code
	}
	code = DecodeResponseCode(responseData)
	if code != 0 {
		return code
	}
	return 0
}

func SetResponseCode(c *gin.Context, code int) {
	c.Set("response_code", code)
}

func DecodeResponseCode(respBs []byte) int {
	var result struct {
		Code int `json:"code"`
	}
	_ = json.Unmarshal(respBs, &result)
	return result.Code
}

type ErrMsg struct {
	code int
	msg  string
	err  error
	data interface{}
}

func (e *ErrMsg) IsErr() bool {
	return e.err != nil
}

func (e *ErrMsg) GetErr() error {
	return e.err
}
