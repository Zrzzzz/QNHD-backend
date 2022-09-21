package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

func GetSetting(c *gin.Context) {
	q := models.GetSetting()

	r.OK(c, e.SUCCESS, map[string]interface{}{"data": q})
}

func EditSetting(c *gin.Context) {
	canVisit := c.PostForm("can_visit")
	flag := canVisit == "1"
	if err := models.EditSetting(flag); err != nil {
		logging.Error("set setting error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
