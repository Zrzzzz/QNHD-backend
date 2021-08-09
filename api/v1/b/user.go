package b

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"strings"

	"qnhd/api/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	uid := c.Query("uid")
	email := c.Query("email")

	valid := validation.Validation{}
	valid.Numeric(uid, "id")
	if email != "" {
		valid.Email(email, "email")
	}
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("Get user error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
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

	list := models.GetUsers(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(code, data))
}

func AddUsers(c *gin.Context) {
	email := c.Query("email")
	password := c.Query("password")

	if !checkMail(email) {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	if models.ExistUserByEmail(email) {
		c.JSON(http.StatusOK, r.H(e.ERROR_EXIST_EMAIL, nil))
		return
	}

	models.AddUser(email, password)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

func EditUsers(c *gin.Context) {
	email := c.Query("email")
	newPass := c.Query("new_password")

	if models.ExistUserByEmail(email) {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	data := make(map[string]interface{})
	data["password"] = newPass
	models.EditUser(email, data)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

func DeleteUsers(c *gin.Context) {
	email := c.Query("email")
	if !models.ExistUserByEmail(email) {
		models.DeleteUser(email)
		c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
	} else {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
	}
}

func checkMail(email string) bool {
	s := strings.Split(email, "@")
	if len(s) != 2 || len(s[0]) != 10 || s[1] != "tju.edu.cn" {
		return false
	}
	return true
}
