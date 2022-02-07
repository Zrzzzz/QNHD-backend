package api

import (
	"qnhd/api/v1/backend"
	"qnhd/api/v1/frontend"
	"qnhd/middleware/crossfield"
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() (r *gin.Engine) {
	gin.SetMode(setting.ServerSetting.RunMode)
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 解决跨域问题
	r.Use(crossfield.CrossField())

	avb := r.Group("/api/v1/b")
	backend.Setup(avb)
	avf := r.Group("/api/v1/f")
	frontend.Setup(avf)

	return r
}
