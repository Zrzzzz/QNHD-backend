package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/upload"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param content page page_size
// @return postList
// @route /f/posts
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
	ok, verr := r.E(&valid, "Get posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	postTypeint := util.AsInt(postType)
	valid.Range(postTypeint, 0, 2, "postType")
	if solved != "" {
		solvedint := util.AsInt(solved)
		valid.Range(solvedint, 0, 1, "solved")
	}
	ok, verr = r.E(&valid, "Get posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	uid := r.GetUid(c)
	maps := map[string]interface{}{
		"type":          models.PostType(postTypeint),
		"content":       content,
		"solved":        solved,
		"department_id": departmentId,
		"tag_id":        tagId,
	}

	list, err := models.GetPostResponses(c, uid, maps)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.Success(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/user
func GetUserPosts(c *gin.Context) {
	uid := r.GetUid(c)

	list, err := models.GetUserPostResponses(c, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.Success(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/fav
func GetFavPosts(c *gin.Context) {
	uid := r.GetUid(c)
	list, err := models.GetFavPostResponses(c, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)

	r.Success(c, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/history
func GetHistoryPosts(c *gin.Context) {
	uid := r.GetUid(c)

	list, err := models.GetHistoryPostResponses(c, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
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
// @param id
// @return post
// @route /f/post
func GetPost(c *gin.Context) {
	id := c.Query("id")
	uid := r.GetUid(c)
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")

	ok, verr := r.E(&valid, "Get Posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	pr, err := models.GetPostResponseAndVisit(id, uid)
	if err != nil {
		logging.Error("Get post error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := map[string]interface{}{
		"post": pr,
	}
	r.Success(c, e.SUCCESS, data)
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
	valid := validation.Validation{}
	valid.Required(content, "content")
	valid.Required(postType, "postType")
	valid.Numeric(postType, "postType")
	valid.Required(campus, "campus")
	valid.Numeric(campus, "campus")
	valid.Required(title, "title")
	valid.MaxSize(title, 30, "title")
	valid.MaxSize(content, 1000, "content")
	ok, verr := r.E(&valid, "Add posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	campusint := util.AsInt(campus)
	valid.Range(campusint, 0, 2, "campus")
	postTypeint := util.AsInt(postType)
	valid.Range(postTypeint, 0, 1, "postType")
	ok, verr = r.E(&valid, "Add posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 需要根据类型判断返回类型
	// 判断type
	if postTypeint == int(models.POST_SCHOOL) {
		// 可选tag
		if tagId != "" {
			valid.Numeric(tagId, "tag_id")
		}
	} else if postTypeint == int(models.POST_HOLE) {
		// 必须要求部门id不为0
		valid.Required(departId, "department_id")
		valid.Numeric(departId, "department_id")
	}

	// 处理图片
	form, err := c.MultipartForm()
	if err != nil {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": err.Error()})
		return
	}
	imgs := form.File["images"]
	if len(imgs) > 3 {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": "images count should less than 3."})
		return
	}
	imageUrls, err := upload.SaveImagesFromFromData(imgs, c)
	if err != nil {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": err.Error()})
		return
	}

	intuid := util.AsUint(uid)
	maps := map[string]interface{}{
		"uid":        intuid,
		"type":       models.PostType(postTypeint),
		"campus":     models.PostCampusType(campusint),
		"title":      title,
		"content":    content,
		"image_urls": imageUrls,
	}

	if postTypeint == 0 && tagId != "" {
		maps["tag_id"] = tagId
	} else if postTypeint == 1 {
		maps["department_id"] = util.AsUint(departId)
	}
	id, err := models.AddPost(maps)
	if err != nil {
		logging.Error("Add post error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	data["image_urls"] = imageUrls
	r.Success(c, e.SUCCESS, data)
}

// @method [put]
// @way [formdata]
// @param post_id, rating
// @return
// @route /f/post/solve
func EditPostSolved(c *gin.Context) {
	uid := r.GetUid(c)
	postId := c.Query("post_id")
	rating := c.Query("rating")
	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(rating, "rating")
	valid.Numeric(rating, "rating")
	ok, verr := r.E(&valid, "Delete posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 限制评分范围
	valid.Range(util.AsInt(rating), 1, 10, "rating")
	ok, verr = r.E(&valid, "Delete posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 判断是否为发帖人
	// 校验有无权限回复
	post, err := models.GetPost(postId)
	if util.AsStrU(post.Uid) != uid {
		r.Success(c, e.ERROR_RIGHT, map[string]interface{}{"error": err.Error()})
		return
	}
	err = models.EditPostSolved(postId, rating)
	if err != nil {
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
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
	ok, verr := r.E(&valid, "Delete posts")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeletePostsUser(postId, uid)
	if err != nil {
		logging.Error("Delete posts error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
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
	ok, verr := r.E(&valid, "fav or unfav post")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, map[string]interface{}{"count": cnt})
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
	ok, verr := r.E(&valid, "like or unlike post")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, map[string]interface{}{"count": cnt})
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
	ok, verr := r.E(&valid, "dis or undis post")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}
