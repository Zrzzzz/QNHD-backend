package main

import (
	"qnhd/api"
	"qnhd/models"
	"qnhd/pkg/cronic"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/segment"
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
	segment.Setup()
	logging.Setup()
	models.Setup()
	filter.Setup()
	api.Setup()
	cronic.Setup()

	defer models.Close()
	defer api.Close()
	defer cronic.Close()
}
