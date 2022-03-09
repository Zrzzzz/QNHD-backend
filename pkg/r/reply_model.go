package r

import (
	"fmt"
	"net/http"

	"qnhd/pkg/e"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func FindToken(c *gin.Context) string {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	return token
}

func GetUid(c *gin.Context) string {
	var claims *util.Claims
	token := FindToken(c)
	if token == "" {
		return ""
	} else {
		claims, _ = util.ParseToken(token)
		return claims.Uid
	}
}

// 返回是否没有错误
func ErrorValid(valid *validation.Validation, errorPhase string) (bool, error) {
	s := errorPhase
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			s += fmt.Sprintf("\n%v %v\n", r.Key, r.Message)
		}
	}
	return !valid.HasErrors(), fmt.Errorf(s)
}

func R(c *gin.Context, httpCode int, code int, data map[string]interface{}) {
	c.JSON(httpCode, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func OK(c *gin.Context, code int, data map[string]interface{}) {
	R(c, http.StatusOK, code, data)
}

func Error(c *gin.Context, code int, err string) {
	R(c, http.StatusOK, code, map[string]interface{}{"error": err})
}
