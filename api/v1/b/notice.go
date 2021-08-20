package b

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

// @Tags backend, notice
// @Summary 获取公告
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=models.Notice}}
// @Router /b/notice [get]
func GetNotices(c *gin.Context) {
	data := make(map[string]interface{})
	list, err := models.GetNotices()
	if err != nil {
		logging.Error("Get notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	data["list"] = list
	data["total"] = len(list)

	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags backend, notice
// @Summary 添加公告
// @Accept json
// @Produce json
// @Param content body string true "公告内容"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.IdRes}
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/notice [post]
func AddNotices(c *gin.Context) {
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(content, "content")
	ok := r.E(&valid, "Add notices")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	maps := make(map[string]interface{})
	maps["content"] = content

	id, err := models.AddNotices(maps)
	if err != nil {
		logging.Error("Add notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags backend, notice
// @Summary 修改公告
// @Accept json
// @Produce json
// @Param id body int true "公告id"
// @Param content body string true "公告内容"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/notice [put]
func EditNotices(c *gin.Context) {
	id := c.PostForm("id")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	valid.Required(content, "content")
	ok := r.E(&valid, "Edit notices")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	intid, _ := strconv.ParseUint(id, 10, 64)
	data := make(map[string]interface{})
	data["content"] = content
	err := models.EditNotices(intid, data)
	if err != nil {
		logging.Error("Edit notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @Tags backend, notice
// @Summary 删除公告
// @Accept json
// @Produce json
// @Param id query int true "公告id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/notice [delete]
func DeleteNotices(c *gin.Context) {
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok := r.E(&valid, "Delete notices")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	intid, _ := strconv.ParseUint(id, 10, 64)
	_, err := models.DeleteNotices(intid)
	if err != nil {
		logging.Error("Delete notices error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
