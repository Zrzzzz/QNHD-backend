package f

import (
	"fmt"
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	sender "qnhd/pkg/email"
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
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response ""
// @Router /f/user [post]
func AddUsers(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	code := c.Query("code")
	rightCode := fmt.Sprintf("%06d", sender.Code[email]%1000000)

	if email == "" || checkMail(email) {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	} else if models.ExistUserByEmail(email) {
		c.JSON(http.StatusOK, r.H(e.ERROR_EXIST_EMAIL, nil))
		return
	}

	if code == "" {
		// 进行邮箱验证码发送
		sender.SendEmail(email,
			func() {
				c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
			},
			func() {
				c.JSON(http.StatusOK, r.H(e.ERROR_SEND_EMAIL, nil))
			})
	} else {
		if rightCode == code {
			// 验证码正确
			models.AddUser(email, password)
			c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
		} else {
			// 验证码错误
			c.JSON(http.StatusOK, r.H(e.ERROR_AUTH, nil))
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

	if models.ValidUser(email, oldPass) {
		data := make(map[string]interface{})
		data["password"] = newPass
		models.EditUser(email, data)
		c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
	} else {
		c.JSON(http.StatusOK, r.H(e.ERROR_AUTH, nil))
	}
}
