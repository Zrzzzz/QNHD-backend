package frontend

import (
	"fmt"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/filter"
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
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param
// @return hottag
// @route /f/tag/recommend
func GetRecommendTag(c *gin.Context) {
	tagId := c.Query("id")
	valid := validation.Validation{}
	valid.Numeric(tagId, "id")
	ok, verr := r.ErrorValid(&valid, "tag id")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	tagint := -1
	if tagId != "" {
		tagint = util.AsInt(tagId)
	}
	tag, err := models.GetRecommendTag(tagint)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"tag": tag})
}

// @method [get]
// @way [query]
// @param
// @return hottagList
// @route /f/tags/hot
func GetHotTag(c *gin.Context) {
	list, err := models.GetHotTags(5)
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("Get hot tag error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param name
// @return
// @route /f/tag
func AddTag(c *gin.Context) {
	uid := r.GetUid(c)
	name := c.PostForm("name")
	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.MaxSize(name, 15, "name")
	ok, verr := r.ErrorValid(&valid, "Add tag")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	ok, s := filter.CommonFilter.Validate(name)
	if !ok {
		r.Error(c, e.INVALID_PARAMS, fmt.Sprintf("Tag触发敏感词%s", s))
		return
	}
	exist, err := models.ExistTagByName(name)
	if err != nil {
		logging.Error("Add tag error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if exist {
		r.OK(c, e.ERROR_EXIST_TAG, nil)
		return
	}
	id, err := models.AddTag(name, uid)
	if err != nil {
		logging.Error("Add tag error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, e.SUCCESS, data)
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
	ok, verr := r.ErrorValid(&valid, "Delete tag")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	intid := util.AsUint(id)
	_, err := models.DeleteTag(intid, uid)
	if err != nil {
		logging.Error("Delete tags error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
