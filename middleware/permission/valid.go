package permission

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

// 前端是否可以访问
func FrontCanVisit() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := models.GetSetting()
		if !q.FrontVisit {
			c.JSON(http.StatusNotFound, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

// 验证封号
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

// 验证禁言
func ValidBlocked() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := r.GetUid(c)
		// 查询是否封禁
		blocked, detail, err := models.IsBlockedByUidDetailed(util.AsUint(uid))
		if err != nil {
			r.Error(c, e.ERROR_DATABASE, err.Error())
			c.Abort()
			return
		}
		if blocked {
			r.OK(c, e.ERROR_BLOCKED_USER, map[string]interface{}{
				"detail": detail,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
