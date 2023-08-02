package common

import (
	"fmt"
	"math/rand"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

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
	uid, err := models.ExistUser("", strings.Trim(result.UserNumber, " "))
	data := make(map[string]interface{})
	if err != nil {
		logging.Error("auth error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 如果不存在就创建一个用户
	if uid == 0 {
		uid, err = models.AddUser(genNickname(), result.UserNumber, "", result.Telephone, result.Realname, true)
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

func genNickname() string {
	letters := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret := ""
	for i := 0; i < 8; i++ {
		ret += letters[r.Intn(len(letters))]
	}
	return ret
}
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
