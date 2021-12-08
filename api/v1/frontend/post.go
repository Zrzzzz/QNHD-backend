package frontend

import (
	"fmt"
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
	CommentCount int            `json:"comment_count"`
	IsLike       bool           `json:"is_like"`
	IsDis        bool           `json:"is_dis"`
	IsFav        bool           `json:"is_fav"`
	Tag          models.Tag     `json:"tag"`
	Floors       []models.Floor `json:"floors"`
	Pictures     []string       `json:"pictures"`
}

// @method [get]
// @way [query]
// @param content page page_size
// @return postList
// @route /f/posts
func GetPosts(c *gin.Context) {
	content := c.Query("content")
	uid := r.GetUid(c)
	base, size := util.HandlePaging(c)
	list, err := models.GetPosts(base, size, content)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []postResponse{}
	for _, p := range list {
		tag, err := models.GetTagInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		floors, err := models.GetFloorInPostShort(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		pics, err := models.GetImageInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		retList = append(retList, postResponse{
			Post:         p,
			Tag:          tag,
			Floors:       floors,
			CommentCount: len(floors),
			IsLike:       models.IsLikePostByUid(uid),
			IsDis:        models.IsDisPostByUid(uid),
			IsFav:        models.IsFavPostByUid(uid),
			Pictures:     pics,
		})
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
		tag, err := models.GetTagInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		floors, err := models.GetFloorInPostShort(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		pics, err := models.GetImageInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		retList = append(retList, postResponse{
			Post:         p,
			Tag:          tag,
			Floors:       floors,
			CommentCount: len(floors),
			IsLike:       models.IsLikePostByUid(uid),
			IsDis:        models.IsDisPostByUid(uid),
			IsFav:        models.IsFavPostByUid(uid),
			Pictures:     pics,
		})
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
		tag, err := models.GetTagInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		floors, err := models.GetFloorInPostShort(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		pics, err := models.GetImageInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		retList = append(retList, postResponse{
			Post:         p,
			Tag:          tag,
			Floors:       floors,
			CommentCount: len(floors),
			IsLike:       models.IsLikePostByUid(uid),
			IsDis:        models.IsDisPostByUid(uid),
			IsFav:        models.IsFavPostByUid(uid),
			Pictures:     pics,
		})
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
		tag, err := models.GetTagInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}
		floors, err := models.GetFloorInPostShort(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		pics, err := models.GetImageInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
			return
		}

		retList = append(retList, postResponse{
			Post:         p,
			Tag:          tag,
			Floors:       floors,
			CommentCount: len(floors),
			IsLike:       models.IsLikePostByUid(uid),
			IsDis:        models.IsDisPostByUid(uid),
			IsFav:        models.IsFavPostByUid(uid),
			Pictures:     pics,
		})
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
	tag, err := models.GetTagInPost(fmt.Sprintf("%d", post.Id))
	if err != nil {
		logging.Error("Get post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	floors, err := models.GetFloorInPostShort(fmt.Sprintf("%d", post.Id))
	if err != nil {
		logging.Error("Get post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	pics, err := models.GetImageInPost(fmt.Sprintf("%d", post.Id))
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data := map[string]interface{}{
		"post": postResponse{
			Post:         post,
			Tag:          tag,
			Floors:       floors,
			CommentCount: len(floors),
			IsLike:       models.IsLikePostByUid(uid),
			IsDis:        models.IsDisPostByUid(uid),
			IsFav:        models.IsFavPostByUid(uid),
			Pictures:     pics,
		},
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
	content := c.PostForm("content")
	tag_id := c.PostForm("tag_id")
	valid := validation.Validation{}
	valid.Required(content, "content")
	if tag_id != "" {
		valid.Numeric(tag_id, "tag_id")
	}
	ok, verr := r.E(&valid, "Add posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
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
	maps["content"] = content
	maps["picture_url"] = imageUrls
	maps["tag_id"] = tag_id
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
	if op == "1" {
		err = models.FavPost(postId, uid)
	} else {
		err = models.UnfavPost(postId, uid)
	}
	if err != nil {
		logging.Error("fav or unfav post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
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
	if op == "1" {
		err = models.LikePost(postId, uid)
	} else {
		err = models.UnLikePost(postId, uid)
	}
	if err != nil {
		logging.Error("like or unlike post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
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
	if op == "1" {
		err = models.DisPost(postId, uid)
	} else {
		err = models.UnDisPost(postId, uid)
	}
	if err != nil {
		logging.Error("dis or undis post error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.R(c, http.StatusOK, e.SUCCESS, nil)
}
