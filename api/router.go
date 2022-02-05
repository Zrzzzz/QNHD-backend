package api

import (
	"qnhd/api/v1/backend"
	"qnhd/api/v1/frontend"
	"qnhd/middleware/crossfield"
	"qnhd/pkg/setting"
	"qnhd/pkg/upload"

	"github.com/gin-gonic/gin"
)

func InitRouter() (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// r.Use(qnhdtls.LoadTls())
	gin.SetMode(setting.ServerSetting.RunMode)

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 解决跨域问题
	r.Use(crossfield.CrossField())
	r.StaticFS("/upload/images", gin.Dir(upload.GetImageFullPath(), false))

	avb := r.Group("/api/v1/b")
	backend.Setup(avb)
	avf := r.Group("/api/v1/f")
	frontend.Setup(avf)

	return r
}
