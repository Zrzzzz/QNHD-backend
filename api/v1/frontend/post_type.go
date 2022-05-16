package frontend

import (
	"qnhd/api/v1/common"

	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /f/posttypes
func GetPostTypes(c *gin.Context) {
	common.GetPostTypes(c)
}
