package frontend

import (
	"qnhd/api/v1/common"

	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/token
func GetAuthToken(c *gin.Context) {
	common.GetAuthToken(c)
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/passwd
func GetAuthPasswd(c *gin.Context) {
	common.GetAuthPasswd(c)
}

// @method [get]
// @way [query]
// @param token
// @return token
// @route /f/auth/:token
func RefreshToken(c *gin.Context) {
	common.RefreshToken(c)
}
