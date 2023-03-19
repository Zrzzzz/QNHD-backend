package frontend

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
// @route /f/frame/my
func GetMyFrame(c *gin.Context) {
	uid := r.GetUid(c)
	avatar_frame, err := models.GetUserAvatarFrameById(util.AsUint(uid))
	logging.Debug("avatar_Frame: %v", avatar_frame)
	if err != nil {
		logging.Error("get user frame error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame": avatar_frame})
}

// @method [get]
// @may [query]
// @param
// @return
// @route /f/frame/set
// SetMyFrame 给当前用户设置头相框
func SetMyFrame(c *gin.Context){
	uid := r.GetUid(c)
	aid := c.PostForm("aid")
	user_avatar_frame, err := models.AddNewUserAvatarFrame(util.AsUint(uid), util.AsUint(aid))
  if err != nil {
		logging.Error("Add user frame Error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
  r.OK(c, e.SUCCESS, map[string]interface{}{"user_avatar_frame": user_avatar_frame})
}

func UpdateMyFrame(c *gin.Context){
  uid := r.GetUid(c)
  aid := c.PostForm("aid")
  user_avatar_frame, err := models.UpdateUserAvatarFrame(util.AsUint(uid), util.AsUint(aid))
  if err != nil {
    logging.Error("Update user frame Error: %v", err)
    r.Error(c, e.ERROR_DATABASE, err.Error())
    return
  }
  r.OK(c, e.SUCCESS, map[string]interface{}{"user_avatar_frame": user_avatar_frame})
}

