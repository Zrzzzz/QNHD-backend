package cronic

import (
	"qnhd/models"
	"qnhd/pkg/logging"

	cron "github.com/robfig/cron/v3"
)

var c *cron.Cron

func Setup() {
	err := models.FlushOldTagLog()
	if err != nil {
		logging.Error(err.Error())
	}
	// 定时任务
	c = cron.New()
	c.AddFunc("00 00 00 * * ?", func() {
		// 清理taglog
		err := models.FlushOldTagLog()
		if err != nil {
			logging.Error(err.Error())
		}
		// 清理已读楼层

		// 清理已读回复
		// 清理已读点赞
		// 清理已读通知
	})
	c.Start()
}

func Close() {
	c.Stop()
}
