package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [post]
// @way [formdata]
// @param
// @return
// @route /b/posttype
func AddPostType(c *gin.Context) {
	short := c.PostForm("short")
	name := c.PostForm("name")
	valid := validation.Validation{}
	valid.MaxSize(short, 5, "short")
	valid.Required(name, "name")
	valid.MaxSize(name, 10, "name")
	ok, verr := r.ErrorValid(&valid, "Add posttype")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	if err := models.AddPostType(short, name); err != nil {
		logging.Error("Add posttype error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
