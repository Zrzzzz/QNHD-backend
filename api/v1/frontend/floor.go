package frontend

import (
	"qnhd/crypto"
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
// @param page, page_size, post_id
// @return floorlist
// @route /f/floors
func GetFloors(c *gin.Context) {
	uid := r.GetUid(c)
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

	list, err := models.GetFloorResponsesWithUid(c, postId, uid, args)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 对楼层uid加密
	for i := range list {
		list[i].Uid = crypto.Encrypt(list[i].Uid, list[i].PostId)
		for j := range list[i].SubFloors {
			list[i].SubFloors[j].Uid = crypto.Encrypt(list[i].SubFloors[j].Uid, list[i].SubFloors[j].PostId)
		}
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param floor_id
// @return floor
// @route /f/floor
func GetFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.Query("floor_id")
	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Get floorreplys")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	floor, err := models.GetFloorResponseWithUid(floorId, uid)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 对楼层uid加密
	floor.Uid = crypto.Encrypt(floor.Uid, floor.PostId)
	for j := range floor.SubFloors {
		floor.SubFloors[j].Uid = crypto.Encrypt(floor.SubFloors[j].Uid, floor.SubFloors[j].PostId)
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"floor": floor})
}

// @method [get]
// @way [query]
// @param page, page_size, floor_id
// @return floorlist
// @route /f/floor/replys
func GetFloorReplys(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.Query("floor_id")
	valid := validation.Validation{}
	valid.Required(floorId, "floor_id")
	valid.Numeric(floorId, "floor_id")
	ok, verr := r.ErrorValid(&valid, "Get floorreplys")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	list, err := models.GetFloorReplyResponsesWithUid(c, floorId, uid)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 对楼层uid加密
	for i := range list {
		list[i].Uid = crypto.Encrypt(list[i].Uid, list[i].PostId)
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
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
	imageURLs := c.PostFormArray("images")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.MaxSize(content, 200, "content")
	valid.MaxSize(imageURLs, 1, "images")
	ok, verr := r.ErrorValid(&valid, "Add floors")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	if content == "" && len(imageURLs) == 0 {
		r.Error(c, e.INVALID_PARAMS, "缺失图片或内容")
		return
	}

	intpostid := util.AsUint(postId)
	intuid := util.AsUint(uid)
	imageURL := ""
	if len(imageURLs) > 0 {
		imageURL = imageURLs[0]
	}
	maps := map[string]interface{}{
		"uid":       intuid,
		"postId":    intpostid,
		"content":   content,
		"image_url": imageURL,
	}

	id, err := models.AddFloor(maps)
	if err != nil {
		logging.Error("Add floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param reply_to_floor, post_id, content
// @return nil
// @route /f/floor/reply
func ReplyFloor(c *gin.Context) {
	uid := r.GetUid(c)
	replyToFloor := c.PostForm("reply_to_floor")
	content := c.PostForm("content")
	imageURLs := c.PostFormArray("images")

	valid := validation.Validation{}
	valid.Required(replyToFloor, "floorId")
	valid.Numeric(replyToFloor, "floorId")
	valid.MaxSize(content, 200, "content")
	valid.MaxSize(imageURLs, 1, "images")
	ok, verr := r.ErrorValid(&valid, "Reply floors")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	if content == "" && len(imageURLs) == 0 {
		r.Error(c, e.INVALID_PARAMS, "缺失图片或内容")
		return
	}
	intuid := util.AsUint(uid)
	intfloor := util.AsUint(replyToFloor)
	imageURL := ""
	if len(imageURLs) > 0 {
		imageURL = imageURLs[0]
	}
	maps := map[string]interface{}{
		"uid":          intuid,
		"replyToFloor": intfloor,
		"content":      content,
		"image_url":    imageURL,
	}

	id, err := models.ReplyFloor(maps)
	if err != nil {
		logging.Error("Reply floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, e.SUCCESS, data)
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
	ok, verr := r.ErrorValid(&valid, "Get floors")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	_, err := models.DeleteFloorByUser(uid, floorId)
	if err != nil {
		logging.Error("Delete floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method post
// @way formdata
// @param uid, floor_id, op
// @return nil
// @route /f/floor/like
func LikeOrUnlikeFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.PostForm("floor_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.ErrorValid(&valid, "like or unlike floor")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	// 代表点赞问题
	var err error
	var cnt uint64
	if op == "1" {
		cnt, err = models.LikeFloor(floorId, uid)
	} else {
		cnt, err = models.UnlikeFloor(floorId, uid)
	}
	if err != nil {
		logging.Error("like or unlike floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method post
// @way formdata
// @param uid, floor_id, op
// @return nil
// @route /f/floor/dis
func DisOrUndisFloor(c *gin.Context) {
	uid := r.GetUid(c)
	floorId := c.PostForm("floor_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.ErrorValid(&valid, "dis or undis floor")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	// 代表点赞问题
	var err error
	var cnt uint64
	if op == "1" {
		cnt, err = models.DisFloor(floorId, uid)
	} else {
		cnt, err = models.UndisFloor(floorId, uid)
	}
	if err != nil {
		logging.Error("dis or undis floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}
