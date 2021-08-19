package f

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Tags front, tag
// @Summary 获取标签
// @Accept json
// @Produce json
// @Param name query string false "标签名称"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=models.Tag}}
// @Failure 400 {object} models.Response ""
// @Router /f/tag [get]
func GetTag(c *gin.Context) {
	name := c.Query("name")

	data := make(map[string]interface{})
	list := models.GetTags(name)
	data["list"] = list
	data["total"] = len(list)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags front, tag
// @Summary 添加标签
// @Accept json
// @Produce json
// @Param name body string true "标签名称"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/tag [post]
func AddTag(c *gin.Context) {
	name := c.PostForm("name")
	valid := validation.Validation{}
	valid.Required(name, "name")
	ok := r.E(&valid, "Add tag")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	if !models.ExistTagByName(name) {
		models.AddTags(name)
		c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
	} else {
		c.JSON(http.StatusOK, r.H(e.ERROR_EXIST_TAG, nil))
	}
}

// @Tags front, tag
// @Summary 删除标签
// @Accept json
// @Produce json
// @Param id query int true "标签id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/tag [get]
func DeleteTag(c *gin.Context) {
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok := r.E(&valid, "Delete tag")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	intid, _ := strconv.ParseUint(id, 10, 64)
	models.DeleteTags(intid)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
