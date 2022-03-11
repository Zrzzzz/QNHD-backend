package main

import (
	"os"
	"qnhd/api"
	"qnhd/models"
	"qnhd/pkg/cronic"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/segment"
	"qnhd/pkg/setting"
)

func main() {
	setting.Setup()
	segment.Setup()
	logging.Setup()
	models.Setup()
	filter.Setup()
	refreshToken()
	api.Setup()
	cronic.Setup()

	defer models.Close()
	defer api.Close()
	defer cronic.Close()
}

func refreshToken() {
	shouldRefresh := os.Getenv("QNHD_REFRESH")
	if shouldRefresh == "1" {
		// 更新未处理的数据
		models.FlushPostsTokens(false)
		models.FlushTagsTokens(false)
	}
}
