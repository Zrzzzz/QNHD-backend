package backend

import (
	"errors"
	"fmt"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"

	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userResponse struct {
	models.User
	IsBlocked     bool   `json:"is_blocked"`
	BlockedStart  string `json:"bloced_start"`
	BlockedRemain uint64 `json:"blocked_remain"`
	BlockedOver   string `json:"blocked_over"`
	IsBanned      bool   `json:"is_banned"`
}

type userInfo struct {
	models.User
	Department models.Department `json:"department"`
}

// @method [get]
// @way [query]
// @param
// @return
// @route /b/user/info
func GetUserInfo(c *gin.Context) {
	uid := r.GetUid(c)
	user, err := models.GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		logging.Error("get user info error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	depart, err := models.GetDepartmentByUid(util.AsUint(uid))
	if !errors.Is(gorm.ErrRecordNotFound, err) {
		logging.Error("get user info error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := map[string]interface{}{
		"user_info": userInfo{User: user, Department: depart},
	}
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param uid, page, page_size
// @return userList
// @route /b/user/common
func GetCommonUser(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "get common user")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	code := e.SUCCESS
	user, err := models.GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		logging.Error("Get users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	nUser := userResponse{User: user}
	isBlocked, detail, err := models.IsBlockedByUidDetailed(user.Uid)
	if err != nil {
		logging.Error("Get users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if isBlocked {
		nUser.BlockedStart = detail.Starttime
		nUser.BlockedOver = detail.Overtime
		nUser.BlockedRemain = detail.Remain
	}
	nUser.IsBlocked = isBlocked
	nUser.IsBanned = user.Active == 0
	data := make(map[string]interface{})
	data["user"] = nUser

	r.OK(c, code, data)
}

// @method [get]
// @way [query]
// @param uid, page, page_size
// @return userList
// @route /b/users/common
func GetCommonUsers(c *gin.Context) {
	name := c.Query("number")
	code := e.SUCCESS
	list, err := models.GetCommonUsers(c, name)
	if err != nil {
		logging.Error("Get users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	retList := []userResponse{}

	for _, user := range list {
		nUser := userResponse{User: user}
		isBlocked, detail, err := models.IsBlockedByUidDetailed(user.Uid)
		if err != nil {
			logging.Error("Get users error: %v", err)
			r.Error(c, e.ERROR_DATABASE, err.Error())
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

	r.OK(c, code, data)
}

// @method [get]
// @way [query]
// @param uid, page, page_size
// @return userList
// @route /b/users/manager
func GetManagers(c *gin.Context) {
	name := c.Query("number")
	list, err := models.GetManagers(c, name)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param number, password
// @return uid
// @route /b/user
func AddUser(c *gin.Context) {
	number := c.PostForm("number")
	password := c.PostForm("password")
	phoneNumber := c.PostForm("phone_number")
	valid := validation.Validation{}
	valid.Required(number, "number")
	valid.Required(password, "password")
	valid.Required(phoneNumber, "phoneNumber")
	fmt.Println(number, password, phoneNumber)
	ok, verr := r.ErrorValid(&valid, "Add backend user")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	uid, err := models.ExistUser(number)
	if err != nil {
		logging.Error("Add user error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if uid > 0 {
		r.OK(c, e.ERROR_EXIST_USER, nil)
		return
	}

	uid, err = models.AddUser(number, password, phoneNumber)
	if err != nil {
		logging.Error("Add users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["uid"] = uid
	r.OK(c, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param new_password
// @return
// @route /b/user/modify/super
func EditUserPasswdBySuper(c *gin.Context) {
	changeid := c.PostForm("uid")
	// 超管需要修改密码
	newPass := c.PostForm("new_password")
	newPhone := c.PostForm("new_phone")
	valid := validation.Validation{}
	ok, verr := r.ErrorValid(&valid, "edit user")
	valid.Required(changeid, "uid")
	valid.Numeric(changeid, "uid")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	data := make(map[string]interface{})
	if newPass != "" {
		data["password"] = newPass
	}
	if newPhone != "" {
		data["phone_number"] = newPhone
	}
	err := models.EditUser(changeid, data)
	if err != nil {
		logging.Error("Edit users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param new_password
// @return
// @route /b/user/passwd/modify
func EditUserPasswd(c *gin.Context) {
	uid := r.GetUid(c)
	// 需要源密码
	newPass := c.PostForm("new_password")
	rawPass := c.PostForm("raw_password")
	valid := validation.Validation{}
	ok, verr := r.ErrorValid(&valid, "edit user")
	valid.Required(rawPass, "raw_password")
	valid.Required(newPass, "new_password")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	data := make(map[string]interface{})
	data["password"] = newPass
	err := models.EditUserPasswd(uid, rawPass, newPass)
	if err != nil {
		logging.Error("Edit users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param new_phone
// @return
// @route /b/user/phone/modify
func EditUserPhone(c *gin.Context) {
	uid := r.GetUid(c)
	// 需要源密码
	newPhone := c.PostForm("new_phone")
	data := make(map[string]interface{})
	data["phone_number"] = newPhone
	err := models.EditUser(uid, data)
	if err != nil {
		logging.Error("Edit users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param uid, sch_admin, stu_admin
// @return
// @route /b/user/right
func EditUserRight(c *gin.Context) {
	uid := c.PostForm("uid")
	schAdmin := c.PostForm("sch_admin")
	stuAdmin := c.PostForm("stu_admin")
	valid := validation.Validation{}
	valid.Required(schAdmin, "schAdmin")
	valid.Required(stuAdmin, "stuAdmin")
	valid.Numeric(schAdmin, "schAdmin")
	valid.Numeric(stuAdmin, "stuAdmin")
	ok, verr := r.ErrorValid(&valid, "Edit user right")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	stui := util.AsUint(schAdmin)
	schi := util.AsUint(stuAdmin)
	valid.Range(int(stui), 0, 1, "stuAdmin")
	valid.Range(int(schi), 0, 1, "schAdmin")
	ok, verr = r.ErrorValid(&valid, "Edit user right")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := map[string]interface{}{
		"sch_admin": schAdmin,
		"stu_admin": stuAdmin,
	}
	if err := models.EditUser(uid, maps); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		logging.Error("Edit user error: %v", err)
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param uid, department_id
// @return
// @route /b/user/department
func EditUserDepartment(c *gin.Context) {
	uid := c.PostForm("uid")
	departmentId := c.PostForm("department_id")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(departmentId, "departmentId")
	valid.Numeric(departmentId, "departmentId")
	ok, verr := r.ErrorValid(&valid, "Edit user right Error")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	if err := models.AddUserToDepartment(util.AsUint(uid), util.AsUint(departmentId)); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		logging.Error("Edit user error: %v", err)
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
