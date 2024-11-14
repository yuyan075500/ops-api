package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
	"io"
	"net/http"
	"ops-api/model"
)

// ExcludedPaths 不记录日志的接口
var ExcludedPaths = map[string]bool{
	"/api/v1/some-path": true,
}

// Oplog 操作日志
func Oplog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 只记录POST、PUT和DELETE请求
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodPut && c.Request.Method != http.MethodDelete {
			c.Next()
			return
		}

		// 跳过不需要记录的接口
		if ExcludedPaths[path] {
			c.Next()
			return
		}

		// 记录请求参数
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		// 获取客户端信息
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		// 获取当前登录用户的用户名
		username, _ := c.Get("username")

		// 执行请求
		c.Next()

		// 记录响应数据
		responseData, exists := c.Get("response")
		if !exists {
			responseData = ""
		}

		// 记录操作日志

		log := model.LogOplog{
			Username:      username.(string),
			Endpoint:      path,
			Method:        c.Request.Method,
			RequestParams: string(reqBodyBytes),
			ResponseData:  responseData.(string),
			ClientIP:      clientIP,
			UserAgent:     userAgent,
		}

		// 将日志保存到数据库
		if err := db.Create(&log).Error; err != nil {
			logger.Warn("保存操作日志失败: %v", err)
		}
	}
}

// getUserID 从上下文中获取用户ID
func getUserID(c *gin.Context) uint {
	// 示例代码，替换为实际的用户ID获取逻辑
	return 1
}
