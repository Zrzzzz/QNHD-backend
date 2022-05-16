package frontend

import (
	"fmt"

	"qnhd/api/v1/common"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/token
func GetAuthToken(c *gin.Context) {
	token := c.Query("token")
	valid := validation.Validation{}
	valid.Required(token, "token")
	ok, verr := r.ErrorValid(&valid, "get auth")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	v, err := twtservice.GetAuthByToken(token)
	if err != nil {
		logging.Error("Auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if v.ErrorCode != 0 {
		r.OK(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": v.Message})
		logging.Error("Auth er%v", v)
		return
	}
	auth(v.Result, c)
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/passwd
func GetAuthPasswd(c *gin.Context) {
	user := c.Query("user")
	password := c.Query("password")
	valid := validation.Validation{}
	valid.Required(user, "user")
	valid.Required(password, "password")
	ok, verr := r.ErrorValid(&valid, "get auth")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	v, err := twtservice.GetAuthByPasswd(user, password)
	if err != nil {
		logging.Error("Auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if v.ErrorCode != 0 {
		r.OK(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": v.Message})
		logging.Error("Auth er%v", v)
		return
	}

	auth(v.Result, c)
}

// 认证过程
func auth(result twtservice.TwTAuthResult, c *gin.Context) {
	uid, err := models.ExistUser("", result.UserNumber)
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 如果不存在就创建一个用户
	if uid == 0 {
		uid, err = models.AddUser("", result.UserNumber, "", result.Telephone, true)
	}

	if err != nil {
		logging.Error("auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	user, _ := models.GetUser(map[string]interface{}{"id": uid})

	token, err := util.GenerateToken(fmt.Sprintf("%d", uid))
	if err != nil {
		logging.Error("auth error: %v", err)
		r.OK(c, e.ERROR_AUTH, map[string]interface{}{"error": err.Error()})
		return
	}
	data["token"] = token
	data["uid"] = uid
	data["user"] = user
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/:token
func RefreshToken(c *gin.Context) {
	common.RefreshToken(c)
}
