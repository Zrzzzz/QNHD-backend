package cronic

import (
	"fmt"
	"qnhd/models"
	"qnhd/pkg/logging"
	"qnhd/request/twtservice"
	"time"

	"github.com/go-co-op/gocron"
)

func Setup() {
	err := models.FlushOldTagLog()
	if err != nil {
		logging.Error(err.Error())
	}
	err = twtservice.SaveToken()
	if err != nil {
		logging.Error(err.Error())
	}
	// 定时任务
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Week().Do(func() {
		// 清理taglog
		err := models.FlushOldTagLog()
		if err != nil {
			logging.Error(err.Error())
		}
		// 更新token
		err = twtservice.SaveToken()
		if err != nil {
			logging.Error(err.Error())
		}
		// 清理已读点赞

		// 计算帖子频率
		err = models.RefreshPostFreq()
		if err != nil {
			logging.Error(err.Error())
		}
	})

	s.StartAsync()
	fmt.Println("cron启动成功")
}

func Close() {
}
