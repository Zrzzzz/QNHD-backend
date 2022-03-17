package backend

import (
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
	nickname := c.Query("user")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(nickname, "nickname")
	valid.Required(password, "password")
	ok, verr := r.ErrorValid(&valid, "Auth")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	user, err := chainAuth(
		map[string]interface{}{
			"nickname": nickname,
			"password": password,
		},
		map[string]interface{}{
			"number":   nickname,
			"password": password,
		},
		map[string]interface{}{
			"phone_number": nickname,
			"password":     password,
		},
	)
	if err != nil {
		logging.Error("check admin error:%v", err)
	}
	auth(c, user)
}

func chainAuth(maps ...map[string]interface{}) (models.User, error) {
	var (
		u models.User
		e error
	)
	for _, m := range maps {
		u, e = models.GetUser(m)
		if e == nil {
			return u, nil
		}
	}
	return u, e
}

func auth(c *gin.Context, user models.User) {
	var code int
	data := make(map[string]interface{})
	if user.Uid > 0 {
		// tag = 0 means ADMIN
		token, err := util.GenerateToken(util.AsStrU(user.Uid))
		if err != nil {
			code = e.ERROR_GENERATE_TOKEN
		} else {
			data["token"] = token
			data["user"] = user
			code = e.SUCCESS
		}
	} else {
		code = e.ERROR_AUTH
	}
	r.OK(c, code, data)
}
