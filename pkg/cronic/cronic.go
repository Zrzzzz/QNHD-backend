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
	c.AddFunc("00 00 08 * * ?", func() {
		err := models.FlushOldTagLog()
		if err != nil {
			logging.Error(err.Error())
		}
	})
	c.Start()
}

func Close() {
	c.Stop()
}
