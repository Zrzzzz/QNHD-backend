package frontend

import (
	"qnhd/api/v1/common"
	"qnhd/enums/PostReplyType"
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
	common.GetPostReplys(c)
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
	imageURLs := c.PostFormArray("images")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.MaxSize(imageURLs, 1, "images")
	ok, verr := r.ErrorValid(&valid, "Get post replys")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	// 校验有无权限回复
	post, err := models.GetPost(postId)
	if util.AsStrU(post.Uid) != uid {
		r.OK(c, e.ERROR_RIGHT, map[string]interface{}{"error": err.Error()})
		return
	}

	// 限制无文字时必须有图
	if content == "" && len(imageURLs) == 0 {
		r.Error(c, e.INVALID_PARAMS, "缺失图片或内容")
		return
	}
	// 添加回复
	_, err = models.AddPostReply(map[string]interface{}{
		"post_id": util.AsUint(postId),
		"sender":  PostReplyType.USER,
		"content": content,
		"urls":    imageURLs,
	})
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
