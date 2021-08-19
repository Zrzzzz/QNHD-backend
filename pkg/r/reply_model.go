package r

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"

	"github.com/astaxie/beego/validation"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

// 通过code和data生成一个gin.H
func H(code int, data map[string]interface{}) gin.H {
	return structs.Map(models.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: data,
	})
}

// 返回是否没有错误
func E(valid *validation.Validation, errorPhase string) bool {
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("%v error: %v", errorPhase, r)
		}
	}
	return !valid.HasErrors()
}

func R(c *gin.Context, httpCode int, code int, data map[string]interface{}) {
	c.JSON(httpCode, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
