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
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	content := c.PostForm("content")
	imageURLs := c.PostFormArray("images")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.MaxSize(content, 1000, "content")
	valid.MaxSize(imageURLs, 3, "images")
	ok, verr := r.ErrorValid(&valid, "Get post replys")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 如果不是超管，看是否为部门对应管理
	if !models.RequireRight(uid, models.UserRight{Super: true}) {
		depart, err := models.GetDepartmentByPostId(util.AsUint(postId))
		if err != nil {
			r.Error(c, e.ERROR_DATABASE, err.Error())
			return
		}
		if !models.IsDepartmentHasUser(util.AsUint(uid), depart.Id) {
			r.Error(c, e.ERROR_RIGHT, "")
			return
		}
	}

	// 限制无文字时必须有图
	if content == "" && len(imageURLs) == 0 {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": "缺失图片或内容"})
		return
	}
	// 添加回复
	id, err := models.AddPostReply(map[string]interface{}{
		"post_id": util.AsUint(postId),
		"sender":  models.PostReplyFromSchool,
		"content": content,
		"urls":    imageURLs,
	})
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if err := models.EditPost(postId, map[string]interface{}{"solved": models.POST_REPLIED}); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 通知回复
	err = models.AddUnreadPostReply(util.AsUint(postId), id)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
