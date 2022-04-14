package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /b/banner
func GetNewestBanner(c *gin.Context) {
	banner, err := models.GetNewestBanner()
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"content": banner.Content})
}

// @method [post]
// @way [formdata]
// @param content
// @return
// @route /b/banner
func AddNewBanner(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		r.Error(c, e.INVALID_PARAMS, "")
		return
	}
	if err := models.AddNewBanner(content); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
