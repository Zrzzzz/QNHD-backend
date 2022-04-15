package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
}

// @method [get]
// @way [query]
// @param name
// @return taglist
// @route /b/tags
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
// @param uid
// @return
// @route /b/tag/detail
func GetTagDetail(c *gin.Context) {
	id := c.Query("id")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "get tag detail")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	tag, err := models.GetTag(id)
	if err != nil {
		logging.Error("get tag detail error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if tag.Id == 0 {
		r.Error(c, e.ERROR_DATABASE, "无此标签")
		return
	}
	u, err := models.GetUser(map[string]interface{}{"id": tag.Uid})
	if err != nil {
		logging.Error("get tag detail error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	detail, err := twtservice.QueryUserDetail(u.Number)
	if err != nil {
		logging.Error("get tag detail error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"detail": detail})
}

// @method [get]
// @way [query]
// @param
// @return
// @route /b/tags/hot
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

// @method [delete]
// @way [query]
// @param id, uid
// @return
// @route /b/tag
func DeleteTag(c *gin.Context) {
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
	_, err := models.DeleteTagAdmin(intid)
	if err != nil {
		logging.Error("Delete tags error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
