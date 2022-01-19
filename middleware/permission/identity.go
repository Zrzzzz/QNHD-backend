package permission

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

type IdentityType int

const (
	ADMIN IdentityType = iota
	USER
)

func IdentityDemand(must IdentityType) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := r.GetUid(c)
		var err error

		if must == ADMIN {
			// 管理员验证
			err = models.RequireAdmin(uid)

		} else {
			err = models.RequireUser(uid)
		}
		if err != nil {
			r.Error(c, e.ERROR_RIGHT, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
