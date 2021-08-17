package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type authRes struct {
	Token string `json:"token"`
	Super bool   `json:"super"`
}

// @Tags backend, auth
// @Summary 后端获取token
// @Accept json
// @Produce json
// @Param name query string true "admin name"
// @Param password query string true "admin password，发送密码的32位小写md5"
// @Success 200 {object} models.Response{data=authRes}
// @Failure 20003 {object} models.Response "失败不返回数据"
// @Router /b/auth [get]
func GetAuth(c *gin.Context) {
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	username := c.Query("name")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if ok {
		isExist := models.CheckAdmin(username, password)
		if isExist {
			token, err := util.GenerateToken(username)
			if err != nil {
				code = e.ERROR_GENERATE_TOKEN
			} else {
				data["token"] = token
				// 传回一个超管字段判断
				data["super"] = username == setting.AppSetting.AdminName
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error("auth error: %v", err)
		}
		c.JSON(http.StatusOK, r.H(e.ERROR_AUTH, nil))
		return
	}
	c.JSON(http.StatusOK, r.H(code, data))
}

// @Tags backend, auth
// @Summary 后端刷新token
// @Accept json
// @Produce json
// @Param token path string true "用户token"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=authRes}
// @Failure 400 {object} models.Response "无token"
// @Failure 20001 {object} models.Response "token检查失败"
// @Failure 20003 {object} models.Response "token生成失败"
// @Router /b/auth/{token} [get]
func RefreshToken(c *gin.Context) {
	token := c.Query("token")
	valid := validation.Validation{}
	valid.Required(token, "token")
	ok := r.E(&valid, "Refresh Token Backend")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		logging.Error(err.Error())
		c.JSON(http.StatusOK, r.H(e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil))
		return
	}

	var code int = e.SUCCESS
	var data = make(map[string]interface{})
	if models.ExistAdmin(claims.Username) {
		token, err := util.GenerateToken(claims.Username)
		if err != nil {
			code = e.ERROR_GENERATE_TOKEN
		} else {
			data["token"] = token
		}
	}

	c.JSON(http.StatusOK, r.H(code, data))
}
