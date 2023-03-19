package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /b/frame/all
func GetAllAvatarFrame(c *gin.Context) {
	avatar_frame_list, err := models.GetAllAvatarFrames()
	if err != nil {
		logging.Error("Get all avatar frame Error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame_list": avatar_frame_list})
}

// @method [post]
// @way [query]
// @param
// @return
// @route /b/frame/upload
// UploadAvatarFrame 存储一下新的头像框
func UploadAvatarFrame(c *gin.Context) {
	addr := c.PostForm("addr")
	comment := c.PostForm("comment")
	Ret, err := models.AddNewAvatarFrame(addr, comment)
	if err != nil{
		logging.Error("Upload Avatar Error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame": Ret})
}

// @method [post]
// @way [query]
// @param
// @return
// @route /b/frame/upload
func UpdateAvatarFrame(c *gin.Context){
	id := c.PostForm("id")
	addr := c.PostForm("addr")
	comment := c.PostForm("comment")
	ret, err := models.UpdateAvatarFrame(util.AsUint(id), addr, comment)
	if err != nil {
		logging.Error("Update Avatar Frame Error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame": ret})
}
