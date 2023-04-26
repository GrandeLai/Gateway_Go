package http_proxy_middleware

import (
	"Gateway-Go/common/public"
	"Gateway-Go/dao"
	"Gateway-Go/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

// HttpStripUriMiddleware 去除uri
func HttpStripUriMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		//比如把 127.0.0.1:8080/test_http_string/abbb =>127.0.0.1:2004/aab
		//只有前缀匹配才能去除
		if serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri == 1 {
			c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule, "", 1)
		}
		c.Next()
	}
}
