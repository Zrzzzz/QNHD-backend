package api

import (
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)
	avb := r.Group("/api/v1/b")
	{
		initUsersBackend(avb)
		initBannedBackend(avb)
		initBlockedBackend(avb)
		initAdminBackend(avb)
	}
	avf := r.Group("api/v1/f")
	{
		initHashTagFront(avf)
		initUsersFront(avf)
	}

	return r
}
