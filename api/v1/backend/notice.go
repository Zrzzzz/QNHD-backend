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
// @param
// @return
// @route /b/notice
func GetNotices(c *gin.Context) {
	list, err := models.GetNotices()
	if err != nil {
		logging.Error("Get notices error: %v", err)
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
// @param content
// @return
// @route /b/notice
func AddNotice(c *gin.Context) {
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(content, "content")
	ok, verr := r.ErrorValid(&valid, "Add notice")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := make(map[string]interface{})
	maps["content"] = content

	id, err := models.AddNotice(maps)
	if err != nil {
		logging.Error("Add notices error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param id, content
// @return
// @route /b/notice
func EditNotice(c *gin.Context) {
	id := c.PostForm("id")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	valid.Required(content, "content")
	ok, verr := r.ErrorValid(&valid, "Edit notices")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	intid := util.AsUint(id)
	data := make(map[string]interface{})
	data["content"] = content
	err := models.EditNotice(intid, data)
	if err != nil {
		logging.Error("Edit notices error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param id
// @return
// @route /b/notice
func DeleteNotice(c *gin.Context) {
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Delete notices")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	intid := util.AsUint(id)
	_, err := models.DeleteNotice(intid)
	if err != nil {
		logging.Error("Delete notices error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
