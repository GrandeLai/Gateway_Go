package http_proxy_middleware

import (
	"Gateway-Go/dao"
	"Gateway-Go/middleware"
	"Gateway-Go/reverse_proxy"
	"errors"
	"github.com/gin-gonic/gin"
)

func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//创建 ReverseProxy
		//需要为每个服务都独立创建一个负载均衡器
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, errors.New("service not found"))
			c.Abort()
			return
		}
		//需要每个服务都有一个独立的连接池
		trans, err := dao.TransporterHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		//使用ReverseProxy.ServeHTTP
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return

	}
}
