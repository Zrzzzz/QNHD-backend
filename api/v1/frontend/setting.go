package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

func GetSetting(c *gin.Context) {
	q := models.GetSetting()

	r.OK(c, e.SUCCESS, map[string]interface{}{"data": q})
}
