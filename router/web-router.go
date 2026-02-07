package router

import (
	"embed"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	router.Use(middleware.Cache())
	router.Use(static.Serve("/", common.EmbedFolder(buildFS, "web/dist")))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// 处理尾部斜杠：如果路径以 / 结尾且不是根路径，重定向到不带 / 的路径
		if len(path) > 1 && strings.HasSuffix(path, "/") {
			// 构建不带尾部斜杠的 URL（保留 query 参数）
			newPath := path[:len(path)-1]
			if c.Request.URL.RawQuery != "" {
				newPath = newPath + "?" + c.Request.URL.RawQuery
			}
			c.Redirect(http.StatusMovedPermanently, newPath)
			return
		}
		// 原有的 API/v1/assets 路由处理
		if strings.HasPrefix(c.Request.RequestURI, "/v1") || strings.HasPrefix(c.Request.RequestURI, "/api") || strings.HasPrefix(c.Request.RequestURI, "/assets") {
			controller.RelayNotFound(c)
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
	})
}
