package http_proxy_router

import (
	"Gateway-Go/controller"
	"Gateway-Go/http_proxy_middleware"
	"Gateway-Go/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	oauth := router.Group("/oauth")
	oauth.Use(middleware.TranslationMiddleware())
	{
		controller.OAuthRegister(oauth)
	}
	router.Use(
		//接入方式的中间件
		http_proxy_middleware.HttpAccessModeMiddleware(),
		//统计和限流的中间件
		http_proxy_middleware.HTTPFlowCountMiddleware(),
		http_proxy_middleware.HTTPFlowLimitMiddleware(),
		//jwt认证,浏览统计控制中间件
		http_proxy_middleware.HTTPJwtAuthTokenMiddleware(),
		http_proxy_middleware.HTTPJwtFlowCountMiddleware(),
		http_proxy_middleware.HTTPJwtFlowLimitMiddleware(),
		//黑白名单校验中间件
		http_proxy_middleware.HttpWhiteListMiddleware(),
		http_proxy_middleware.HttpBlackListMiddleware(),
		//header头与uri，url处理中间件
		http_proxy_middleware.HttpHeaderTransferMiddleware(),
		http_proxy_middleware.HttpStripUriMiddleware(),
		http_proxy_middleware.HttpUrlRewriteMiddleware(),

		http_proxy_middleware.HttpReverseProxyMiddleware(),
	)
	return router
}
