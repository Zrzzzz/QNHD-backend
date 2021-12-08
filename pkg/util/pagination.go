package util

import (
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
)

// require content have page and page_size param
// return overnum, neednum
func HandlePaging(c *gin.Context) (int, int) {
	pageNum := 0
	pn := c.Query("page")
	if pn != "" {
		pageNum = AsInt(pn)
	}

	pageSize := setting.AppSetting.PageSize
	ps := c.Query("page_size")
	if ps != "" {
		pageSize = AsInt(ps)
	}
	return int(pageNum) * pageSize, pageSize
}
