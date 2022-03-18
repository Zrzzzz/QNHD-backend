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
// @param uid
// @return
// @route /b/blocked
func GetBlocked(c *gin.Context) {
	uid := c.Query("uid")

	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "Get blocked")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if uid != "" {
		maps["uid"] = uid
	}

	list, err := models.GetBlocked(maps)
	if err != nil {
		logging.Error("Get blocked error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param uid, last, reason
// @return
// @route /b/blocked
func AddBlocked(c *gin.Context) {
	doer := r.GetUid(c)
	uid := c.PostForm("uid")
	last := c.PostForm("last")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(last, "last")
	valid.Numeric(last, "last")
	ok, verr := r.ErrorValid(&valid, "Add blocked")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	reason := c.PostForm("reason")
	// 因为做过valid了不必考虑错误
	intuid := util.AsUint(uid)
	intlast := util.AsUint(last)
	code := e.SUCCESS
	ifBlocked := models.IsBlockedByUid(intuid)
	var id uint64
	var err error
	if !ifBlocked {
		id, err = models.AddBlockedByUid(intuid, doer, reason, uint8(intlast))
		if err != nil {
			logging.Error("Add blocked error: %v", err)
			code = e.ERROR_DATABASE
		}
	} else {
		code = e.ERROR_BLOCKED_USER
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, code, data)
}

// @method [delete]
// @way [query]
// @param uid
// @return
// @route /b/blocked
func DeleteBlocked(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "Delete blocked")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	intuid := util.AsUint(uid)

	code := e.SUCCESS
	ifBlocked := models.IsBlockedByUid(intuid)
	var err error
	if err != nil {
		logging.Error("Add blocked error: %v", err)
		code = e.ERROR_DATABASE
	}
	if ifBlocked {
		_, err := models.DeleteBlockedByUid(intuid)
		if err != nil {
			logging.Error("Delete blocked error: %v", err)
			code = e.ERROR_DATABASE
		}
	} else {
		code = e.ERROR_NOT_BLOCKED_USER
	}
	r.OK(c, code, nil)
}
