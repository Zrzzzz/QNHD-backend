package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [post]
// @way [query]
// @param
// @return
// @route
// @method [delete]
// @way [query]
// @param id
// @return
// @route /b/post/delete
func DeletePosts(c *gin.Context) {
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.E(&valid, "Delete post")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeletePostsAdmin(id)
	if err != nil {
		logging.Error("Delete posts error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
