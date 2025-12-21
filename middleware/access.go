package middleware

import (
	"axe-backend/pkg/ginutil"
	"axe-backend/util"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"axe-backend/pkg/serialize"

	"github.com/gin-gonic/gin"
)

var accessLogger *logrus.Logger

type commonParams struct {
	Uid int `form:"uid"`
}

var ignoreKeyMap = map[string]struct{}{
	"channel":          {},
	"device_id":        {},
	"device_name":      {},
	"distinct_id":      {},
	"market":           {},
	"nonce":            {},
	"push_id":          {},
	"shumei_device_id": {},
	"sid":              {},
	//"version_code":       {},
	//"version_name":       {},
	"full_pkg_name":      {},
	"pkg_name":           {},
	"millisecond":        {},
	"wespy_access_token": {},
	"wespy_sign":         {},
	"android_version":    {},
	"is_google_weplay":   {},
	"audio_data":         {},
}

var ignorePanicMap = map[string]struct{}{
	"connection reset by peer":         {},
	"i/o timeout":                      {},
	"broken pipe":                      {},
	"use of closed network connection": {},
}

func AccessLogWithRecover(c *gin.Context) {
	ginutil.AccessLoggerF(func(params *ginutil.AccessParams) {
		accessLogMap := make(map[string]interface{})
		//req meta
		accessLogMap["method"] = params.Method
		accessLogMap["protocol"] = params.Protocol
		accessLogMap["ua"] = params.UA
		accessLogMap["referer"] = params.Referer
		accessLogMap["client_ip"] = params.ClientIP
		accessLogMap["content_type"] = params.ContentType

		//req data
		accessLogMap["uri"] = params.URI
		accessLogMap["post_data"] = params.PostData
		//common params
		commonParam := new(commonParams)
		_ = c.ShouldBind(commonParam)
		accessLogMap["uid"] = commonParam.Uid
		accessLogMap["product_version"] = VersionFromReq(c)
		accessLogMap["product_name"] = c.GetHeader("WepieProduct")
		accessLogMap["header"] = c.Request.Header
		accessLogMap["RemoteAddr"] = c.Request.RemoteAddr

		//resp
		readStep := c.GetInt64("wait_body_step")
		if readStep > 0 {
			readStep = (readStep - params.StartTime) / 1e6
			accessLogMap["wait_body"] = readStep
		}
		writeStep := c.GetInt64("write_body_step")
		if writeStep > 0 {
			writeStep = (params.EndTime - writeStep) / 1e6
			accessLogMap["write_body"] = writeStep
		}
		accessLogMap["response_ms"] = (params.EndTime-params.StartTime)/1e6 - readStep - writeStep
		accessLogMap["bytes_sent"] = len(params.ResponseData)
		accessLogMap["http_code"] = params.HttpCode
		accessLogMap["panic_err"] = params.PanicErr
		accessLogMap["response_code"] = serialize.GetResponseCode(c, params.ResponseData)
		accessLogMap["x_header"] = params.ExtraHeader // 压缩算法+加密算法
		if params.PanicErr != "" {
			entry := logrus.WithField("path", params.URI).WithField("err", params.PanicErr).WithField("stack", params.PanicStack).
				WithField("costTimeMs", (params.EndTime-params.StartTime)/1e6)
			if isValidPanic(params.PanicErr) {
				localIp := util.GetLocalIp()
				_ = fmt.Sprintf("服务器出现异常panic, 服务器ip: %s, path: %s，err: %v, stack: %s", localIp, params.URI, params.PanicErr, params.PanicStack)
				// sms.SendITAlarmMsg(msg)
				// sms.SendToSlsGroup(msg)
				entry.Errorln("PanicRecover")
			} else {
				entry.Infoln("PanicRecover")
			}
		}

		accessLogger.WithFields(accessLogMap).Infoln("access_log")
	}, true, ignoreKeyMap, c)
}

func containsRiskInfo(uri string, uid int, respValue string) (bool, string) {
	// return false, ""
	// riskMonitorConfig := weconfig.GetRiskConfig2023().RiskMonitorConfig
	// for _, path := range riskMonitorConfig.WhiteUriList {
	// 	if strings.Contains(uri, path) {
	// 		return false, ""
	// 	}
	// }
	// if !riskMonitorConfig.InSaltGray(uid) {
	// 	return false, ""
	// }
	// if find, findStr := helper.ContainsPhoneNumber(respValue); find {
	// 	return true, findStr
	// }
	// if find, findStr := helper.ContainsIdCard(respValue); find {
	// 	return true, findStr
	// }
	return false, ""
}

func isValidPanic(err string) bool {
	for ignoreErr := range ignorePanicMap {
		if strings.Contains(err, ignoreErr) {
			return false
		}
	}
	return true
}

func VersionFromReq(c *gin.Context) string {
	return c.GetHeader("WepieVersion")
}

var ignorePanic = []error{
	http.ErrAbortHandler, // https://github.com/golang/go/issues/28239
}

func AdminRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		defer func() {
			if err := recover(); err != nil {
				_ = c.Request.ParseForm()
				buf := strings.Builder{}
				buf.Grow(128)
				for k, vs := range c.Request.PostForm {
					for _, v := range vs {
						buf.WriteString(k)
						buf.WriteByte('^')
						buf.WriteString(v)
						buf.WriteByte(' ')
						break
					}
				}

				logMap := logrus.Fields{}
				logMap["err"] = err
				logMap["post_data"] = buf.String()
				logMap["url"] = c.Request.RequestURI
				logMap["client_ip"] = c.ClientIP()
				logMap["cost"] = time.Since(now).Milliseconds()
				// 如果是可忽略的panic,这里截住
				for _, targetErr := range ignorePanic {
					if targetErr.Error() == err {
						logrus.WithFields(logMap).Warnln("AdminRecovery ignore panic")
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}
				}

				panicErr := fmt.Sprintf("%v", err)
				if panicErr != "" {
					entry := logrus.WithField("path", c.Request.RequestURI).WithField("err", panicErr).WithField("stack", ginutil.GetStack()).
						WithField("costTimeMs", time.Since(now).Milliseconds())

					if isValidPanic(panicErr) {
						localIp := util.GetLocalIp()
						_ = fmt.Sprintf("服务器出现异常panic, 服务器ip: %s, path: %s，err: %v, stack: %s", localIp, c.Request.RequestURI, panicErr, ginutil.GetStack())

						// sms.SendAdminWarnGroup(msg, []string{})
						entry.Errorln("PanicRecover")
					} else {
						entry.Errorln("PanicRecover")
					}
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
