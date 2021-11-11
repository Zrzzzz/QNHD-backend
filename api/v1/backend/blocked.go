package backend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"strconv"

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
	valid.Numeric(uid, "uid")
	ok, verr := r.E(&valid, "Get blocked")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)

	r.R(c, http.StatusOK, e.SUCCESS, data)
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
	ok, verr := r.E(&valid, "Add blocked")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	reason := c.PostForm("reason")
	// 因为做过valid了不必考虑错误
	intuid, _ := strconv.ParseUint(uid, 10, 64)
	intlast, _ := strconv.ParseUint(last, 10, 8)
	code := e.SUCCESS
	ifBlocked, err := models.IfBlockedByUid(intuid)
	if err != nil {
		logging.Error("Add blocked error: %v", err)
		code = e.ERROR_DATABASE
	}
	var id uint64
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
	r.R(c, http.StatusOK, code, data)
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
	ok, verr := r.E(&valid, "Delete blocked")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := e.SUCCESS
	ifBlocked, err := models.IfBlockedByUid(intuid)
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
	r.R(c, http.StatusOK, code, nil)
}
