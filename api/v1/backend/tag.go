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
// @param name
// @return taglist
// @route /b/tags
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
// @route /b/tags/hot
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
	ok, verr := r.E(&valid, "Delete tag")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	intid := util.AsUint(id)
	_, err := models.DeleteTagAdmin(intid)
	if err != nil {
		logging.Error("Delete tags error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
