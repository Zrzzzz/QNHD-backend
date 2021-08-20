package f

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Email    string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type authRes struct {
	Token string `json:"token"`
}

// @Tags front, auth
// @Summary 前端获取token
// @Accept json
// @Produce json
// @Param name query string true "admin name"
// @Param password query string true "admin password，发送密码的32位小写md5"
// @Param token query string false "jwt token，如果用此参数无需传递name和password，可用于刷新token"
// @Success 200 {object} models.Response{data=authRes}
// @Failure 20003 {object} models.Response "失败不返回数据"
// @Router /f/auth [get]
func GetAuth(c *gin.Context) {
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	email := c.Query("email")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Email: email, Password: password}
	ok, _ := valid.Valid(&a)

	if ok {
		isExist, err := models.CheckUser(email, password)
		if err != nil {
			logging.Error("Auth user error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
			return
		}
		if isExist {
			token, err := util.GenerateToken(email)
			if err != nil {
				code = e.ERROR_GENERATE_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error("auth error: %v", err)
		}
		code = e.ERROR_AUTH
	}
	r.R(c, http.StatusOK, code, data)
}

// @Tags front, auth
// @Summary 前端刷新token
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=authRes}
// @Failure 400 {object} models.Response "无token"
// @Failure 20001 {object} models.Response "token检查失败"
// @Failure 20003 {object} models.Response "token生成失败"
// @Router /f/auth/{token} [get]
func RefreshToken(c *gin.Context) {
	token := c.Param("token")
	valid := validation.Validation{}
	valid.Required(token, "token")
	ok := r.E(&valid, "Refresh Token Front")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		logging.Error(err.Error())
		r.R(c, http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	var code int = e.SUCCESS
	var data = make(map[string]interface{})
	exist, err := models.ExistUser(claims.Username)
	if err != nil {
		logging.Error("Refresh token error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
		return
	}
	if exist {
		token, err := util.GenerateToken(claims.Username)
		if err != nil {
			code = e.ERROR_GENERATE_TOKEN
		} else {
			data["token"] = token
		}
	}
	r.R(c, http.StatusOK, code, data)
}
