package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Tags backend, admin
// @Summary 获取所有的管理员
// @Accept json
// @Produce json
// @Param token query string true "用于验证用户"
// @Success 200 {object} models.Response{data=models.ListRes{list=[]models.Admin}}
// @Router /b/admin [get]
func GetAdmins(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	maps["name"] = name
	list := models.GetAdmins(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags backend, admin
// @Summary 增加管理员
// @Accept json
// @Produce json
// @Param name query string true "管理员昵称"
// @Param password query string true "管理员密码, 32位小写md5"
// @Param token query string true "用于验证用户"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "参数错误"
// @Router /b/admin [post]
func AddAdmins(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("Add admin error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.AddAdmins(name, password)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

// @Tags backend, admin
// @Summary 修改管理员密码
// @Accept json
// @Produce json
// @Param name query string true "管理员昵称"
// @Param password query string true "管理员密码, 32位小写md5"
// @Param token query string true "用于验证用户"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "参数错误"
// @Router /b/admin [put]
func EditAdmins(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("Edit admin error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.EditAdmins(name, password)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

// @Tags backend, admin
// @Summary 删除管理员
// @Accept json
// @Produce json
// @Param name query string true "管理员昵称"
// @Param token query string true "用于验证用户"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "失败不返回数据"
// @Router /b/admin [delete]
func DeleteAdmins(c *gin.Context) {
	name := c.Query("name")

	valid := validation.Validation{}
	valid.Required(name, "name")

	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("Delete admin error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.DeleteAdmins(name)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
