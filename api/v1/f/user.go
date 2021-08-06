package f

import (
	"fmt"
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
	sender "qnhd/pkg/email"
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

func AddUsers(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")
	code := c.Query("code")
	rightCode := fmt.Sprintf("%06d", sender.Code[email]%1000000)

	if email == "" || checkMail(email) {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
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

func DeleteUsers(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")

	if models.ValidUser(email, password) {
		models.DeleteUser(email)
		c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
	} else {
		c.JSON(http.StatusOK, r.H(e.ERROR_AUTH, nil))
	}
}
