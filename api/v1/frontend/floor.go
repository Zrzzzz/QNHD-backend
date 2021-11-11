package frontend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [formdata]
// @param page, page_size, post_id
// @return floorlist
// @route /f/floors
func GetFloors(c *gin.Context) {

	postId := c.Query("post_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	ok, verr := r.E(&valid, "Get floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	base, size := util.HandlePaging(c)
	list, err := models.GetFloorInPost(base, size, postId)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param post_id, content
// @return nil
// @route /f/floor
func AddFloor(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(content, "content")
	ok, verr := r.E(&valid, "Add floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	intpostid, _ := strconv.ParseUint(postId, 10, 64)
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	maps := map[string]interface{}{
		"uid":     intuid,
		"postId":  intpostid,
		"content": content,
	}

	id, err := models.AddFloor(maps)
	if err != nil {
		logging.Error("Add floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param reply_to_floor, post_id, content
// @return nil
// @route /f/floor/reply
func ReplyFloor(c *gin.Context) {
	uid := r.GetUid(c)

	postId := c.PostForm("post_id")
	replyToFloor := c.PostForm("reply_to_floor")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(replyToFloor, "floorId")
	valid.Numeric(replyToFloor, "floorId")
	valid.Required(content, "content")
	ok, verr := r.E(&valid, "Reply floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	intpostid, _ := strconv.ParseUint(postId, 10, 64)
	intuid, _ := strconv.ParseUint(uid, 10, 64)
	intfloor, _ := strconv.ParseUint(replyToFloor, 10, 64)

	maps := map[string]interface{}{
		"uid":          intuid,
		"postId":       intpostid,
		"replyToFloor": intfloor,
		"content":      content,
	}

	_, err := models.ReplyFloor(maps)
	if err != nil {
		logging.Error("Reply floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param post_id, floor_id
// @return
// @route /f/floor
func DeleteFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.Query("floor_id")

	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	ok, verr := r.E(&valid, "Get floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeleteFloorByUser(uid, floorId)
	if err != nil {
		logging.Error("Delete floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method post
// @way formdata
// @param uid, floor_id, op
// @return nil
// @route /f/floor/likeOrUnlike
func LikeOrUnlikeFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.PostForm("floor_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.E(&valid, "like or unlike floor")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	// 代表点赞问题
	var err error
	if op == "1" {
		err = models.LikeFloor(floorId, uid)
	} else {
		err = models.UnlikeFloor(floorId, uid)
	}
	if err != nil {
		logging.Error("like or unlike floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method post
// @way formdata
// @param uid, floor_id, op
// @return nil
// @route /f/floor/disOrUndis
func DisOrUndisFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.PostForm("floor_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.E(&valid, "dis or undis floor")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	// 代表点赞问题
	var err error
	if op == "1" {
		err = models.DisFloor(floorId, uid)
	} else {
		err = models.UndisFloor(floorId, uid)
	}
	if err != nil {
		logging.Error("dis or undis floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
