package b

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

// @Tags backend, report
// @Summary 获取举报列表
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=[]models.Report}}
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/report [get]
func GetReports(c *gin.Context) {
	data := make(map[string]interface{})

	list := models.GetReports()
	data["list"] = list
	data["total"] = len(list)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}
