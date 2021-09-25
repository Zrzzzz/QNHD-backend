package main

import (
	"fmt"
	"net/http"
	"qnhd/api"
	"qnhd/models"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"
)

// @title QNHD API
// @version 1.0
// @schemes http
// @description 青年湖底api
// @host 116.62.107.46:7013
// @BasePath /api/v1
// @license.name Apache 2.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
func main() {
	// endless.DefaultReadTimeOut = setting.ReadTimeout
	// endless.DefaultWriteTimeOut = setting.WriteTimeout
	// endless.DefaultMaxHeaderBytes = 1 << 20
	// endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	// s := endless.NewServer(endPoint, api.InitRouter())
	// s.BeforeBegin = func(add string) {
	// 	logging.Info("Actual pid is %d", syscall.Getpid())
	// }

	// if err := s.ListenAndServe(); err != nil {
	// 	logging.Error("Server err: %v", err)
	// }
	setting.Setup()
	logging.Setup()
	models.Setup()

	router := api.InitRouter()
	// tlscfg := api.InitTlsConfig()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
		// TLSConfig:      tlscfg,
	}
	s.ListenAndServeTLS("cert/5193613_zrzz.site.pem", "cert/5193613_zrzz.site.key")
	// s.ListenAndServe()
	defer models.CloseDB()
}
