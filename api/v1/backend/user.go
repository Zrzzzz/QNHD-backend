package backend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"strings"

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

// @Tags backend, user
// @Summary 获取用户列表
// @Accept json
// @Produce json
// @Param uid query int false "用户id"
// @Param email query string false "用户邮箱"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=UserResponse}}
// @Failure 400 {object} models.Response "无效的参数"
// @Router /b/user [get]
func GetUsers(c *gin.Context) {
	uid := c.Query("uid")
	email := c.Query("email")

	valid := validation.Validation{}
	valid.Numeric(uid, "id")
	if email != "" {
		valid.Email(email, "email")
	}
	ok, verr := r.E(&valid, "Get user")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if uid != "" {
		maps["uid"] = uid
	}
	if email != "" {
		maps["email"] = email
	}

	code := e.SUCCESS

	list, err := models.GetUsers(maps)
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
		nUser.IsBanned = user.Status == 0
		retList = append(retList, nUser)
	}
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, code, data)
}

// @Tags backend, user
// @Summary 添加用户，不应使用此接口
// @Accept json
// @Produce json
// @Param email body string true "用户邮箱, 必须10位学号的tju邮箱"
// @Param password body string true "用户密码, 32位小写md5"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.IdRes}
// @Failure 400 {object} models.Response "无效的参数"
// @Router /b/user [post]
func AddUsers(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if !checkMail(email) {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	uid, err := models.ExistUser(email)
	if err != nil {
		logging.Error("Add user error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if uid > 0 {
		c.JSON(http.StatusOK, r.H(e.ERROR_EXIST_EMAIL, nil))
		return
	}

	id, err := models.AddUser(email, password)
	if err != nil {
		logging.Error("Add users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags backend, user
// @Summary 修改用户数据
// @Accept json
// @Produce json
// @Param email body string true "用户邮箱"
// @Param new_password body string true "用户新密码, 32位小写md5"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效的参数"
// @Router /b/user [put]
func EditUsers(c *gin.Context) {
	email := c.PostForm("email")
	newPass := c.PostForm("new_password")

	uid, err := models.ExistUser(email)
	if err != nil {
		logging.Error("Edit user error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if !(uid > 0) {
		r.R(c, http.StatusOK, e.ERROR_NOT_EXIST_EMAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["password"] = newPass
	err = models.EditUser(email, data)
	if err != nil {
		logging.Error("Edit users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @Tags backend, user
// @Summary 删除用户，不应使用此接口
// @Accept json
// @Produce json
// @Param email query string true "用户邮箱"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效的参数"
// @Router /b/user [delete]
func DeleteUsers(c *gin.Context) {
	email := c.Query("email")
	uid, err := models.ExistUser(email)
	if err != nil {
		logging.Error("Delete users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	if !(uid > 0) {
		r.R(c, http.StatusOK, e.ERROR_NOT_EXIST_EMAIL, nil)
		return
	}
	err = models.DeleteUser(email)
	if err != nil {
		logging.Error("Delete users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

func checkMail(email string) bool {
	s := strings.Split(email, "@")
	if len(s) != 2 || len(s[0]) != 10 || s[1] != "tju.edu.cn" {
		return false
	}
	return true
}
