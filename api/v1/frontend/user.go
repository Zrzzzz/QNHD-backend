package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func EditUserName(c *gin.Context) {
	uid := r.GetUid(c)
	name := c.PostForm("name")
	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.MaxSize(name, 20, "name")
	ok, verr := r.ErrorValid(&valid, "Edit user name")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.EditUserName(uid, name)
	if err != nil {
		logging.Error("edit user name error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
