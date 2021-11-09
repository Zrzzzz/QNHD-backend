package r

import (
	"fmt"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

func GetUid(c *gin.Context) string {
	var claims *util.Claims
	token := c.GetHeader("token")
	if token == "" {
		return ""
	} else {
		claims, _ = util.ParseToken(token)
		return claims.Uid
	}
}

// 通过code和data生成一个gin.H
func H(code int, data map[string]interface{}) gin.H {
	return structs.Map(models.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: data,
	})
}

// 返回是否没有错误
func E(valid *validation.Validation, errorPhase string) (bool, error) {
	s := errorPhase
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("%v error: %v", errorPhase, r)
			s += r.Error()
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
