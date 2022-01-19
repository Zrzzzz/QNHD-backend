package permission

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

type RightType int

const (
	ADMIN RightType = iota
	USER
)

func RightDemand(must RightType) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := r.GetUid(c)
		var err error

		if must == ADMIN {
			// 管理员验证
			_, err = models.AdminRightDemand(uid, models.UserRight{Super: true, SchAdmin: true, StuAdmin: true})

		} else {
			_, err = models.UserRightDemand(uid)
		}
		if err != nil {
			r.OK(c, e.ERROR_RIGHT, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
