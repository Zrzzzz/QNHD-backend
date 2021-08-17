package api

import (
	"net/http"
	"qnhd/api/v1/b"
	"qnhd/api/v1/f"
	"qnhd/middleware/crossfield"
	"qnhd/pkg/setting"
	"qnhd/pkg/upload"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "qnhd/docs"
)

func InitRouter() (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 解决跨域问题
	r.Use(crossfield.CrossField())

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	avb := r.Group("/api/v1/b")
	b.Setup(avb)
	avf := r.Group("api/v1/f")
	f.Setup(avf)

	return r
}
