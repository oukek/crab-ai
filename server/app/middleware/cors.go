package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			baseAllowHeaders := "x-market-user-id,x-market-client-id,if-match,if-none-match, cache-control, Origin, X-Account-Type, X-User-Session-Id, X-Machine-Session-Id, X-Client-Name,X-Client-Version,X-Client-Commit,X-Execution-Id, X-Requested-With, Content-Type, Accept, Authorization, token"

			addHeader := c.Request.Header.Get("Access-Control-Request-Headers")
			if len(addHeader) > 0 {
				baseAllowHeaders += "," + addHeader
			}
			c.Header("Access-Control-Allow-Headers", baseAllowHeaders)
			c.Header("Access-Control-Expose-Headers", "ETag, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
