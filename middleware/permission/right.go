package permission

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

func RightDemand(right models.UserRight) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := r.GetUid(c)
		if !models.RequireRight(uid, right) {
			r.Error(c, e.ERROR_RIGHT, "")
			c.Abort()
			return
		}
		c.Next()
	}
}
