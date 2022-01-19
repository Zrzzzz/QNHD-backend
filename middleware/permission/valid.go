package permission

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

func ValidBanned() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := r.GetUid(c)
		// 查询是否封号
		if models.IsBannedByUid(util.AsUint(uid)) {
			r.OK(c, e.ERROR_BANNED_USER, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidBlocked() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := r.GetUid(c)
		// 查询是否封禁
		ok, detail, err := models.IsBlockedByUidDetailed(util.AsUint(uid))
		if err != nil {
			r.Error(c, e.ERROR_DATABASE, err.Error())
			return
		}
		if !ok {
			r.OK(c, e.ERROR_BANNED_USER, map[string]interface{}{
				"detail": detail,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
