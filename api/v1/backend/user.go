package backend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"

	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	models.User
	IsBlocked     bool   `json:"is_blocked"`
	BlockedStart  string `json:"bloced_start"`
	BlockedRemain uint64 `json:"blocked_remain"`
	BlockedOver   string `json:"blocked_over"`
	IsBanned      bool   `json:"is_banned"`
}

// @method [get]
// @way [query]
// @param uid, page, page_size
// @return userList
// @route /b/users/common
func GetCommonUsers(c *gin.Context) {
	uid := c.Query("uid")
	over, pageSize := util.HandlePaging(c)

	code := e.SUCCESS
	list, err := models.GetCommonUsers(uid, over, pageSize)
	if err != nil {
		logging.Error("Get users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []UserResponse{}

	for _, user := range list {
		nUser := UserResponse{User: user}
		isBlocked, detail, err := models.IfBlockedByUidDetailed(user.Uid)
		if err != nil {
			logging.Error("Get users error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		if isBlocked {
			nUser.BlockedStart = detail.Starttime
			nUser.BlockedOver = detail.Overtime
			nUser.BlockedRemain = detail.Remain
		}
		nUser.IsBlocked = isBlocked
		nUser.IsBanned = user.Active == 0
		retList = append(retList, nUser)
	}
	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, code, data)
}

// @method [get]
// @way [query]
// @param uid, page, page_size
// @return userList
// @route /b/users/all
func GetAllUsers(c *gin.Context) {
	uid := r.GetUid(c)
	// 权限控制
	ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
	if err != nil {
		logging.Error("Edit user right error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	// 需要管理员权限
	if !ok {
		r.R(c, http.StatusOK, e.ERROR_RIGHT, nil)
		return
	}

	uid = c.Query("uid")
	over, pageSize := util.HandlePaging(c)

	code := e.SUCCESS
	list, err := models.GetAllUsers(uid, over, pageSize)
	if err != nil {
		logging.Error("Get users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []UserResponse{}

	for _, user := range list {
		nUser := UserResponse{User: user}
		isBlocked, detail, err := models.IfBlockedByUidDetailed(user.Uid)
		if err != nil {
			logging.Error("Get users error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		if isBlocked {
			nUser.BlockedStart = detail.Starttime
			nUser.BlockedOver = detail.Overtime
			nUser.BlockedRemain = detail.Remain
		}
		nUser.IsBlocked = isBlocked
		nUser.IsBanned = user.Active == 0
		retList = append(retList, nUser)
	}
	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, code, data)
}

// @method [post]
// @way [formdata]
// @param number, password
// @return uid
// @route /b/user
func AddUser(c *gin.Context) {
	number := c.PostForm("number")
	password := c.PostForm("password")

	uid, err := models.ExistUser(number)
	if err != nil {
		logging.Error("Add user error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if uid > 0 {
		r.R(c, http.StatusOK, e.ERROR_EXIST_USER, nil)
		return
	}

	uid, err = models.AddUser(number, password)
	if err != nil {
		logging.Error("Add users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["uid"] = uid
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param new_password
// @return
// @route /b/user
func EditUser(c *gin.Context) {
	uid := r.GetUid(c)
	changeid := c.PostForm("uid")
	// 超管需要修改密码
	if changeid != "" {
		ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
		if err != nil {
			logging.Error("check right error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		if !ok {
			r.R(c, http.StatusOK, e.ERROR_RIGHT, nil)
			return
		}
	}
	newPass := c.PostForm("new_password")

	data := make(map[string]interface{})
	data["password"] = newPass
	var err error
	// 看是超管更改还是自己更改
	if changeid != "" {
		err = models.EditUser(changeid, data)
	} else {
		err = models.EditUser(uid, data)
	}
	if err != nil {
		logging.Error("Edit users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param uid, sch_admin, stu_admin
// @return
// @route /b/user/right
func EditUserRight(c *gin.Context) {
	uid := r.GetUid(c)
	// 权限控制
	ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
	if err != nil {
		logging.Error("Edit user right error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	// 需要超管才能修改
	if !ok {
		r.R(c, http.StatusOK, e.ERROR_RIGHT, nil)
		return
	}

	uid = c.PostForm("uid")
	schAdmin := c.PostForm("sch_admin")
	stuAdmin := c.PostForm("stu_admin")
	valid := validation.Validation{}
	valid.Numeric(schAdmin, "schAdmin")
	valid.Numeric(stuAdmin, "stuAdmin")
	ok, verr := r.E(&valid, "Edit user right")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	stui := util.AsUint(schAdmin)
	schi := util.AsUint(stuAdmin)
	valid.Range(int(stui), 0, 1, "stuAdmin range")
	valid.Range(int(schi), 0, 1, "schAdmin range")
	ok, verr = r.E(&valid, "Edit user right")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := map[string]interface{}{
		"sch_admin": schAdmin,
		"stu_admin": stuAdmin,
	}
	if err := models.EditUser(uid, maps); err != nil {
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		logging.Error("Edit user error: %v", err)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param uid, department_id
// @return
// @route /b/user/department
func EditUserDepartment(c *gin.Context) {
	uid := r.GetUid(c)
	// 权限控制
	ok, err := models.AdminRightDemand(uid, models.UserRight{Super: true})
	if err != nil {
		logging.Error("Edit user right error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	// 需要超管才能修改
	if !ok {
		r.R(c, http.StatusOK, e.ERROR_RIGHT, nil)
		return
	}

	uid = c.PostForm("uid")
	departmentId := c.PostForm("department_id")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(departmentId, "departmentId")
	valid.Numeric(departmentId, "departmentId")
	ok, verr := r.E(&valid, "Edit user right Error")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	if err := models.AddUserToDepartment(util.AsUint(uid), util.AsUint(departmentId)); err != nil {
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		logging.Error("Edit user error: %v", err)
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
