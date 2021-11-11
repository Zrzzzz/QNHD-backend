package backend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /b/notice
func GetNotices(c *gin.Context) {
	data := make(map[string]interface{})
	list, err := models.GetNotices()
	if err != nil {
		logging.Error("Get notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)

	r.R(c, http.StatusOK, e.SUCCESS, data)
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
	ok, verr := r.E(&valid, "Add notice")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := make(map[string]interface{})
	maps["content"] = content

	id, err := models.AddNotices(maps)
	if err != nil {
		logging.Error("Add notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
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
	ok, verr := r.E(&valid, "Edit notices")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	intid, _ := strconv.ParseUint(id, 10, 64)
	data := make(map[string]interface{})
	data["content"] = content
	err := models.EditNotices(intid, data)
	if err != nil {
		logging.Error("Edit notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
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
	ok, verr := r.E(&valid, "Delete notices")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	intid, _ := strconv.ParseUint(id, 10, 64)
	_, err := models.DeleteNotices(intid)
	if err != nil {
		logging.Error("Delete notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}