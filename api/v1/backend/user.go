package backend

import (
	"encoding/json"
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"

	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type userResponse struct {
	models.User
	IsBlocked     bool   `json:"is_blocked"`
	BlockedStart  string `json:"bloced_start"`
	BlockedRemain int    `json:"blocked_remain"`
	BlockedOver   string `json:"blocked_over"`
	IsBanned      bool   `json:"is_banned"`
}

type userInfo struct {
	models.User
	Department models.Department `json:"department"`
}

// @method [get]
// @way [query]
// @param uid
// @return
// @route /b/user/detail
func GetUserDetail(c *gin.Context) {
	doer := r.GetUid(c)
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	ok, verr := r.ErrorValid(&valid, "get user detail")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	u, err := models.GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		logging.Error("get user error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	detail, err := twtservice.QueryUserDetail(u.Number)
	if err != nil {
		logging.Error("get user error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	models.AddManagerLogWithDetail(util.AsUint(doer), u.Uid, ManagerLogType.USER_DETAIL, "")
	r.OK(c, e.SUCCESS, map[string]interface{}{"detail": detail})
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
	depart, _ := models.GetDepartmentByUid(util.AsUint(uid))
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
	nUser.IsBanned = !user.Active
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
	list, err := models.GetCommonUsers(c, map[string]interface{}{
		"is_blocked": c.Query("is_blocked"),
		"is_banned":  c.Query("is_banned"),
	})
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
		nUser.IsBanned = !user.Active
		retList = append(retList, nUser)
	}
	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param uid, page, page_size
// @return userList
// @route /b/users/manager
func GetManagers(c *gin.Context) {
	name := c.Query("user")
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
// @param user, password
// @return uid
// @route /b/user
func AddUser(c *gin.Context) {
	nickname := c.PostForm("nickname")
	password := c.PostForm("password")
	phoneNumber := c.PostForm("phone_number")
	valid := validation.Validation{}
	valid.Required(nickname, "number")
	valid.Required(password, "password")
	valid.Required(phoneNumber, "phoneNumber")
	ok, verr := r.ErrorValid(&valid, "Add backend user")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	uid, err := models.ExistUser(nickname, "")
	if err != nil {
		logging.Error("Add user error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if uid > 0 {
		r.OK(c, e.ERROR_EXIST_USER, nil)
		return
	}

	uid, err = models.AddUser(nickname, "", password, phoneNumber, "", false)
	if err != nil {
		logging.Error("Add users error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["uid"] = uid
	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param json
// @return
// @route /b/user
func AddUsers(c *gin.Context) {
	var users []models.NewUserData
	content := c.PostForm("content")

	if err := json.Unmarshal([]byte(content), &users); err != nil {
		r.Error(c, e.INVALID_PARAMS, err.Error())
		return
	}
	if err := models.AddUsers(users); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
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
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	stui := util.AsUint(schAdmin)
	schi := util.AsUint(stuAdmin)
	valid.Range(int(stui), 0, 1, "stuAdmin")
	valid.Range(int(schi), 0, 1, "schAdmin")
	ok, verr = r.ErrorValid(&valid, "Edit user right")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	maps := map[string]interface{}{
		"is_sch_admin": schAdmin,
		"is_stu_admin": stuAdmin,
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
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	if err := models.AddUserToDepartment(util.AsUint(uid), util.AsUint(departmentId)); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		logging.Error("Edit user error: %v", err)
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [get]
// @way [query]
// @param user_id
// @return
// @route /b/user/manager/delete
func DeleteManager(c *gin.Context) {
	userId := c.Query("user_id")
	valid := validation.Validation{}
	valid.Required(userId, "user_id")
	valid.Numeric(userId, "user_id")
	ok, verr := r.ErrorValid(&valid, "Delete manager Error")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	if err := models.DeleteUser(util.AsUint(userId)); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		logging.Error("Delete manager error: %v", err)
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
