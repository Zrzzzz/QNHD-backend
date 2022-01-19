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
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(content, "content")
	ok, verr := r.ErrorValid(&valid, "Get post replys")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 如果不是超管，看是否为部门对应管理
	if !models.IsUserSuperAdmin(uid) {
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
	// 添加回复
	id, err := models.AddPostReply(map[string]interface{}{
		"post_id": util.AsUint(postId),
		"from":    models.PostReplyType(1),
		"content": content,
	})
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 通知回复
	err = models.AddUnreadPostReply(util.AsUint(uid), id)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
