package permission

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

type RightType int

const (
	ADMIN RightType = iota
	USER
)

func RightDemand(must RightType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		token := c.GetHeader("token")
		claims, err := util.ParseToken(token)
		if err != nil {
			c.Abort()
			return
		}

		if must == ADMIN {
			// 管理员验证
			_, err = models.AdminRightDemand(claims.Uid, models.UserRight{Super: true, SchAdmin: true, StuAdmin: true})

		} else {
			_, err = models.UserRightDemand(claims.Uid)
		}
		if err != nil {
			r.Success(c, e.ERROR_RIGHT, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
