package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /f/notice
func GetNotices(c *gin.Context) {
	list, err := models.GetNotices()
	if err != nil {
		logging.Error("Get notices error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.Success(c, e.SUCCESS, data)
}
