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
// @route /b/notices
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

// @method [post]
// @way [formdata]
// @param sender, title, content, pub_at
// @return
// @route /b/notice
func AddNotice(c *gin.Context) {
	uid := r.GetUid(c)
	sender := c.PostForm("sender")
	title := c.PostForm("title")
	content := c.PostForm("content")
	pubAt := c.PostForm("pub_at")
	valid := validation.Validation{}
	valid.Required(sender, "sender")
	valid.MaxSize(sender, 30, "sender")
	valid.Required(title, "title")
	valid.MaxSize(title, 30, "title")
	valid.Required(content, "content")
	valid.MaxSize(content, 2000, "content")
	ok, verr := r.ErrorValid(&valid, "Add notice")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.AddNoticeToAllUsers(uid, map[string]interface{}{
		"sender":  sender,
		"title":   title,
		"content": content,
		"pub_at":  pubAt,
	})
	if err != nil {
		logging.Error("Add notice error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param sender, title, content, pub_at
// @return
// @route /b/notice/template
func AddNoticeTemplate(c *gin.Context) {
	sender := c.PostForm("sender")
	title := c.PostForm("title")
	content := c.PostForm("content")
	symbol := c.PostForm("symbol")
	valid := validation.Validation{}
	valid.Required(sender, "sender")
	valid.MaxSize(sender, 30, "sender")
	valid.Required(title, "title")
	valid.MaxSize(title, 30, "title")
	valid.Required(content, "content")
	valid.MaxSize(content, 2000, "content")
	valid.Required(symbol, "symbol")
	valid.MaxSize(symbol, 50, "symbol")
	ok, verr := r.ErrorValid(&valid, "Add notice template")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	id, err := models.AddNoticeTemplate(map[string]interface{}{
		"sender":  sender,
		"title":   title,
		"content": content,
		"symbol":  symbol,
	})
	if err != nil {
		logging.Error("Add notice template error: %v", err)
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
// @route /b/notice/modify
func EditNoticeTemplate(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.PostForm("id")
	sender := c.PostForm("sender")
	title := c.PostForm("title")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Edit notice")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	intid := util.AsUint(id)

	err := models.EditNoticeTemplate(uid, intid, map[string]interface{}{
		"sender":  sender,
		"title":   title,
		"content": content,
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
// @route /b/notice/delete
func DeleteNotice(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Delete notices")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	intid := util.AsUint(id)
	_, err := models.DeleteNoticeTemplate(uid, intid)
	if err != nil {
		logging.Error("Delete notices error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
