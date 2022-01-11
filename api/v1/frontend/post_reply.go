package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param post_id
// @return
// @route /f/post/replys
func GetPostReplys(c *gin.Context) {
	postId := c.Query("post_id")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	ok, verr := r.E(&valid, "Get post replys")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	list, err := models.GetPostReplys(postId)
	if err != nil {
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.Success(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param post_id, content
// @return
// @route /f/post/reply
func AddPostReply(c *gin.Context) {
	uid := r.GetUid(c)
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
	// 校验有无权限回复
	post, err := models.GetPost(postId)
	if util.AsStrU(post.Uid) != uid {
		r.Success(c, e.ERROR_RIGHT, map[string]interface{}{"error": err.Error()})
		return
	}
	// 添加回复
	_, err = models.AddPostReply(map[string]interface{}{
		"post_id": util.AsUint(postId),
		"from":    models.PostReplyType(0),
		"content": content,
	})
	if err != nil {
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
