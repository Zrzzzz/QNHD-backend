package models

import (
  "errors"

  "gorm.io/gorm"
)

type UserAvatarFrame struct{
  UId             uint64 `json:"uid"`
  AvatarFrameId   uint64 `json:"avatar_frame_id"`
  CreatedAt 	string `json:"created_at" gorm:"default:null;"`
}

// GetAddrById 通过 id 获取整个 UserAvatarFrame
func GetUserAvatarFrameById(id uint64) (ret UserAvatarFrame, err error) {
  err = db.Where("uid= ?", id).First(&ret).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return ret, nil
  }
  return ret, err
}

// AddNewUserAvatarFrame uid 对应用户获取 avatar frame
func AddNewUserAvatarFrame(uid, aid uint64) (user_avatar_frame UserAvatarFrame, err error) {
  user_avatar_frame = UserAvatarFrame{UId: uid, AvatarFrameId: aid}
  err = db.Select("uid", "avatar_frame_id").Create(&user_avatar_frame).Error
  return
}

// UpdateUserAvatarFrame uid 对应用户切换新的 avatar frame
func UpdateUserAvatarFrame(uid, aid uint64) (user_avatar_frame UserAvatarFrame, err error){
  db.First(&user_avatar_frame, "uid = ?", uid)
  err = db.Model(&user_avatar_frame).Where("uid", uid).Update("avatar_frame_id", aid).Error
  return
} 

func GetUserAvatarFrameAddr(id uint64) (addr string) {
  uaf := UserAvatarFrame{}
  err := db.Where("uid = ?", id).First(&uaf).Error
  if err != nil{
    addr = ""
  } else {
    af, err := GetAddrById(uaf.AvatarFrameId)
    if err != nil{
      addr = ""
    } else {
      addr = af.Addr
    }
  }
  return
}
