package r

import (
	"qnhd/pkg/e"

	"github.com/gin-gonic/gin"
)

func H(code int, data map[string]interface{}) gin.H {
	return gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	}
}
