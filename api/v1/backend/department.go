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
// @param name
// @return departmentList
// @route /b/departments
func GetDepartments(c *gin.Context) {
	name := c.Query("name")

	data := make(map[string]interface{})
	list, err := models.GetDepartments(name)
	if err != nil {
		logging.Error("Get department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param name, contact, contact_phone, introduction
// @return
// @route /b/department
func AddDepartment(c *gin.Context) {
	// 仅超管可用

	name := c.PostForm("name")
	introduction := c.PostForm("introduction")
	valid := validation.Validation{}
	valid.Required(name, "name")
	ok, verr := r.ErrorValid(&valid, "Add department")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	exist, err := models.ExistDepartmentByName(name)
	if err != nil {
		logging.Error("Add department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if exist {
		r.OK(c, e.ERROR_EXIST_DEPARTMENT, nil)
    return
	}
	maps := map[string]interface{}{
		"name":         name,
		"introduction": introduction,
	}
	id, err := models.AddDepartment(maps)
	if err != nil {
		logging.Error("Add department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param department_id, introduction
// @return
// @route /b/department/modify
func EditDepartment(c *gin.Context) {
	uid := r.GetUid(c)
	// 权限管理，仅学校管理
	departmentId := c.PostForm("department_id")
	introduction := c.PostForm("introduction")
	valid := validation.Validation{}
	valid.Required(departmentId, "department_id")
	valid.Numeric(departmentId, "department_id")
	ok, verr := r.ErrorValid(&valid, "Edit department")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	// 校管或者超管
	if !models.RequireRight(uid, models.UserRight{Super: true}) {
		hasRight := models.IsDepartmentHasUser(util.AsUint(uid), util.AsUint(departmentId))
		if !hasRight {
			r.OK(c, e.ERROR_RIGHT, nil)
			return
		}
	}
	err := models.EditDepartment(departmentId, introduction)
	if err != nil {
		logging.Error("Add department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param id
// @return
// @route /b/department/delete
func DeleteDepartment(c *gin.Context) {
	// 要求超管权限
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Delete department")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	_, err := models.DeleteDepartment(id)
	if err != nil {
		logging.Error("Delete departments error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
