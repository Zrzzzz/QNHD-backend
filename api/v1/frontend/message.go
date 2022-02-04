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

type messageReplyResponse struct {
	Post  models.Post              `json:"post"`
	Reply models.PostReplyResponse `json:"reply"`
}

// @method [get]
// @way [query]
// @param
// @return
// @route /f/message/notices
func GetMessageNotices(c *gin.Context) {
	uid := util.AsUint(r.GetUid(c))
	list, err := models.GetNotices()
	if err != nil {
		logging.Error("Get notices error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 对每个查询是否已读
	for idx, i := range list {
		if models.IsReadFloor(uid, i.Id) {
			list[idx].Read = 0
		} else {
			list[idx].Read = 1
		}
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
	list, err := models.GetMessageFloors(c, uid)
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
// @param
// @return
// @route /f/message/replys
func GetMessagePostReplys(c *gin.Context) {
	uid := r.GetUid(c)
	var err error
	// 先获取记录
	list, err := models.GetMessagePostReplys(c, uid)
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
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
	ok, verr := r.ErrorValid(&valid, "Read notice")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
	ok, verr := r.ErrorValid(&valid, "Read notice")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
