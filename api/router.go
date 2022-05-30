package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"qnhd/api/v1/backend"
	"qnhd/api/v1/frontend"
	"qnhd/middleware/crossfield"
	"qnhd/middleware/safety"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
)

var s *http.Server

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

func initRouter() (r *gin.Engine) {
	gin.SetMode(setting.ServerSetting.RunMode)
	r = gin.New()

	r.Use(logging.GinLogger())
	r.Use(gin.Recovery())

	// 解决跨域问题
	r.Use(crossfield.CrossField())
	// 解决安全问题
	r.Use(safety.Safety())
	// 头像服务转发
	r.GET("/avatar/*p", avatarReverse)

	avb := r.Group("/api/v1/b")
	backend.Setup(avb)
	avf := r.Group("/api/v1/f")
	frontend.Setup(avf)
	r.Static("src", "pages/jump")

	r.Static("message", "pages/message")
	return r
}

func Setup() {
	router := initRouter()
	s = &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
		// TLSConfig:      tlscfg,
	}
	if setting.EnvironmentSetting.RELEASE == "1" {
		s.ListenAndServe()
	} else {
		s.ListenAndServeTLS("cert/cert.pem", "cert/cert.key")
	}
}

func Close() {
	s.Close()
}
