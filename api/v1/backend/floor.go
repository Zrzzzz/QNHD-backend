package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param floor_id
// @return floor
// @route /b/floor
func GetFloor(c *gin.Context) {
	floorId := c.Query("floor_id")
	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Get floorreplys")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	floor, err := models.GetFloorResponse(floorId)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"floor": floor})
}

// @method [get]
// @way [query]
// @param post_id, page=0, page_size
// @return floors
// @route /b/floors
func GetFloors(c *gin.Context) {
	postId := c.Query("post_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	ok, verr := r.ErrorValid(&valid, "Get floors")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	list, err := models.GetFloorResponses(c, postId)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page, page_size, floor_id
// @return floorlist
// @route /b/floor/replys
func GetFloorReplys(c *gin.Context) {
	floorId := c.Query("floor_id")
	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Get floorreplys")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	list, err := models.GetFloorReplyResponses(c, floorId)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [delete]
// @way [query]
// @param floor_id
// @return nil
func DeleteFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.Query("floor_id")

	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	ok, verr := r.ErrorValid(&valid, "Get floors")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeleteFloorByAdmin(uid, floorId)
	if err != nil {
		logging.Error("Delete floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
