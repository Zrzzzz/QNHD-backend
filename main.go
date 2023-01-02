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

// func main() {
// 	setting.Setup()
// 	if err := twtservice.NotifyFloor("haha", "3019244343"); err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("ok")
// 	}
// }

func main() {
	setting.Setup()
	segment.Setup()
	logging.Setup()
	setupModels()
	filter.Setup()
	refreshToken()
	cronic.Setup()
	api.Setup()

	defer models.Close()
	defer api.Close()
	defer cronic.Close()
}

func setupModels() {
	models.Setup(setting.EnvironmentSetting.DB_DEBUG == "1")
}

func refreshToken() {
	refreshAll := setting.EnvironmentSetting.QNHD_REFRESH_ALL == "1"
	if setting.EnvironmentSetting.QNHD_REFRESH == "1" {
		// 更新未处理的数据
		models.FlushPostsTokens(refreshAll)
		models.FlushTagsTokens(refreshAll)
	}
}
