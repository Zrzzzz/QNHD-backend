package util

import (
	"qnhd/pkg/setting"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	if page > 0 {
		result = int(page-1) * setting.AppSetting.PageSize
	}
	return result
}
