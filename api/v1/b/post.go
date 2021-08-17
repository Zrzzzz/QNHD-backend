package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/setting"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @Tags backend, post
// @Summary 获取帖子
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=models.Post}}
// @Router /b/post [get]
func GetPosts(c *gin.Context) {
	var pageSize = setting.AppSetting.PageSize
	content := c.Query("content")

	data := make(map[string]interface{})

	list := models.GetPosts(util.GetPage(c), pageSize, content)

	data["list"] = list
	data["total"] = len(list)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags backend, post
// @Summary 删除帖子
// @Accept json
// @Produce json
// @Param id query int true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /b/post [delete]
func DeletePosts(c *gin.Context) {
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok := r.E(&valid, "Delete notices")
	if !ok {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.DeletePostsAdmin(id)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
