package jwt

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ADMIN = 0x53
	USER  = 0x16
)

func JWT(must int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var claims *util.Claims
		var err error
		code = e.SUCCESS
		token := c.GetHeader("token")
		if token == "" {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else {
			claims, err = util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				// 时间判断
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			} else {
				var ok bool
				if must == ADMIN {
					// 管理员验证
					ok, err = models.AdminRightDemand(claims.Uid, models.UserRight{Super: true, SchAdmin: true, StuAdmin: true})

				} else {
					ok, err = models.UserRightDemand(claims.Uid)
				}
				if !ok {
					code = e.ERROR_RIGHT
				}
			}
		}

		if err != nil {
			code = e.ERROR_DATABASE
		}

		if code != e.SUCCESS {
			if err != nil {
				r.R(c, http.StatusUnauthorized, code, map[string]interface{}{"error": err.Error()})
			} else {
				r.R(c, http.StatusUnauthorized, code, nil)
			}
			logging.Error("Auth Fail: %v, reason: %v", code, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
