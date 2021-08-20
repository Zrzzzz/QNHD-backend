package f

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
	list, err := models.GetTags(name)
	if err != nil {
		logging.Error("Get tag error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
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
	exist, err := models.ExistTagByName(name)
	if err != nil {
		logging.Error("Add tag error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	if exist {
		r.R(c, http.StatusOK, e.ERROR_EXIST_TAG, nil)
	}
	id, err := models.AddTags(name)
	if err != nil {
		logging.Error("Add tag error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
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
	_, err := models.DeleteTags(intid)
	if err != nil {
		logging.Error("Delete tags error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
