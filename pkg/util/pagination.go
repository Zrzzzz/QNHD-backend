package util

import (
	"qnhd/pkg/setting"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	return int(pageNum-1) * pageSize, pageSize
}

// func GetPageInfo(c *gin.Context) (int, int) {
// 	result := 0
// 	var pageSize = setting.AppSetting.PageSize
// 	page, _ := strconv.ParseInt(, 10, 64)
// 	if page > 0 {
// 		result = int(page-1) * pageSize
// 	}
// 	return result,
// }
