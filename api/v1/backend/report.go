package backend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
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

	list, err := models.GetReports()
	if err != nil {
		logging.Error("Get report error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
}
