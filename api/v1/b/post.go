package b

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method get
// @way query
// @param content page page_size
// @return list
func GetPosts(c *gin.Context) {

	content := c.Query("content")

	data := make(map[string]interface{})
	base, size := util.HandlePaging(c)
	list, err := models.GetPosts(base, size, content)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	data["list"] = list
	data["total"] = len(list)
	r.R(c, http.StatusOK, e.SUCCESS, data)
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
	ok, verr := r.E(&valid, "Delete notices")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeletePostsAdmin(id)
	if err != nil {
		logging.Error("Delete posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
