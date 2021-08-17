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
	ok := r.E(&valid, "Get banned")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if uid != "" {
		maps["uid"] = uid
	}

	list := models.GetBanned(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags backend, banned
// @Summary 添加封号用户
// @Accept json
// @Produce json
// @Param uid body int true "用户id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/banned [post]
func AddBanned(c *gin.Context) {
	uid := c.PostForm("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok := r.E(&valid, "Add banned")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if !models.IfBannedByUid(intuid) {
		models.AddBannedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_BANNED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))

}

// @Tags backend, banned
// @Summary 删除封号用户(解禁), 此接口不使用
// @Accept json
// @Produce json
// @Param uid query int true "用户id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/banned [delete]
func DeleteBanned(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok := r.E(&valid, "Delete banned")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if models.IfBannedByUid(intuid) {
		models.DeleteBannedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_BANNED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))
}
