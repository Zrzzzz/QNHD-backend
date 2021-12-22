package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

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
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.Success(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param name, contact, contact_phone, introduction
// @return
// @route /b/department
func AddDepartment(c *gin.Context) {
	// 仅超管可用
	uid := r.GetUid(c)
	ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
	if err != nil {
		logging.Error("Check right error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if !ok {
		r.Success(c, e.ERROR_RIGHT, nil)
		return
	}
	name := c.PostForm("name")
	introduction := c.PostForm("introduction")
	valid := validation.Validation{}
	valid.Required(name, "name")
	ok, verr := r.E(&valid, "Add department")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	exist, err := models.ExistDepartmentByName(name)
	if err != nil {
		logging.Error("Add department error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if exist {
		r.Success(c, e.ERROR_EXIST_DEPARTMENT, nil)
	}
	maps := map[string]interface{}{
		"name":         name,
		"introduction": introduction,
	}
	id, err := models.AddDepartment(maps)
	if err != nil {
		logging.Error("Add department error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.Success(c, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param department_id, introduction
// @return
// @route /b/department/modify
func EditDepartment(c *gin.Context) {
	uid := r.GetUid(c)
	// 权限管理，仅学校管理
	ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
	if err != nil {
		logging.Error("Check right error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if !ok {
		r.Success(c, e.ERROR_RIGHT, nil)
		return
	}
	departmentId := c.PostForm("department_id")
	introduction := c.PostForm("introduction")
	valid := validation.Validation{}
	valid.Required(departmentId, "department_id")
	valid.Numeric(departmentId, "department_id")
	ok, verr := r.E(&valid, "Edit department")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	hasRight, err := models.IsUserInDepartment(uid, departmentId)
	if err != nil {
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if !hasRight {
		r.Success(c, e.ERROR_RIGHT, nil)
		return
	}
	err = models.EditDepartment(departmentId, introduction)
	if err != nil {
		logging.Error("Add department error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param id
// @return
// @route /b/department/delete
func DeleteDepartment(c *gin.Context) {
	// 要求超管权限
	uid := r.GetUid(c)
	ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
	if err != nil {
		logging.Error("Check right error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if !ok {
		r.Success(c, e.ERROR_RIGHT, nil)
		return
	}
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.E(&valid, "Delete department")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err = models.DeleteDepartment(id)
	if err != nil {
		logging.Error("Delete departments error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
