package frontend

import (
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/upload"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type postResponse struct {
	models.Post
	Tag          models.Tag        `json:"tag"`
	Floors       []models.Floor    `json:"floors"`
	CommentCount int               `json:"comment_count"`
	IsLike       bool              `json:"is_like"`
	IsDis        bool              `json:"is_dis"`
	IsFav        bool              `json:"is_fav"`
	Pictures     []string          `json:"pictures"`
	Department   models.Department `json:"department"`
}

func makePostResponse(p models.Post, uid string) (postResponse, error) {
	var pr postResponse
	tag, err := models.GetTagInPost(util.AsStrU(p.Id))
	if err != nil {
		return pr, err
	}
	floors, err := models.GetFloorInPostShort(util.AsStrU(p.Id))
	if err != nil {
		return pr, err
	}
	pics, err := models.GetImageInPost(util.AsStrU(p.Id))
	if err != nil {
		return pr, err
	}
	var depart models.Department
	if p.DepartmentId > 0 {
		d, err := models.GetDepartment(p.DepartmentId)
		if err != nil {
			return pr, err
		}
		depart = d
	}
	return postResponse{
		Post:         p,
		Tag:          tag,
		Floors:       floors,
		CommentCount: len(floors),
		IsLike:       models.IsLikePostByUid(uid, util.AsStrU(p.Id)),
		IsDis:        models.IsDisPostByUid(uid, util.AsStrU(p.Id)),
		IsFav:        models.IsFavPostByUid(uid, util.AsStrU(p.Id)),
		Pictures:     pics,
		Department:   depart,
	}, nil
}

// @method [get]
// @way [query]
// @param content page page_size
// @return postList
// @route /f/posts
func GetPosts(c *gin.Context) {
	postType := c.Query("type")
	content := c.Query("content")
	departmentId := c.Query("department_id")
	valid := validation.Validation{}
	valid.Required(postType, "type")
	valid.Numeric(postType, "type")
	if departmentId != "" {
		valid.Numeric(departmentId, "department_id")
	}
	ok, verr := r.E(&valid, "Get posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	postTypeint := util.AsInt(postType)
	valid.Range(postTypeint, 0, 2, "postType")
	ok, verr = r.E(&valid, "Get posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	uid := r.GetUid(c)
	base, size := util.HandlePaging(c)
	maps := map[string]interface{}{
		"type":          postTypeint,
		"content":       content,
		"department_id": departmentId,
	}
	list, err := models.GetPosts(base, size, maps)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []postResponse{}
	for _, p := range list {
		pr, err := makePostResponse(p, uid)
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		retList = append(retList, pr)
	}

	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/user
func GetUserPosts(c *gin.Context) {
	uid := r.GetUid(c)
	base, size := util.HandlePaging(c)
	list, err := models.GetUserPosts(base, size, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []postResponse{}
	for _, p := range list {
		pr, err := makePostResponse(p, uid)
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		retList = append(retList, pr)
	}

	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/fav
func GetFavPosts(c *gin.Context) {
	uid := r.GetUid(c)
	base, size := util.HandlePaging(c)
	list, err := models.GetFavPosts(base, size, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []postResponse{}
	for _, p := range list {
		pr, err := makePostResponse(p, uid)
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		retList = append(retList, pr)
	}

	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [get]
// @way [query]
// @param page page_size
// @return postList
// @route /f/posts/history
func GetHistoryPosts(c *gin.Context) {
	uid := r.GetUid(c)
	base, size := util.HandlePaging(c)
	list, err := models.GetHistoryPosts(base, size, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []postResponse{}
	for _, p := range list {
		pr, err := makePostResponse(p, uid)
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		retList = append(retList, pr)
	}

	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, e.SUCCESS, data)
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
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	post, err := models.GetPost(id, uid)
	if err != nil {
		logging.Error("Get post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	pr, err := makePostResponse(post, uid)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := map[string]interface{}{
		"post": pr,
	}
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param uid content picture tag_id
// @return uploadres
// @route /f/post
func AddPost(c *gin.Context) {
	uid := r.GetUid(c)
	postType := c.PostForm("type")
	content := c.PostForm("content")
	tagId := c.PostForm("tag_id")
	campus := c.PostForm("campus")
	departId := c.PostForm("department_id")
	valid := validation.Validation{}
	valid.Required(content, "content")
	valid.Required(postType, "postType")
	valid.Required(campus, "campus")
	valid.Numeric(campus, "campus")
	valid.Numeric(postType, "postType")
	ok, verr := r.E(&valid, "Add posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	campusint := util.AsInt(campus)
	valid.Range(campusint, 0, 2, "campus")
	postTypeint := util.AsInt(postType)
	valid.Range(postTypeint, 0, 1, "postType")
	ok, verr = r.E(&valid, "Add posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	// 需要根据类型判断返回类型
	// 0为树洞帖子
	// 1为校务帖子

	// 判断type
	if postTypeint == 0 {
		// 可选tag
		if tagId != "" {
			valid.Numeric(tagId, "tag_id")
		}
	} else if postTypeint == 1 {
		// 必须要求部门id不为0
		valid.Required(departId, "department_id")
		valid.Numeric(departId, "department_id")
	}

	// 处理图片
	form, err := c.MultipartForm()
	if err != nil {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": err.Error()})
		return
	}
	pics := form.File["pictures"]
	if len(pics) > 3 {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": "pictures count should less than 3."})
		return
	}
	imageUrls, err := upload.SaveImagesFromFromData(pics, c)
	if err != nil {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": err.Error()})
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	intuid := util.AsUint(uid)
	maps["uid"] = intuid
	maps["type"] = postTypeint
	maps["campus"] = campusint
	maps["content"] = content
	maps["picture_urls"] = imageUrls

	if postTypeint == 0 {
		maps["tag_id"] = tagId
	} else if postTypeint == 1 {
		maps["department_id"] = util.AsUint(departId)
	}
	id, err := models.AddPost(maps)
	if err != nil {
		logging.Error("Add post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["id"] = id
	data["pictrue_url"] = imageUrls
	r.R(c, http.StatusOK, e.SUCCESS, data)
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
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeletePostsUser(postId, uid)
	if err != nil {
		logging.Error("Delete posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param post_id, op
// @return nil
// @route /f/favOrUnfav
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
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method [post]
// @way [formdata]
// @param post_id, op
// @return nil
// @route /f/likeOrUnlike
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
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method [post]
// @way [formdata]
// @param post_id, op
// @return nil
// @route /f/disOrUndis
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
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
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
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, map[string]interface{}{"count": cnt})
}
