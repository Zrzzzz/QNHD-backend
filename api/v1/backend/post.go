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

// @method [get]
// @way [query]
// @param content page page_size
// @return postList
// @route /b/posts
func GetPosts(c *gin.Context) {
	postType := c.Query("type")
	content := c.Query("content")
	departmentId := c.Query("department_id")
	solved := c.Query("solved")
	tagId := c.Query("tag_id")

	valid := validation.Validation{}
	valid.Required(postType, "type")
	valid.Numeric(postType, "type")
	if solved != "" {
		valid.Numeric(solved, "solved")
	}
	if departmentId != "" {
		valid.Numeric(departmentId, "department_id")
	}
	if tagId != "" {
		valid.Numeric(tagId, "tag_id")
	}
	ok, verr := r.ErrorValid(&valid, "Get posts")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	postTypeint := util.AsInt(postType)
	valid.Range(postTypeint, 0, 2, "postType")
	if solved != "" {
		solvedint := util.AsInt(solved)
		valid.Range(solvedint, 0, 1, "solved")
	}
	ok, verr = r.ErrorValid(&valid, "Get posts")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	maps := map[string]interface{}{
		"type":          models.PostType(postTypeint),
		"content":       content,
		"solved":        solved,
		"department_id": departmentId,
		"tag_id":        tagId,
	}

	list, cnt, err := models.GetPosts(c, maps)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = cnt

	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param id
// @return post
// @route /b/post
func GetPost(c *gin.Context) {
	id := c.Query("id")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")

	ok, verr := r.ErrorValid(&valid, "Get Posts")
	if !ok {
		r.OK(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	pr, err := models.GetPostResponse(id)
	if err != nil {
		logging.Error("Get post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := map[string]interface{}{
		"post": pr,
	}
	r.OK(c, e.SUCCESS, data)
}

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
	if ok := models.IsDepartmentHasUser(util.AsUint(uid), post.DepartmentId); !ok {
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

// @method [delete]
// @way [query]
// @param id
// @return
// @route /b/post/delete
func DeletePosts(c *gin.Context) {
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
