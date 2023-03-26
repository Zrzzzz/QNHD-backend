package models

import (
	"errors"
	"qnhd/pkg/logging"

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
	err :=db.Model(&AvatarFrame{}).Select("avatar_frame.addr").Joins("JOIN qnhd.user_avatar_frame ON avatar_frame.id = qnhd.user_avatar_frame.avatar_frame_id").Where("user_avatar_frame.uid = ?", id).First(&addr).Error
	if err != nil{
		logging.Error("Get User Avatar Frame Addr by Uid (%v) Error: %v" ,id ,err)
		addr = ""
	}
	return
}
