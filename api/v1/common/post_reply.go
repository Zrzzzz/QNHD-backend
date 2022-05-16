package common

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetPostReplys(c *gin.Context) {
	postId := c.Query("post_id")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	ok, verr := r.ErrorValid(&valid, "Get post replys")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	list, err := models.GetPostReplyResponses(postId)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}
