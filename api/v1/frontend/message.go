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

type messageFloorResponseUser struct {
	Type    int          `json:"type"`
	ToFloor models.Floor `json:"to_floor"`
	Post    models.Post  `json:"post"`
	Floor   models.Floor `json:"floor"`
}

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
	var err error
	// 先获取记录
	logs, err := models.GetMessageFloors(c, uid)
	if err != nil {
		logging.Error("Get message floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	var floors = []models.Floor{}

	// 根据记录查询出楼层 由于可能有删除的楼层 这里判断一下
	for _, log := range logs {
		f, e := models.GetFloor(util.AsStrU(log.FloorId))
		if f.Id > 0 {
			floors = append(floors, f)
		}
		logging.Error(e.Error())
	}
	if err != nil {
		logging.Error("Get message floor error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 再根据楼层是否为回复帖子还是回复评论的做查询
	var list = []messageFloorResponseUser{}
	for _, f := range floors {
		var r = messageFloorResponseUser{Floor: f}
		// 搜索floor
		if f.SubTo > 0 {
			tof, e := models.GetFloor(util.AsStrU(f.ReplyTo))
			if e != nil {
				err = e
				break
			}
			r.Type = 1
			r.ToFloor = tof
		} else {
			r.Type = 0
		}
		// 搜索帖子
		p, e := models.GetPost(util.AsStrU(f.PostId))
		if e != nil {
			err = e
			break
		}
		r.Post = p
		list = append(list, r)
	}
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
	logs, err := models.GetMessagePostReplys(c, uid)
	if err != nil {
		logging.Error("Get message postReply error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	var replys = []models.PostReplyResponse{}

	// 根据记录查询出楼层
	for _, log := range logs {
		prr, e := models.GetPostReplyResponse(util.AsStrU(log.ReplyId))
		if e != nil {
			err = e
			break
		}
		replys = append(replys, prr)
	}
	if err != nil {
		logging.Error("Get message postReply error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	var list = []messageReplyResponse{}
	for _, pr := range replys {
		var r = messageReplyResponse{Reply: pr}
		// 搜索帖子
		p, e := models.GetPost(util.AsStrU(pr.PostId))
		if e != nil {
			err = e
			break
		}
		r.Post = p
		list = append(list, r)
	}
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
