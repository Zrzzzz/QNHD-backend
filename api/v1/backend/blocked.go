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

// @Tags backend, blocked
// @Summary 获取禁言用户
// @Accept json
// @Produce json
// @Param uid query string false "用户id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=[]models.Blocked}}
// @Failure 400 {object} models.Response "失败不返回数据"
// @Router /b/blocked [get]
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

// @Tags backend, blocked
// @Summary 添加禁言用户
// @Accept json
// @Produce json
// @Param uid body int true "用户id"
// @Param last body int true "持续天数 0<?<=30"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.IdRes}
// @Failure 400 {object} models.Response "失败不返回数据"
// @Router /b/blocked [post]
func AddBlocked(c *gin.Context) {
	uid := c.PostForm("uid")
	last := c.PostForm("last")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(last, "last")
	valid.Numeric(last, "last")
	ok, verr := r.E(&valid, "Add blocked")
	reason := c.PostForm("reason")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
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
		id, err = models.AddBlockedByUid(intuid, reason, uint8(intlast))
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

// @Tags backend, blocked
// @Summary 删除禁言用户
// @Accept json
// @Produce json
// @Param uid query string true "用户id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "失败不返回数据"
// @Router /b/blocked [delete]
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
