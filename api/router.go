package api

import (
	"qnhd/middleware/jwt"
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "qnhd/docs"
)

func InitRouter() (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gin.SetMode(setting.RunMode)

	avb := r.Group("/api/v1/b")
	{
		initAuthBackend(avb)
		avb.Use(jwt.JWT())
		// 这之后的需要jwt验证
		initUsersBackend(avb)
		initBannedBackend(avb)
		initBlockedBackend(avb)
		initAdminBackend(avb)
	}
	avf := r.Group("api/v1/f")
	{
		initAuthFront(avf)
		avb.Use(jwt.JWT())
		// 这之后的需要jwt验证
		initHashTagFront(avf)
		initUsersFront(avf)
	}

	return r
}
