package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetEpiInfos(c *gin.Context) {
	infos, cnt, err := models.GetEpiInfos(c)
	if err != nil {
		logging.Error("get epiinfos error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["data"] = infos
	data["total"] = cnt
	r.OK(c, e.SUCCESS, data)
}

func AddEpiInfoReadCount(c *gin.Context) {
	id := c.PostForm("id")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "add info read count")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	cnt, err := models.AddEpiInfoReadCount(util.AsUint(id))
	if err != nil {
		logging.Error("add info read count error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}
