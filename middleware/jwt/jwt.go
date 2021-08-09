package jwt

import (
	"net/http"
	"qnhd/pkg/e"
	"qnhd/pkg/setting"
	"qnhd/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var claims *util.Claims
		var err error
		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err = util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})
			c.Abort()
			return
		}
		// 管理员验证
		if c.FullPath() == "/api/v1/b/admin" {
			if claims.Username != setting.AdminName || claims.Password != setting.AdminPass {
				code = e.ERROR_AUTH
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": code,
					"msg":  e.GetMsg(code),
					"data": data,
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
