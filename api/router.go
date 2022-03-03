package api

import (
	"net/http"
	"net/http/httputil"
	"qnhd/api/v1/backend"
	"qnhd/api/v1/frontend"
	"qnhd/middleware/crossfield"
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
)

func avatarReverse(c *gin.Context) {
	realPath := c.Param("p")
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:7014"
		req.URL.Path = realPath
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func InitRouter() (r *gin.Engine) {
	gin.SetMode(setting.ServerSetting.RunMode)
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 解决跨域问题
	r.Use(crossfield.CrossField())
	// 头像服务转发
	r.GET("/avatar/*p", avatarReverse)
	avb := r.Group("/api/v1/b")
	backend.Setup(avb)
	avf := r.Group("/api/v1/f")
	frontend.Setup(avf)

	return r
}
