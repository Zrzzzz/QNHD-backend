package backend

import (
	"qnhd/api/v1/common"
	"qnhd/enums/PostEtagType"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"
	"qnhd/request/yunpian"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param content page page_size
// @return postList
// @route /b/posts
func GetPosts() gin.HandlerFunc {
	return common.GetPosts(false)
}

// @method [get]
// @way [query]
// @param id
// @return post
// @route /b/post
func GetPost() gin.HandlerFunc {
	return common.GetPost(false)
}

// @method [get]
// @way [query]
// @param
// @return
// @route /b/posts/undistributed
func GetUndistributedPosts(c *gin.Context) {
	data := make(map[string]interface{})
	list, err := models.GetUndistributedPosts(c)
	if err != nil {
		logging.Error("get undistributed posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param uid
// @return
// @route /b/posts/user
func GetUserPosts(c *gin.Context) {
	uid := c.Query("uid")
	t := c.Query("type")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(t, "type")
	valid.Numeric(t, "type")
	ok, verr := r.ErrorValid(&valid, "Get user history")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	data := make(map[string]interface{})
	list, err := models.GetUserPostResponses(c, uid, t == "1")
	if err != nil {
		logging.Error("get user history error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data["list"] = list
	r.OK(c, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param post_id, new_department_id
// @return
// @route /b/post/transfer/department
func TransferPostDepartment(c *gin.Context) {
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
	uid := r.GetUid(c)
	post, err := models.GetPost(postId)
	if err != nil {
		logging.Error("transfer department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 判断是否为校务帖子
	if post.Type != models.POST_SCHOOL_TYPE {
		r.Error(c, e.ERROR_POST_TYPE, "")
		return
	}
	// 判断是否是这个帖子的管理员
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
	// 记录历史
	if err := models.AddPostDepartmentTransferLog(util.AsUint(uid), post.Id, post.DepartmentId, util.AsUint(newDepartmentId)); err != nil {
		logging.Error(err.Error())
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param post_id, new_department_id
// @return
// @route /b/post/distribute
func DistributePost(c *gin.Context) {
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
	uid := r.GetUid(c)
	post, err := models.GetPost(postId)
	if err != nil {
		logging.Error("transfer department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 判断是否为校务帖子
	if post.Type != models.POST_SCHOOL_TYPE {
		r.Error(c, e.ERROR_POST_TYPE, "")
		return
	}
	err = models.DistributePost(postId, newDepartmentId)
	if err != nil {
		logging.Error("transfer department error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 向新的部门的管理员发通知
	if err := yunpian.NotifyNewPost(util.AsUint(newDepartmentId), post.Title); err != nil {
		logging.Error(err.Error())
	}
	// 记录历史
	if err := models.AddPostDepartmentTransferLog(util.AsUint(uid), post.Id, post.DepartmentId, util.AsUint(newDepartmentId)); err != nil {
		logging.Error(err.Error())
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param post_id, new_department_id
// @return
// @route /b/post/transfer/type
func TransferPostType(c *gin.Context) {
	postId := c.PostForm("post_id")
	newTypeId := c.PostForm("new_type_id")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(newTypeId, "new_type_id")
	valid.Numeric(newTypeId, "new_type_id")
	ok, verr := r.ErrorValid(&valid, "transfer post")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	post, err := models.GetPost(postId)
	if err != nil {
		logging.Error("transfer type error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	if post.Type == util.AsInt(newTypeId) {
		r.Error(c, e.ERROR_POST_TYPE, "")
		return
	}
	err = models.EditPostType(postId, newTypeId)
	if err != nil {
		logging.Error("transfer type error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	uid := r.GetUid(c)
	// 记录历史
	if err := models.AddPostTypeTransferLog(util.AsUint(uid), post.Id, post.Type, util.AsInt(newTypeId)); err != nil {
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
	err := models.EditPostValue(postId, util.AsInt(value))
	if err != nil {
		logging.Error("edit post value error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param
// @return
// @route /b/post/etag
func EditPostEtag(c *gin.Context) {
	postId := c.PostForm("post_id")
	value := c.PostForm("value")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(value, "value")
	valid.Numeric(value, "value")
	ok, verr := r.ErrorValid(&valid, "edit post etag")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	valid.Range(util.AsInt(value), 0, 2, "value")
	ok, verr = r.ErrorValid(&valid, "edit post etag")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.EditPostEtag(postId, PostEtagType.Enum(util.AsInt(value)))
	if err != nil {
		logging.Error("edit post etag error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param post_id, tag_id
// @return
// @route /b/post_tag
func AddPostTag(c *gin.Context) {
	postId := c.PostForm("post_id")
	tagId := c.PostForm("tag_id")
	valid := validation.Validation{}
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	valid.Required(tagId, "post_id")
	valid.Numeric(tagId, "post_id")
	ok, verr := r.ErrorValid(&valid, "Add post tag")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.AddPostWithTag(nil, util.AsUint(postId), util.AsUint(tagId))
	if err != nil {
		logging.Error("Add post tag error: %v", err)
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
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	_, err := models.DeletePostAdmin(uid, id)
	if err != nil {
		logging.Error("Delete post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param post_id
// @return
// @route /b/post/recover
func RecoverPost(c *gin.Context) {
	postId := c.PostForm("post_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	ok, verr := r.ErrorValid(&valid, "Recover post")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	err := models.RecoverPost(postId)
	if err != nil {
		logging.Error("Recover post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [get]
// @way [query]
// @param post_id
// @return
// @route /b/post_tag/delete
func DeletePostTag(c *gin.Context) {
	id := c.Query("post_id")
	valid := validation.Validation{}
	valid.Required(id, "post_id")
	valid.Numeric(id, "post_id")
	ok, verr := r.ErrorValid(&valid, "Delete post tag")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.DeleteTagInPost(nil, util.AsUint(id))
	if err != nil {
		logging.Error("Delete post tag error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [get]
// @way [query]
// @param post_id
// @return
// @route /b/post_tag/delete
func DeletePostImages(c *gin.Context) {
	id := c.Query("post_id")
	valid := validation.Validation{}
	valid.Required(id, "post_id")
	valid.Numeric(id, "post_id")
	ok, verr := r.ErrorValid(&valid, "Delete post images")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.DeletePostImages(id)
	if err != nil {
		logging.Error("Delete post images error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
