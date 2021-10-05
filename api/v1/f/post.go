package f

import (
	"fmt"
	"log"
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/upload"
	"qnhd/pkg/util"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type postRes struct {
	post postResponse
}

type postResponse struct {
	models.Post
	Tags         []models.Tag   `json:"tags"`
	Floors       []models.Floor `json:"floors"`
	StarCount    int            `json:"star_count"`
	CommentCount int            `json:"comment_count"`
}

type uploadRes struct {
	Id         int    `json:"id"`
	PictureUrl string `json:"picture_url"`
}

// @Tags front, post
// @Summary 获取多个简短post
// @Accept json
// @Produce json
// @Param content query string false "帖子内容"
// @Param page query string false "页数, 从0开始 默认为0"
// @Param page_size query int false "页面大小，默认为10"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=postResponse}}
// @Failure 400 {object} models.Response "无效参数"
// @Router /f/posts [get]
func GetPosts(c *gin.Context) {
	content := c.Query("content")
	base, size := util.HandlePaging(c)
	list, err := models.GetPosts(base, size, content)
	if err != nil {
		logging.Error("Get posts error: %v", err)
		r.R(c, http.StatusOK, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	retList := []postResponse{}
	for _, p := range list {
		tags, err := models.GetTagsInPost(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
			return
		}
		floors, err := models.GetFloorInPostShort(fmt.Sprintf("%d", p.Id))
		if err != nil {
			logging.Error("Get posts error: %v", err)
			r.R(c, http.StatusOK, e.ERROR_DATABASE, nil)
			return
		}
		retList = append(retList, postResponse{
			Post:         p,
			Tags:         tags,
			Floors:       floors,
			CommentCount: len(floors),
		})
	}

	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags front, post
// @Summary 获取单个post
// @Accept json
// @Produce json
// @Param id query int true "帖子id"
// @Param uid query int true "uid"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=postRes}
// @Failure 400 {object} models.Response "无效参数"
// @Router /f/post [get]
func GetPost(c *gin.Context) {
	id := c.Query("id")
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")

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
	tags, err := models.GetTagsInPost(fmt.Sprintf("%d", post.Id))
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
	data := map[string]interface{}{
		"post": postResponse{
			Post:   post,
			Tags:   tags,
			Floors: floors,
		},
	}
	r.R(c, http.StatusOK, e.SUCCESS, data)
}

// @Tags front, post
// @Summary 添加帖子
// @Accept json
// @Produce json
// @Param uid formData string true "发帖人id"
// @Param content formData string true "帖子内容"
// @Param picture formData string false "图片data，最大5MB"
// @Param tags formData []int false "标签id数组"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=uploadRes}
// @Failure 400 {object} models.Response "无效参数"
// @Router /f/post [post]
func AddPosts(c *gin.Context) {
	uid := c.PostForm("uid")
	content := c.PostForm("content")
	f, image, err := c.Request.FormFile("picture")
	hasImage := err == nil
	imageUrl := ""
	tags := c.PostFormArray("tags")

	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(content, "content")
	ok, verr := r.E(&valid, "Add posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	if hasImage {
		src, err := upload.CheckImage(&f, image)
		if err != nil {
			logging.Error("Add post error: %v", err)
			c.JSON(http.StatusOK, r.H(e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil))
			return
		}
		if err := c.SaveUploadedFile(image, src); err != nil {
			logging.Error("Add post error: %v", err)
			c.JSON(http.StatusOK, r.H(e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil))
			return
		}
		imageName := upload.GetImageName(image.Filename)
		imageUrl = upload.GetImagePath() + imageName
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	intuid, _ := strconv.ParseUint(uid, 10, 64)
	maps["uid"] = intuid
	maps["content"] = content
	maps["picture_url"] = imageUrl
	maps["tags"] = tags
	log.Println(tags)
	id, err := models.AddPosts(maps)
	data["id"] = id
	data["pictrue_url"] = imageUrl
	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @method get
// @way query
// @param content page page_size
// @return list
func FavOrUnFavPost(c *gin.Context) {
	uid := c.PostForm("uid")
	postId := c.PostForm("post_id")
	fav := c.PostForm("fav")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	valid.Required(fav, "fav")
	valid.Numeric(fav, "fav")
	ok, verr := r.E(&valid, "fav or unfav post")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	// 代表点赞问题
	var err error
	if fav == "1" {
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

// @method [get]
// @way [query]
// @param
// @return

// @Tags front, post
// @Summary 删除帖子
// @Accept json
// @Produce json
// @Param uid query string true "发帖人id"
// @Param post_id query string true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "无效参数"
// @Router /f/post [delete]
func DeletePosts(c *gin.Context) {
	uid := c.Query("uid")
	postId := c.Query("post_id")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
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
