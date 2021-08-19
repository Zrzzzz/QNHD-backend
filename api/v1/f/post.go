package f

import (
	"fmt"
	"log"
	"net/http"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/setting"
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
	Tags   []models.Tag
	Floors []models.Floor
}

// @Tags front, post
// @Summary 获取多个简短post
// @Accept json
// @Produce json
// @Param content query string false "帖子内容"
// @Param page query string false "页数, 从0开始 默认为0"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=models.ListRes{list=postResponse}}
// @Failure 400 {object} models.Response "无效参数"
// @Router /f/posts [get]
func GetPosts(c *gin.Context) {
	var pageSize = setting.AppSetting.PageSize
	content := c.Query("content")
	list := models.GetPosts(util.GetPage(c), pageSize, content)
	retList := []postResponse{}
	for _, p := range list {
		tags := models.GetTagsInPost(fmt.Sprintf("%d", p.Id))
		floors := models.GetFloorInPostShort(fmt.Sprintf("%d", p.Id))
		retList = append(retList, postResponse{
			Post:   p,
			Tags:   tags,
			Floors: floors,
		})
	}
	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

// @Tags front, post
// @Summary 获取单个post
// @Accept json
// @Produce json
// @Param id query string true "帖子id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response{data=postRes}
// @Failure 400 {object} models.Response "无效参数"
// @Router /f/post [get]
func GetPost(c *gin.Context) {
	id := c.Query("id")

	valid := validation.Validation{}
	valid.Required(id, "id")
	valid.Numeric(id, "id")

	ok := r.E(&valid, "Get Posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	post := models.GetPost(id)
	tags := models.GetTagsInPost(fmt.Sprintf("%d", post.Id))
	floors := models.GetFloorInPostShort(fmt.Sprintf("%d", post.Id))
	data := map[string]interface{}{
		"post": postResponse{
			Post:   post,
			Tags:   tags,
			Floors: floors,
		},
	}

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
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
// @Success 200 {object} models.Response
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
	ok := r.E(&valid, "Add posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
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
	models.AddPosts(maps)

	data["pictrue_url"] = imageUrl
	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

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
	ok := r.E(&valid, "Delete posts")
	if !ok {
		r.R(c, http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	models.DeletePostsUser(postId, uid)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
