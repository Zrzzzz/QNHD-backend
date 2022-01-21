package main

import (
	"fmt"
	"net/http"
	"qnhd/api"
	"qnhd/models"
	"qnhd/pkg/cronic"
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
	setting.Setup()
	logging.Setup()
	models.Setup()
	router := api.InitRouter()
	cronic.Setup()
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

	defer models.Close()
	defer cronic.Close()
}
