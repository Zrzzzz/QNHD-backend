package frontend

import (
	"fmt"

	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"qnhd/request/twtauth"

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
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	v, err := twtauth.GetAuthByToken(token)
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
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	v, err := twtauth.GetAuthByPasswd(user, password)
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
func auth(result twtauth.TwTAuthResult, c *gin.Context) {
	uid, err := models.ExistUser(result.UserNumber)
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 如果不存在就创建一个用户
	if uid == 0 {
		uid, err = models.AddUser(result.UserNumber, "", result.Telephone, true)
	}

	if err != nil {
		logging.Error("auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	token, err := util.GenerateToken(fmt.Sprintf("%d", uid))
	if err != nil {
		logging.Error("auth error: %v", err)
		r.OK(c, e.ERROR_AUTH, map[string]interface{}{"error": err.Error()})
		return
	}
	data["token"] = token
	data["uid"] = uid
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/:token | /b/auth/:token
func RefreshToken(c *gin.Context) {
	token := c.Param("token")
	valid := validation.Validation{}
	valid.Required(token, "token")
	ok, verr := r.ErrorValid(&valid, "Refresh Token")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		logging.Error(err.Error())
		r.OK(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": err.Error()})
		return
	}

	var code = e.SUCCESS
	var data = make(map[string]interface{})
	// tag = 1 means is USER
	token, err = util.GenerateToken(claims.Uid)
	if err != nil {
		code = e.ERROR_GENERATE_TOKEN
	} else {
		data["token"] = token
		data["uid"] = util.AsUint(claims.Uid)
	}
	r.OK(c, code, data)
}
