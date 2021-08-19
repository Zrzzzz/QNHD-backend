package f

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"qnhd/pkg/setting"
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
	var pageSize = setting.AppSetting.PageSize
	postId := c.Query("post_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "podsId")
	ok := r.E(&valid, "Get floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	list := models.GetFloorInPost(util.GetPage(c), pageSize, postId)

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags front, floor
// @Summary 添加楼层
// @Accept json
// @Produce json
// @Param uid body string true "用户id"
// @Param post_id body string true "帖子id"
// @Param content body string true "内容"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/floor [post]
func AddFloors(c *gin.Context) {
	postId := c.PostForm("post_id")
	uid := c.PostForm("uid")
	content := c.PostForm("content")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "podsId")
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(content, "content")
	ok := r.E(&valid, "Add floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	intpostid, _ := strconv.ParseUint(postId, 10, 64)
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	maps := map[string]interface{}{
		"uid":     intuid,
		"postId":  intpostid,
		"content": content,
	}

	models.AddFloor(maps)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
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
	valid.Numeric(postId, "podsId")
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(replyToFloor, "floorId")
	valid.Numeric(replyToFloor, "floorId")
	valid.Required(content, "content")
	ok := r.E(&valid, "Reply floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
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

	models.ReplyFloor(maps)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

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
	ok := r.E(&valid, "Get floors")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	models.DeleteFloorByUser(postId, uid, floorId)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
