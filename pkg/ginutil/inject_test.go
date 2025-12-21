package ginutil

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestInject(t *testing.T) {
	go func() {
		for {
			time.Sleep(time.Second)
			go func() {
				t.Log("posting")
				resp, err := http.Post("http://127.0.0.1:7072/test", "", nil)
				if err == nil {
					respBody, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						t.Log(string(respBody))
					}
				}
			}()
		}
	}()

	e := gin.New()
	e.Use(func(context *gin.Context) {
		InjectRespF(func(c *gin.Context, oldResp string) (k string, v interface{}) {
			return "injected_time", time.Now().Unix()
		}, context)
	})
	e.POST("test", testGin)
	e.Run(":7072")
}

func testGin(c *gin.Context) {
	c.JSON(200, map[string]interface{}{"a": "a", "b": "b"})
}
