package b

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

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
	list, err := models.GetAdmins(maps)
	if err != nil {
		logging.Error("Get admins error:%v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags backend, admin
// @Summary 增加管理员
// @Accept json
// @Produce json
// @Param name body string true "管理员昵称"
// @Param password body string true "管理员密码, 32位小写md5"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.IdRes}
// @Failure 400 {object} models.Response "参数错误"
// @Router /b/admin [post]
func AddAdmins(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	ok, verr := r.E(&valid, "Add admin")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	id, err := models.AddAdmins(name, password)

	if err != nil {
		logging.Error("Add admin error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
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

	ok, verr := r.E(&valid, "Edit admin")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	err := models.EditAdmins(name, password)
	if err != nil {
		logging.Error("Edit admin error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
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

	ok, verr := r.E(&valid, "Delete admin")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeleteAdmins(name)
	if err != nil {
		logging.Error("Delete admin error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
