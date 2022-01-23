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
	sender := c.PostForm("sender")
	title := c.PostForm("title")
	content := c.PostForm("content")
	url := c.PostForm("url")

	valid := validation.Validation{}
	valid.Required(sender, "sender")
	valid.MaxSize(sender, 20, "sender")
	valid.Required(title, "title")
	valid.MaxSize(title, 20, "title")
	valid.Required(content, "content")
	valid.MaxSize(content, 200, "content")
	ok, verr := r.ErrorValid(&valid, "Add notice")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	id, err := models.AddNotice(map[string]interface{}{
		"sender":  sender,
		"title":   title,
		"content": content,
		"url":     url,
	})
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
	sender := c.PostForm("sender")
	title := c.PostForm("title")
	content := c.PostForm("content")
	url := c.PostForm("url")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Edit notices")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	intid := util.AsUint(id)

	err := models.EditNotice(intid, map[string]interface{}{
		"sender":  sender,
		"title":   title,
		"content": content,
		"url":     url,
	})
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
