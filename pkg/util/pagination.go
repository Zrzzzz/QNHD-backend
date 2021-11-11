package util

import (
	"qnhd/pkg/setting"
	"strconv"

	"github.com/gin-gonic/gin"
)

// require content have page and page_size param
// return overnum, neednum
func HandlePaging(c *gin.Context) (int, int) {
	pageNum := 0
	pn := c.Query("page")
	if pn != "" {
		pni, _ := strconv.ParseInt(pn, 10, 64)
		pageNum = int(pni)
	}

	pageSize := setting.AppSetting.PageSize
	ps := c.Query("page_size")
	if ps != "" {
		psi, _ := strconv.ParseInt(ps, 10, 64)
		pageSize = int(psi)
	}
	return int(pageNum) * pageSize, pageSize
}
