package f

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

// @Tags front, floor
// @Summary 获取楼层
// @Accept json
// @Produce json
// @Param page query string false "分页数量"
// @Param post_id query string true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=models.Floor}}
// @Failure 400 {object} models.Response ""
// @Router /f/floors [get]
func GetFloors(c *gin.Context) {
	postId := c.Query("post_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "podsId")
	ok, verr := r.E(&valid, "Get floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	base, size := util.HandlePaging(c)
	list, err := models.GetFloorInPost(base, size, postId)
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags front, floor
// @Summary 添加楼层
// @Accept json
// @Produce json
// @Param uid body string true "用户id"
// @Param post_id body string true "帖子id"
// @Param content body string true "内容"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.IdRes}
// @Failure 400 {object} models.Response ""
// @Router /f/floor [post]
func AddFloors(c *gin.Context) {
	postId := c.PostForm("post_id")
	uid := c.PostForm("uid")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
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
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags front, floor
// @Summary 回复楼层
// @Accept json
// @Produce json
// @Param uid body string true "用户id"
// @Param reply_to_floor body string true "回复楼层id"
// @Param post_id body string true "帖子id"
// @Param content body string true "内容"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/floor/reply [post]
func ReplyFloor(c *gin.Context) {
	postId := c.PostForm("post_id")
	uid := c.PostForm("uid")
	replyToFloor := c.PostForm("reply_to_floor")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
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
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method post
// @way formdata
// @param uid, floor_id, like
// @return nil
func LikeOrUnlikeFloor(c *gin.Context) {
	uid := c.PostForm("uid")
	floorId := c.PostForm("floor_id")
	like := c.PostForm("like")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	valid.Required(like, "like")
	valid.Numeric(like, "like")
	ok, verr := r.E(&valid, "like or unlike floor")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	// 代表点赞问题
	var err error
	if like == "1" {
		err = models.LikeFloor(floorId, uid)
	} else {
		err = models.UnLikeFloor(floorId, uid)
	}
	if err != nil {
		logging.Error("like or unlike floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// func LikeFloor(c *gin.Context) {
// 	fId := c.PostForm("floor_id")
// }

// @Tags front, floor
// @Summary 删除楼层
// @Accept json
// @Produce json
// @Param uid query string true "用户id"
// @Param post_id query string true "帖子id"
// @Param floor_id query string true "楼层id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/floor [delete]
func DeleteFloor(c *gin.Context) {
	postId := c.Query("post_id")
	uid := c.Query("uid")
	floorId := c.Query("floor_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "podsId")
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	ok, verr := r.E(&valid, "Get floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeleteFloorByUser(postId, uid, floorId)
	if err != nil {
		logging.Error("Delete floor error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
