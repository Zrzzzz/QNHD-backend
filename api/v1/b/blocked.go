package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
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
	ok := r.E(&valid, "Get blocked")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if uid != "" {
		maps["uid"] = uid
	}

	list := models.GetBlocked(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags backend, blocked
// @Summary 添加禁言用户
// @Accept json
// @Produce json
// @Param uid body int true "用户id"
// @Param last body int true "持续天数 0<?<=30"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
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
	ok := r.E(&valid, "Add blocked")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}
	// 因为做过valid了不必考虑错误
	intuid, _ := strconv.ParseUint(uid, 10, 64)
	intlast, _ := strconv.ParseUint(last, 10, 8)
	code := 0
	if !models.IfBlockedByUid(intuid) {
		models.AddBlockedByUid(intuid, uint8(intlast))
		code = e.SUCCESS
	} else {
		code = e.ERROR_BLOCKED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))

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
	ok := r.E(&valid, "Delete blocked")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if models.IfBlockedByUid(intuid) {
		models.DeleteBlockedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_BLOCKED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))
}
