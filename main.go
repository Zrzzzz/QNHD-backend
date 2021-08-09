package main

import (
	"fmt"
	"net/http"
	"qnhd/api"
	"qnhd/pkg/setting"
)

// @title QNHD API
// @version 1.0
// @schemes http
// @description 青年湖底api
// @host 116.62.107.46:7013
// @BasePath /api/v1
// @license.name Apache 2.0

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
	router := api.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
