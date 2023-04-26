package http_proxy_middleware

import (
	"Gateway-Go/dao"
	"Gateway-Go/middleware"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

// HttpUrlRewriteMiddleware 重写url
func HttpUrlRewriteMiddleware() gin.HandlerFunc {
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
		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			regexp, err := regexp.Compile(items[0])
			if err != nil {
				fmt.Println("regexp.Compile err:", err)
				continue
			}
			fmt.Println("before rewrite", c.Request.URL.Path)
			replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
			fmt.Println("after rewrite", c.Request.URL.Path)
		}
		c.Next()
	}
}
