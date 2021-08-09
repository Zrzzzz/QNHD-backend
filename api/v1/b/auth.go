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
	Super bool   `json:"super" example:"false"`
}

// @Tags backend
// @Summary 后端获取token
// @Description getAuth
// @Accept json
// @Produce json
// @Param name query string true "admin name"
// @Param password query string true "admin password，发送密码的32位小写md5"
// @Param token query string false "jwt token，如果用此参数无需传递name和password，可用于刷新token"
// @Success 200 {object} models.Response{data=authRes}
// @Failure 20003 {object} models.Response "失败不返回数据"
// @Router /b/auth [get]
func GetAuth(c *gin.Context) {
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if token := c.Query("token"); token == "" {
		username := c.Query("name")
		password := c.Query("password")
		valid := validation.Validation{}
		a := auth{Username: username, Password: password}
		ok, _ := valid.Valid(&a)

		if ok {
			isExist := models.CheckAdmin(username, password)
			if isExist {
				token, err := util.GenerateToken(username, password)
				if err != nil {
					code = e.ERROR_AUTH_TOKEN
				} else {
					data["token"] = token
					// 传回一个超管字段判断
					data["super"] = username == setting.AdminName
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
	} else {
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, r.H(e.ERROR, nil))
		}
		if models.CheckAdmin(claims.Username, claims.Password) {
			token, err := util.GenerateToken(claims.Username, claims.Password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				// 传回一个超管字段判断
				data["super"] = claims.Username == setting.AdminName
				code = e.SUCCESS
			}
		}
	}
	c.JSON(http.StatusOK, r.H(code, data))
}
