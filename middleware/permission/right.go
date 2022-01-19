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
		if err := models.RequireRight(uid, right); err != nil {
			r.Error(c, e.ERROR_RIGHT, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}
