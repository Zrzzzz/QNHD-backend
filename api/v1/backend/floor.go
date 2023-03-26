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
// @param floor_id
// @return floor
// @route /b/floor
func GetFloor(c *gin.Context) {
	floorId := c.Query("floor_id")
	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Get floor")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
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
	order := c.Query("order")
	onlyOwner := c.Query("only_owner")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Numeric(order, "order")
	valid.Numeric(onlyOwner, "only_owner")
	ok, verr := r.ErrorValid(&valid, "Get floors")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	args := make(map[string]interface{})
	args["order"] = order
	args["only_owner"] = onlyOwner

	list, err := models.GetFloorResponses(c, postId, args)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = models.GetCommentCount(util.AsUint(postId), false, true)
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param uid
// @return
// @route /b/floors/user
func GetUserFloors(c *gin.Context) {
	uid := c.Query("uid")
	t := c.Query("type")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(t, "type")
	valid.Numeric(t, "type")
	ok, verr := r.ErrorValid(&valid, "Get user history")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	data := make(map[string]interface{})
	list, err := models.GetUserFloorResponses(c, uid, t == "1")
	if err != nil {
		logging.Error("get user history error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
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
		r.Error(c, e.INVALID_PARAMS, verr.Error())
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
// @route /b/floor/delete
func DeleteFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.Query("floor_id")
	reason := c.Query("reason")

	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	ok, verr := r.ErrorValid(&valid, "Get floors")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	_, err := models.DeleteFloorByAdmin(uid, floorId, reason)
	if err != nil {
		logging.Error("Delete floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param post_id
// @return
// @route /b/floor/recover
func RecoverFloor(c *gin.Context) {
	floorId := c.PostForm("floor_id")

	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Recover floor")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	err := models.RecoverFloor(floorId)
	if err != nil {
		logging.Error("Recover floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param floor_id, commentable
// @return
// @route /b/floor/commentable/edit
func EditFloorCommentable(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.PostForm("floor_id")
	commentable := c.PostForm("commentable")
	valid := validation.Validation{}
	valid.Required(id, "floor_id")
	valid.Numeric(id, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Edit floor commentable")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.EditFloorCommentable(uid, id, commentable == "1")
	if err != nil {
		logging.Error("Edit floor commentable error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param floor_id, value
// @return
// @route /b/floor/value
func EditFloorValue(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.PostForm("floor_id")
	value := c.PostForm("value")
	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	valid.Required(value, "value")
	valid.Numeric(value, "value")
	ok, verr := r.ErrorValid(&valid, "edit floor value")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	valid.Range(util.AsInt(value), 0, 30000, "value")
	ok, verr = r.ErrorValid(&valid, "edit floor value")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.EditFloorValue(uid, floorId, util.AsInt(value))
	if err != nil {
		logging.Error("edit floor value error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
