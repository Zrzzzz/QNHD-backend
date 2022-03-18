package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"qnhd/request/yunpian"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [put]
// @way [formdata]
// @param post_id, transfer_id
// @return
// @route /b/post/transfer
func TransferPost(c *gin.Context) {
	postId := c.PostForm("post_id")
	newDepartmentId := c.PostForm("new_department_id")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(newDepartmentId, "new_department_id")
	valid.Numeric(newDepartmentId, "new_department_id")
	ok, verr := r.ErrorValid(&valid, "transfer post")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	// 先判断是否是这个帖子的管理员
	uid := r.GetUid(c)
	post, err := models.GetPost(postId)
	if err != nil {
		logging.Error("transfer department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if !models.RequireRight(uid, models.UserRight{Super: true}) && !models.IsDepartmentHasUser(util.AsUint(uid), post.DepartmentId) {
		r.Error(c, e.ERROR_RIGHT, "")
		return
	}
	err = models.EditPostDepartment(postId, newDepartmentId)
	if err != nil {
		logging.Error("transfer department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 向新的部门的管理员发通知
	if err := yunpian.NotifyNewPost(util.AsUint(newDepartmentId), post.Title); err != nil {
		logging.Error(err.Error())
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param
// @return
// @route /b/post/value
func EditPostValue(c *gin.Context) {
	postId := c.PostForm("post_id")
	value := c.PostForm("value")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(value, "value")
	valid.Numeric(value, "value")
	ok, verr := r.ErrorValid(&valid, "edit post value")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	valid.Range(util.AsInt(value), 0, 30000, "value")
	ok, verr = r.ErrorValid(&valid, "edit post value")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.EditPostValue(postId, util.AsUint(value))
	if err != nil {
		logging.Error("edit post value error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param id
// @return
// @route /b/post/delete
func DeletePost(c *gin.Context) {
	uid := r.GetUid(c)
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	ok, verr := r.ErrorValid(&valid, "Delete post")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeletePostsAdmin(uid, id)
	if err != nil {
		logging.Error("Delete posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
