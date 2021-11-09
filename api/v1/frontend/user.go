package frontend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	sender "qnhd/pkg/email"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkMail(email string) bool {
	s := strings.Split(email, "@")
	if len(s) != 2 || len(s[0]) == 10 || s[1] != "tju.edu.cn" {
		return false
	}
	return true
}

// @Tags front, user
// @Summary 前端新建用户，没有带code，会发送邮件，带code的，会验证code并且新建用户
// @Accept json
// @Produce json
// @Param email query string true "tju邮箱，必须是10位学号"
// @Param password query string true "密码，32位md5"
// @Param code query string false "邮箱验证码"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.IdRes}
// @Failure 400 {object} models.Response ""
// @Router /f/user [post]
func AddUsers(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	code := c.Query("code")
	rightCode := sender.Code[email]

	if email == "" || checkMail(email) {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	uid, err := models.ExistUser(email)
	if err != nil {
		logging.Error("Add users error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if uid > 0 {
		r.R(c, http.StatusOK, e.ERROR_EXIST_EMAIL, nil)
		return
	}

	if code == "" {
		// 进行邮箱验证码发送
		sender.SendEmail(email,
			func() {
				r.R(c, http.StatusOK, e.SUCCESS, nil)
			},
			func() {
				r.R(c, http.StatusOK, e.ERROR_SEND_EMAIL, nil)
			})
	} else {
		if rightCode == code {
			// 验证码正确
			id, err := models.AddUser(email, password)
			if err != nil {
				logging.Error("Add user error: %v", err)
				r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
				return
			}
			data := make(map[string]interface{})
			data["id"] = id
			r.R(c, http.StatusOK, e.SUCCESS, data)
		} else {
			// 验证码错误
			r.R(c, http.StatusOK, e.ERROR_EMAIL_CODE_CHECK, map[string]interface{}{"error": err.Error()})
		}
	}
}

// @Tags front, user
// @Summary 更改密码
// @Accept json
// @Produce json
// @Param email query string true "用户邮箱"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/user [put]
func EditUsers(c *gin.Context) {
	email := c.Query("email")
	oldPass := c.Query("old_password")
	newPass := c.Query("new_password")

	uid, err := models.CheckUser(email, oldPass)
	if err != nil {
		logging.Error("Edit user error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	if uid > 0 {
		data := make(map[string]interface{})
		data["password"] = newPass
		err := models.EditUser(email, data)
		if err != nil {
			logging.Error("Edit user error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		r.R(c, http.StatusOK, e.SUCCESS, nil)
	} else {
		r.R(c, http.StatusOK, e.ERROR_AUTH, nil)
	}
}

// get
