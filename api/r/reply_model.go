package r

import (
	"qnhd/models"
	"qnhd/pkg/e"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
)

func H(code int, data map[string]interface{}) gin.H {
	return structs.Map(models.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: data,
	})
}
