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

// @method [get]
// @way [query]
// @param
// @return
// @route /f/frame/all
func GetAllAvatarFrame(c *gin.Context) {
  avatar_frame_list, err := models.GetAllAvatarFrames(1)
  if err != nil {
    logging.Error("Get all avatar frame Error: %v", err)
    r.Error(c, e.ERROR_DATABASE, err.Error())
    return
  }
  r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame_list": avatar_frame_list, "total": len(avatar_frame_list)})
}


// @method [get]
// @way [query]
// @param
// @return
// @route /f/frame/id_url
func GetAvatarFrameUrlById(c *gin.Context) {
	aid := c.Query("aid")
  avatar_frame, err := models.GetAddrById(util.AsUint(aid))
  if err != nil {
    logging.Error("Get Avatar Frame Error: %v", aid)
    r.Error(c, e.ERROR_DATABASE, err.Error())
    return
  }
  r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame": avatar_frame})
}

// @method [get]
// @way [query]
// @param
// @return
// @route /f/frame/type_url
func GetAvatarFrameUrlByType(c *gin.Context) {
	t := c.Query("type")
  avatar_frame_list, err := models.GetAddrByType(t)
  if err != nil {
    logging.Error("Get Avatar Frame Error Type: %v", t)
    r.Error(c, e.ERROR_DATABASE, err.Error())
    return
  }
  r.OK(c, e.SUCCESS, map[string]interface{}{"avatar_frame_list": avatar_frame_list, "total": len(avatar_frame_list)})
}




