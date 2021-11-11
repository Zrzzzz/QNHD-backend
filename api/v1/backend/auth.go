package backend

import (
	"fmt"
	"net/http"
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
// @route /b/auth
func GetAuth(c *gin.Context) {
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	number := c.Query("number")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(number, "number")
	valid.Required(password, "password")
	ok, verr := r.E(&valid, "Auth")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
	if user.Uid > 0 {
		// tag = 0 means ADMIN
		token, err := util.GenerateToken(fmt.Sprintf("%d", user.Uid))
		if err != nil {
			code = e.ERROR_GENERATE_TOKEN
		} else {
			data["token"] = token
			data["uid"] = user.Uid
			code = e.SUCCESS
		}
	} else {
		code = e.ERROR_AUTH
	}
	r.R(c, http.StatusOK, code, data)
}
