package backend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
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
	rType := c.Query("type")
	valid := validation.Validation{}
	valid.Required(rType, "type")
	valid.Numeric(rType, "type")
	ok, verr := r.E(&valid, "Add report")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	rTypeint := util.AsInt(rType)
	valid.Range(rTypeint, 1, 2, "type")
	ok, verr = r.E(&valid, "Add report")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	list, err := models.GetReports(rType)
	if err != nil {
		logging.Error("Get report error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
}
