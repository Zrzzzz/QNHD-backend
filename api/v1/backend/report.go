package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param type
// @return
// @route /b/reports
func GetReports(c *gin.Context) {
	rType := c.Query("type")
	valid := validation.Validation{}
	valid.Required(rType, "type")
	valid.Numeric(rType, "type")
	ok, verr := r.ErrorValid(&valid, "Add report")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	rTypeint := util.AsInt(rType)
	valid.Range(rTypeint, 1, 2, "type")
	ok, verr = r.ErrorValid(&valid, "Add report")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	list, err := models.GetReports(models.ReportType(rTypeint))
	if err != nil {
		logging.Error("Get report error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param id
// @return
// @route /b/report/delete
func SolveReport(c *gin.Context) {
	reportType := c.Query("type")
	id := c.Query("id")
	valid := validation.Validation{}
	valid.Required(reportType, "type")
	valid.Numeric(reportType, "type")
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "delete report")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	if err := models.SolveReport(reportType, id); err != nil {
		logging.Error("Delete report error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
