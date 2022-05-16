package common

import (
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {
	token := c.Param("token")
	valid := validation.Validation{}
	valid.Required(token, "token")
	ok, verr := r.ErrorValid(&valid, "Refresh Token")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
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
