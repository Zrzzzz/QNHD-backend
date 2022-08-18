package frontend

import (
	"qnhd/enums/ShareLogType"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [post]
// @way [formdata]
// @param
// @return
// @route /f/share
func ShareLog(c *gin.Context) {
	uid := r.GetUid(c)
	t := c.PostForm("type")
	oid := c.PostForm("object_id")
	valid := validation.Validation{}
	valid.Required(t, "type")
	valid.Numeric(t, "type")
	valid.Required(oid, "object_id")
	valid.Numeric(oid, "object_id")
	ok, verr := r.ErrorValid(&valid, "Add share history")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.AddShareLog(uid, util.AsUint(oid), ShareLogType.Enum(util.AsInt(t)))
	if err != nil {
		logging.Error("add share log error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
