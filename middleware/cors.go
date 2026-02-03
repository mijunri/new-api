package middleware

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	// 必须明确列出所有允许的请求头，通配符 "*" 在某些库版本中不生效
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Authorization",
		"Cache-Control",
		"Pragma",
		"X-Requested-With",
		"New-API-User", // 前端自定义头
		"X-New-Api-Version",
		"Priority", // Chrome 浏览器的优先级提示头
	}
	// 暴露响应头给前端
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"X-New-Api-Version",
	}
	// AllowAllOrigins 和 AllowCredentials 不能同时为 true
	// 使用 AllowOriginFunc 动态允许所有来源
	config.AllowOriginFunc = func(origin string) bool {
		return true // 允许所有来源
	}
	return cors.New(config)
}

func PoweredBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-New-Api-Version", common.Version)
		c.Next()
	}
}
