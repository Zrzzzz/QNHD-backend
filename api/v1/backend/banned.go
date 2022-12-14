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

// @Tags backend, banned
// @Summary 获取封号用户
// @Accept json
// @Produce json
// @Param uid query int false "用户id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=[]models.Banned}}
// @Failure 400 {object} models.Response "失败不返回数据"
// @Router /b/banned [get]
func GetBanned(c *gin.Context) {
	uid := c.Query("uid")

	valid := validation.Validation{}
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "Get banned")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if uid != "" {
		maps["uid"] = uid
	}

	list, err := models.GetBanned(maps)
	if err != nil {
		logging.Error("get banned error:%v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param uid, reason
// @return
// @route /b/banned
func AddBanned(c *gin.Context) {
	doer := r.GetUid(c)
	uid := c.PostForm("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "Add banned")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	reason := c.PostForm("reason")
	intuid := util.AsUint(uid)

	code := e.SUCCESS
	ifBanned := models.IsBannedByUid(intuid)
	var err error
	if err != nil {
		logging.Error("Judging banned failed: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	var id uint64
	if !ifBanned {
		id, err = models.AddBannedByUid(intuid, util.AsUint(doer), reason)
		if err != nil {
			logging.Error("Add banned error: %v", err)
			r.Error(c, e.ERROR_DATABASE, err.Error())
			return
		}
	} else {
		code = e.ERROR_BANNED_USER
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, code, data)
}

// @method [delete]
// @way [query]
// @param uid
// @return
// @route /b/banned
func DeleteBanned(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "Delete banned")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	intuid := util.AsUint(uid)

	code := e.SUCCESS
	ifBanned := models.IsBannedByUid(intuid)
	var err error
	if err != nil {
		logging.Error("Judging banned failed: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if ifBanned {
		_, err := models.DeleteBannedByUid(intuid)
		if err != nil {
			logging.Error("Delete banned error: %v", err)
			r.Error(c, e.ERROR_DATABASE, err.Error())
			return
		}
	} else {
		code = e.ERROR_NOT_BANNED_USER
	}
	r.OK(c, code, nil)
}
