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
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Result    struct {
		UserNumber string `json:"userNumber"`
	} `json:"result"`
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth
func GetAuth(c *gin.Context) {
	token := c.Query("token")
	user := c.Query("user")
	password := c.Query("password")
	var pass = false
	var req *http.Request
	if token != "" {
		pass = true
		req, _ = http.NewRequest("GET", "https://api.twt.edu.cn/api/user/single", nil)
		req.Header.Add("token", token)
	}
	if user != "" || password != "" {
		pass = true
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("account", user)
		_ = writer.WriteField("password", password)
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
	}

	if !pass {
		r.R(c, http.StatusUnauthorized, e.INVALID_PARAMS, nil)
		return
	}

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

	uid, err := models.ExistUser(v.Result.UserNumber)
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("auth error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	// 如果不存在就创建一个用户
	if uid == 0 {
		uid, err = models.AddUser(v.Result.UserNumber, "", "")
	}

	if err != nil {
		logging.Error("auth error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	token, err = util.GenerateToken(fmt.Sprintf("%d", uid))
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
