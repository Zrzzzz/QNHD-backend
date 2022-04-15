package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/request/yunpian"

	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

const POST_SCHOOL_TYPE = 1

// @method [get]
// @way [query]
// @param content page page_size
// @return postList
// @route /f/posts
func GetPosts(front bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		postType := c.Query("type")
		searchMode := c.Query("search_mode")
		content := c.Query("content")
		departmentId := c.Query("department_id")
		solved := c.Query("solved")
		tagId := c.Query("tag_id")
		valueMode := c.Query("value_mode")

		valid := validation.Validation{}
		valid.Required(postType, "type")
		valid.Numeric(postType, "type")
		valid.Required(searchMode, "search_mode")
		valid.Numeric(searchMode, "search_mode")
		if valueMode == "" {
			// 默认值
			valueMode = "0"
		}
		valid.Numeric(valueMode, "value_mode")
		ok, verr := r.ErrorValid(&valid, "Get posts")
		if !ok {
			r.Error(c, e.INVALID_PARAMS, verr.Error())
			return
		}
		valid.Numeric(solved, "solved")
		valid.Numeric(departmentId, "department_id")
		valid.Numeric(tagId, "tag_id")
		postTypeint := util.AsInt(postType)
		searchModeint := util.AsInt(searchMode)
		valid.Range(searchModeint, 0, 1, "search_mode")
		if solved != "" {
			solvedint := util.AsInt(solved)
			valid.Range(solvedint, 0, 2, "solved")
		}
		ok, verr = r.ErrorValid(&valid, "Get posts")
		if !ok {
			r.Error(c, e.INVALID_PARAMS, verr.Error())
			return
		}

		data := make(map[string]interface{})
		maps := map[string]interface{}{
			"type":          postTypeint,
			"search_mode":   models.PostSearchModeType(searchModeint),
			"content":       content,
			"solved":        solved,
			"department_id": departmentId,
			"tag_id":        tagId,
			"value_mode":    models.PostValueModeType(util.AsInt(valueMode)),
		}
		if front {
			uid := r.GetUid(c)
			list, err := models.GetPostResponsesWithUid(c, uid, maps)
			if err != nil {
				logging.Error("Get posts error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			data["list"] = list
			data["total"] = len(list)
		} else {
			list, cnt, err := models.GetPostResponses(c, maps)
			if err != nil {
				logging.Error("Get posts error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			data["list"] = list
			data["total"] = cnt
		}

		r.OK(c, e.SUCCESS, data)
	}
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/user
func GetUserPosts(c *gin.Context) {
	uid := r.GetUid(c)

	list, err := models.GetUserPostResponseUsers(c, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/fav
func GetFavPosts(c *gin.Context) {
	uid := r.GetUid(c)
	list, err := models.GetFavPostResponseUsers(c, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/history
func GetHistoryPosts(c *gin.Context) {
	uid := r.GetUid(c)

	list, err := models.GetHistoryPostResponseUsers(c, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param id
// @return post
// @route /f/post
func GetPost(front bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		valid := validation.Validation{}
		valid.Required(id, "id")
		valid.Numeric(id, "id")

		ok, verr := r.ErrorValid(&valid, "Get Posts")
		if !ok {
			r.Error(c, e.INVALID_PARAMS, verr.Error())
			return
		}
		data := make(map[string]interface{})
		if front {
			uid := r.GetUid(c)
			pr, err := models.GetPostResponseUserAndVisit(id, uid)
			if err != nil {
				logging.Error("Get post error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			data["post"] = pr
		} else {
			pr, err := models.GetPostResponse(id)
			if err != nil {
				logging.Error("Get post error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			data["post"] = pr
		}
		r.OK(c, e.SUCCESS, data)
	}
}

// @method [post]
// @way [formdata]
// @param uid, type, title, content, campus, department_id, images
// @return uploadres
// @route /f/post
func AddPost(c *gin.Context) {
	uid := r.GetUid(c)
	postType := c.PostForm("type")
	title := c.PostForm("title")
	content := c.PostForm("content")
	tagId := c.PostForm("tag_id")
	campus := c.PostForm("campus")
	departId := c.PostForm("department_id")
	imageURLs := c.PostFormArray("images")
	valid := validation.Validation{}
	valid.Required(postType, "postType")
	valid.Numeric(postType, "postType")
	valid.Required(campus, "campus")
	valid.Numeric(campus, "campus")
	valid.Required(title, "title")
	valid.MaxSize(title, 30, "title")
	valid.MaxSize(content, 1000, "content")
	valid.MaxSize(imageURLs, 3, "images")
	ok, verr := r.ErrorValid(&valid, "Add posts")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	campusint := util.AsInt(campus)
	valid.Range(campusint, 0, 2, "campus")
	postTypeint := util.AsInt(postType)
	ok, verr = r.ErrorValid(&valid, "Add posts")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	// 需要根据类型判断返回类型
	// 判断type
	if postTypeint == POST_SCHOOL_TYPE {
		// 必须要求部门id不为0
		valid.Required(departId, "department_id")
		valid.Numeric(departId, "department_id")
	} else if models.IsValidPostType(postTypeint) {
		// 可选tag
		if tagId != "" {
			valid.Numeric(tagId, "tag_id")
		}
	}
	ok, verr = r.ErrorValid(&valid, "Add posts")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	// 限制无文字时必须有图
	if content == "" && len(imageURLs) == 0 {
		r.Error(c, e.INVALID_PARAMS, "缺失图片或内容")
		return
	}
	intuid := util.AsUint(uid)
	maps := map[string]interface{}{
		"uid":        intuid,
		"type":       postTypeint,
		"campus":     models.PostCampusType(campusint),
		"title":      title,
		"content":    content,
		"image_urls": imageURLs,
	}
	if postTypeint == POST_SCHOOL_TYPE {
		maps["department_id"] = util.AsUint(departId)
	} else if tagId != "" {
		maps["tag_id"] = tagId
	}
	id, err := models.AddPost(maps)
	if err != nil {
		logging.Error("Add post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 如果是校务贴，需要对部门发出通知
	if postTypeint == POST_SCHOOL_TYPE {
		err = yunpian.NotifyNewPost(util.AsUint(departId), title)
		if err != nil {
			logging.Error(err.Error())
		}
	}
	data := make(map[string]interface{})
	data["id"] = id
	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param post_id
// @return
// @route /f/post/visit
func VisitPost(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	valid := validation.Validation{}
	valid.Numeric(postId, "post_id")
	ok, verr := r.ErrorValid(&valid, "add post visit")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	if err := models.AddVisitHistory(uid, postId); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [put]
// @way [formdata]
// @param post_id, rating
// @return
// @route /f/post/solve
func EditPostSolved(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	rating := c.PostForm("rating")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(rating, "rating")
	valid.Numeric(rating, "rating")
	ok, verr := r.ErrorValid(&valid, "Delete posts")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	// 限制评分范围
	valid.Range(util.AsInt(rating), 1, 10, "rating")
	ok, verr = r.ErrorValid(&valid, "Delete posts")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	// 判断是否为发帖人
	// 校验有无权限修改
	post, err := models.GetPost(postId)
	if util.AsStrU(post.Uid) != uid {
		r.OK(c, e.ERROR_RIGHT, map[string]interface{}{"error": err.Error()})
		return
	}
	err = models.EditPost(postId, map[string]interface{}{
		"solved": models.POST_SOLVED,
		"rating": rating,
	})
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param post_id
// @return nil
// @route /f/post
func DeletePost(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.Query("post_id")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	ok, verr := r.ErrorValid(&valid, "Delete posts")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	_, err := models.DeletePostsUser(postId, uid)
	if err != nil {
		logging.Error("Delete posts error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param post_id, op
// @return nil
// @route /f/post/fav
func FavOrUnfavPost(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.ErrorValid(&valid, "fav or unfav post")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	var err error
	var cnt uint64
	if op == "1" {
		cnt, err = models.FavPost(postId, uid)
	} else {
		cnt, err = models.UnfavPost(postId, uid)
	}
	if err != nil {
		logging.Error("fav or unfav post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method [post]
// @way [formdata]
// @param post_id, op
// @return nil
// @route /f/post/like
func LikeOrUnlikePost(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.ErrorValid(&valid, "like or unlike post")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	var err error
	var cnt uint64
	if op == "1" {
		cnt, err = models.LikePost(postId, uid)
	} else {
		cnt, err = models.UnLikePost(postId, uid)
	}
	if err != nil {
		logging.Error("like or unlike post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method [post]
// @way [formdata]
// @param post_id, op
// @return nil
// @route /f/post/dis
func DisOrUndisPost(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.PostForm("post_id")
	op := c.PostForm("op")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(op, "op")
	valid.Numeric(op, "op")
	ok, verr := r.ErrorValid(&valid, "dis or undis post")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	var err error
	var cnt uint64
	if op == "1" {
		cnt, err = models.DisPost(postId, uid)
	} else {
		cnt, err = models.UnDisPost(postId, uid)
	}
	if err != nil {
		logging.Error("dis or undis post error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}
