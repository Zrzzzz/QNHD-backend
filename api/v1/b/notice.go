package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
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
	list := models.GetNotices()
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags backend, notice
// @Summary 添加公告
// @Accept json
// @Produce json
// @Param content body string true "公告内容"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/notice [post]
func AddNotices(c *gin.Context) {
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(content, "content")
	ok := r.E(&valid, "Add notices")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	data := make(map[string]interface{})
	data["content"] = content
	models.AddNotices(data)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
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
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	intid, _ := strconv.ParseUint(id, 10, 64)
	data := make(map[string]interface{})
	data["content"] = content
	models.EditNotices(intid, data)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
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
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}
	intid, _ := strconv.ParseUint(id, 10, 64)
	models.DeleteNotices(intid)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
