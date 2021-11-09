package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/setting"
	"qnhd/pkg/util"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Email    string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type authRes struct {
	Token string `json:"token"`
	Uid   int    `json:"uid"`
}

type authRes1 struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		UserNumber string `json:"userNumber"`
	} `json:"result"`
}

func GetAuth1(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		r.R(c, http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}
	var err error
	// 解析出用户的number
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.twt.edu.cn/api/user/single", nil)
	req.Header.Add("domain", setting.AppSetting.WPYDomain)
	ticket := b64.StdEncoding.EncodeToString([]byte(setting.AppSetting.WPYAppSecret + "." + setting.AppSetting.WPYAppKey))
	req.Header.Add("ticket", ticket)
	req.Header.Add("token", token)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	// var v map[string]interface{}
	var v authRes1
	err = json.Unmarshal(body, &v)
	if err != nil {
		logging.Error(" error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if v.ErrorCode != 0 {
		fmt.Println("hhhh")
		r.R(c, http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		logging.Error("wpy token verify failed.")
		return
	}

	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @Tags front, auth
// @Summary 前端获取token
// @Accept json
// @Produce json
// @Param email query string true "admin name"
// @Param password query string true "admin password，发送密码的32位小写md5"
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
		uid, err := models.CheckUser(email, password)
		if err != nil {
			logging.Error("Auth user error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		if uid > 0 {
			// tag = 1 means is USER
			token, err := util.GenerateToken(fmt.Sprintf("%d", uid), 1)
			if err != nil {
				code = e.ERROR_GENERATE_TOKEN
			} else {
				data["token"] = token
				data["uid"] = uid
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
	ok, verr := r.E(&valid, "Refresh Token Front")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		logging.Error(err.Error())
		r.R(c, http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": err.Error()})
		return
	}

	// 判断是否为用户
	if claims.Tag != util.USER {
		logging.Error("权限错误, not user")
		r.R(c, http.StatusOK, e.ERROR_AUTH, map[string]interface{}{"error": "权限错误, not user"})
		return
	}
	var code = e.SUCCESS
	var data = make(map[string]interface{})
	// tag = 1 means is USER
	token, err = util.GenerateToken(claims.Uid, 1)
	if err != nil {
		code = e.ERROR_GENERATE_TOKEN
	} else {
		data["token"] = token
		data["uid"] = claims.Uid
	}
	r.R(c, http.StatusOK, code, data)
}
