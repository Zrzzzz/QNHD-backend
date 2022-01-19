package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param number, password
// @return token
// @route /b/auth/number
func GetAuthNumber(c *gin.Context) {
	code := e.INVALID_PARAMS
	number := c.Query("number")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(number, "number")
	valid.Required(password, "password")
	ok, verr := r.ErrorValid(&valid, "Auth")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	user, err := models.GetUser(map[string]interface{}{
		"number":   number,
		"password": password,
	})
	if err != nil {
		logging.Error("check admin error:%v", err)
		code = e.ERROR_DATABASE
	}
	auth(c, user, code)
}

// @method [get]
// @way [query]
// @param number, password
// @return token
// @route /b/auth/phone
func GetAuthPhone(c *gin.Context) {
	code := e.INVALID_PARAMS
	phone_number := c.Query("phone_number")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(phone_number, "phone_number")
	valid.Required(password, "password")
	ok, verr := r.ErrorValid(&valid, "Auth")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	user, err := models.GetUser(map[string]interface{}{
		"phone_number": phone_number,
		"password":     password,
	})
	if err != nil {
		logging.Error("check admin error:%v", err)
		code = e.ERROR_DATABASE
	}
	auth(c, user, code)
}

func auth(c *gin.Context, user models.User, code int) {
	data := make(map[string]interface{})
	if user.Uid > 0 {
		// tag = 0 means ADMIN
		token, err := util.GenerateToken(util.AsStrU(user.Uid))
		if err != nil {
			code = e.ERROR_GENERATE_TOKEN
		} else {
			data["token"] = token
			data["user"] = user
			code = e.SUCCESS
		}
	} else {
		code = e.ERROR_AUTH
	}
	r.OK(c, code, data)
}
