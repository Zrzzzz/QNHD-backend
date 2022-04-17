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

// @method [get]
// @way [query]
// @param
// @return
// @route /f/message/notices
func GetMessageNotices(c *gin.Context) {
	uid := util.AsUint(r.GetUid(c))
	list, err := models.GetUnreadNotices(c, uid)
	if err != nil {
		logging.Error("Get notices error: %v", err)
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
// @param uid
// @return
// @route /f/message/floors
func GetMessageFloors(c *gin.Context) {
	uid := r.GetUid(c)
	list, err := models.GetUnreadFloors(c, uid)
	if err != nil {
		logging.Error("Get message floor error: %v", err)
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
// @param uid
// @return
// @route /f/message/likes
func GetMessageLikes(c *gin.Context) {
	uid := r.GetUid(c)
	list, err := models.GetUnreadLikes(c, uid)
	if err != nil {
		logging.Error("Get message likes error: %v", err)
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
// @param
// @return
// @route /f/message/replys
func GetMessagePostReplys(c *gin.Context) {
	uid := r.GetUid(c)
	var err error
	// 先获取记录
	list, err := models.GetUnreadPostReplys(c, uid)
	if err != nil {
		logging.Error("Get message postReply error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param id
// @return
// @route /f/message/notice/read
func ReadNotice(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.PostForm("id")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Read notice")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.ReadNotice(util.AsUint(uid), util.AsUint(id))
	if err != nil {
		logging.Error("Read notice error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param ids
// @return
// @route /f/message/notice/delete
func DeleteMessageNotices(c *gin.Context) {
	uid := r.GetUid(c)
	ids := c.QueryArray("ids")
	err := models.DeleteMessageNotices(uid, ids)
	if err != nil {
		logging.Error("Delete notices error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param id
// @return
// @route /f/message/floor/read
func ReadFloor(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.PostForm("id")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Read floor")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.ReadFloor(util.AsUint(uid), util.AsUint(id))
	if err != nil {
		logging.Error("Read floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param postId, id
// @return
// @route /f/message/reply/read
func ReadReply(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.PostForm("id")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Read reply")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.ReadPostReply(util.AsUint(uid), util.AsUint(id))
	if err != nil {
		logging.Error("Read reply error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param type, id
// @return
// @route /f/message/like/read
func ReadLike(c *gin.Context) {
	uid := r.GetUid(c)
	likeType := c.PostForm("type")
	id := c.PostForm("id")
	valid := validation.Validation{}
	valid.Required(likeType, "type")
	valid.Numeric(likeType, "type")
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Read like")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	typeint := util.AsInt(likeType)
	valid.Range(typeint, 0, 1, "type")
	ok, verr = r.ErrorValid(&valid, "Read like")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.ReadLike(util.AsUint(uid), models.LikeType(typeint), util.AsUint(id))
	if err != nil {
		logging.Error("Read like error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [get]
// @way [query]
// @param
// @return
// @route /f/message/count
func GetMessageCount(c *gin.Context) {
	uid := r.GetUid(c)
	cnt, err := models.GetMessageCount(uid)
	if err != nil {
		logging.Error("Get message count error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{
		"count": cnt,
	})
}

// @method [post]
// @way [formdata]
// @param
// @return
// @route /f/message/all
func ReadAllMessage(c *gin.Context) {
	uid := r.GetUid(c)
	err := models.ReadAllMessage(util.AsUint(uid))
	if err != nil {
		logging.Error("Read reply error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
