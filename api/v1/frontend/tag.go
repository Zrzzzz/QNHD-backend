package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

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
func GetTags(c *gin.Context) {
	name := c.Query("name")

	data := make(map[string]interface{})
	list, err := models.GetTags(name)
	if err != nil {
		logging.Error("Get tag error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.Success(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param
// @return
func GetHotTag(c *gin.Context) {
	list, err := models.GetHotTags()
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("Get hot tag error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.Success(c, e.SUCCESS, data)
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
	uid := r.GetUid(c)
	name := c.PostForm("name")
	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.MaxSize(name, 15, "name")
	ok, verr := r.E(&valid, "Add tag")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	exist, err := models.ExistTagByName(name)
	if err != nil {
		logging.Error("Add tag error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if exist {
		r.Success(c, e.ERROR_EXIST_TAG, nil)
	}
	id, err := models.AddTag(name, uid)
	if err != nil {
		logging.Error("Add tag error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.Success(c, e.SUCCESS, data)
}

// @method [delete]
// @way [query]
// @param id, uid
// @return
// @route /f/tag
func DeleteTag(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.E(&valid, "Delete tag")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	intid := util.AsUint(id)
	_, err := models.DeleteTag(intid, uid)
	if err != nil {
		logging.Error("Delete tags error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
