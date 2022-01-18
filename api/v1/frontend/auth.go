package frontend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
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

type authRes struct {
	ErrorCode int        `json:"error_code"`
	Message   string     `json:"message"`
	Result    authResult `json:"result"`
}

type authResult struct {
	UserNumber string `json:"userNumber"`
	Telephone  string `json:"telephone"`
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/token
func GetAuthToken(c *gin.Context) {
	token := c.Query("token")
	valid := validation.Validation{}
	valid.Required(token, "token")
	ok, verr := r.E(&valid, "get auth")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	var req *http.Request
	req, _ = http.NewRequest("GET", "https://api.twt.edu.cn/api/user/single", nil)
	req.Header.Add("token", token)

	var err error
	// 解析出用户的number
	client := &http.Client{}
	req.Header.Add("domain", setting.AppSetting.WPYDomain)
	ticket := b64.StdEncoding.EncodeToString([]byte(setting.AppSetting.WPYAppSecret + "." + setting.AppSetting.WPYAppKey))
	req.Header.Add("ticket", ticket)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	// var v map[string]interface{}
	var v authRes
	err = json.Unmarshal(body, &v)
	if err != nil {
		logging.Error("Auth error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if v.ErrorCode != 0 {
		r.Success(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": v.Message})
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
	ok, verr := r.E(&valid, "get auth")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 请求服务器
	var req *http.Request
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("account", user)
	writer.WriteField("password", password)

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err = http.NewRequest("POST", "https://api.twt.edu.cn/api/auth/common", payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 解析出用户的number
	client := &http.Client{}
	req.Header.Add("domain", setting.AppSetting.WPYDomain)
	ticket := b64.StdEncoding.EncodeToString([]byte(setting.AppSetting.WPYAppSecret + "." + setting.AppSetting.WPYAppKey))
	req.Header.Add("ticket", ticket)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	// var v map[string]interface{}
	var v authRes
	err = json.Unmarshal(body, &v)
	if err != nil {
		logging.Error("Auth error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	if v.ErrorCode != 0 {
		r.Success(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": v.Message})
		logging.Error("Auth er%v", v)
		return
	}

	auth(v.Result, c)
}

func auth(result authResult, c *gin.Context) {
	uid, err := models.ExistUser(result.UserNumber)
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("auth error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	// 如果不存在就创建一个用户
	if uid == 0 {
		uid, err = models.AddUser(result.UserNumber, result.Telephone, "")
	}

	if err != nil {
		logging.Error("auth error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	token, err := util.GenerateToken(fmt.Sprintf("%d", uid))
	if err != nil {
		logging.Error("auth error: %v", err)
		r.Success(c, e.ERROR_AUTH, map[string]interface{}{"error": err.Error()})
		return
	}
	data["token"] = token
	data["uid"] = uid
	r.Success(c, e.SUCCESS, data)
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
	ok, verr := r.E(&valid, "Refresh Token")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		logging.Error(err.Error())
		r.Success(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, map[string]interface{}{"error": err.Error()})
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
		data["uid"] = claims.Uid
	}
	r.Success(c, code, data)
}
