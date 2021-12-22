package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [post]
// @way [formdata]
// @param post_id, content
// @return
// @route /b/post/reply
func AddPostReply(c *gin.Context) {
	postId := c.PostForm("post_id")
	content := c.PostForm("content")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(content, "content")
	ok, verr := r.E(&valid, "Get post replys")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 添加回复
	err := models.AddPostReply(map[string]interface{}{
		"post_id": util.AsUint(postId),
		"from":    models.PostReplyType(1),
		"content": content,
	})
	if err != nil {
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
