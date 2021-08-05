package v1

import (
	"log"
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	sender "qnhd/pkg/email"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	id := c.Query("uid")
	email := c.Query("emial")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if id != "" {
		maps["UID"] = id
		maps["email"] = email
	}

	code := e.SUCCESS

	data["list"] = models.GetUsers(maps)
	data["total"] = models.GetUsersAll(maps)

	c.JSON(http.StatusOK, H(code, data))
}

func AddUsers(c *gin.Context) {
	email := c.Query("email")
	// password := c.Query("password")
	code := c.Query("code")
	log.Printf("adding user, email: %v, code : %v", email, code)

	if email == "" {
		c.JSON(http.StatusOK, H(e.ERROR_EXIST_EMAIL, nil))
	}

	if code == "" {
		// 进行邮箱验证码发送
		sender.SendEmail(email,
			func() {
				c.JSON(http.StatusOK, H(e.SUCCESS, nil))
			},
			func() {
				c.JSON(http.StatusOK, H(e.ERROR_SEND_EMAIL, nil))
			})
	} else {
		if sender.Code[email] == code {
			// 验证码正确
		} else {
			// 验证码错误
		}
	}

}

func EditUsers(c *gin.Context) {

}

func DeleteUsers(c *gin.Context) {

}
