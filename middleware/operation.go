package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"server/global"
	modelMonitor "server/model/monitor"
	"server/service"
	"server/utils"
)

var (
	operationLogService = service.GroupApp.Monitor.OperationLogService
)

type responseProxyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseProxyWriter) Write(b []byte) (int, error) {
	if r.body == nil {
		r.body = bytes.NewBufferString("")
	}
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求参数
		var reqParam []byte
		if c.Request.Method == http.MethodGet {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			reqParam, _ = json.Marshal(&m)
		} else {
			var err error
			reqParam, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.Log.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(reqParam))
			}
		}

		// 解析token
		claims, _ := utils.GetClaims(c)

		record := modelMonitor.OperationLogModel{
			Ip:        c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			UserAgent: c.Request.UserAgent(),
			ReqParam:  string(reqParam),
			UserID:    claims.ID,
			UserName:  claims.Username,
		}

		writer := responseProxyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		record.Status = c.Writer.Status()
		record.RespTime = time.Since(now).Milliseconds()
		record.RespData = writer.body.String()

		//所有操作均在这里记录日志了。
		if err := operationLogService.CreateOperationLog(record); err != nil {
			global.Log.Error("create operation record error:", zap.Error(err))
		}
	}
}
