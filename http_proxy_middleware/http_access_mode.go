package http_proxy_middleware

import (
	"Gateway-Go/common/public"
	"Gateway-Go/dao"
	"Gateway-Go/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

// HttpAccessModeMiddleware 匹配接入方式 基于请求信息
func HttpAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceDetail, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}
		fmt.Println("match service", public.Obj2Json(serviceDetail))
		c.Set("service", serviceDetail)
		c.Next()
	}
}
