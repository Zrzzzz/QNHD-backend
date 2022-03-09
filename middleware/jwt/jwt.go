package jwt

import (
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var claims *util.Claims
		var err error
		code = e.SUCCESS
		token := r.FindToken(c)
		// token不空
		if token == "" {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else {
			claims, err = util.ParseToken(token)
			if err != nil {
				// 如果检查token失败
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				// 是否过期
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			// 如果有错误
			if err != nil {
				r.OK(c, code, map[string]interface{}{"error": err.Error()})
			} else {
				r.OK(c, code, nil)
			}
			logging.Error("Auth Fail: %v, reason: %v", code, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
