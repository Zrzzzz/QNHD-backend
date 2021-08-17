package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Tags backend, admin
// @Summary 获取所有的管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "管理员昵称"
// @Success 200 {object} models.Response{data=models.ListRes{list=[]models.Admin}}
// @Router /b/admin [get]
func GetAdmins(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}
	list := models.GetAdmins(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags backend, admin
// @Summary 增加管理员
// @Accept json
// @Produce json
// @Param name body string true "管理员昵称"
// @Param password body string true "管理员密码, 32位小写md5"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "参数错误"
// @Router /b/admin [post]
func AddAdmins(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	ok := r.E(&valid, "Add admin")
	if !ok {
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
// @Param name body string true "管理员昵称"
// @Param password body string true "管理员密码, 32位小写md5"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "参数错误"
// @Router /b/admin [put]
func EditAdmins(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	ok := r.E(&valid, "Edit admin")
	if !ok {
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
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "失败不返回数据"
// @Router /b/admin [delete]
func DeleteAdmins(c *gin.Context) {
	name := c.Query("name")

	valid := validation.Validation{}
	valid.Required(name, "name")

	ok := r.E(&valid, "Delete admin")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.DeleteAdmins(name)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
