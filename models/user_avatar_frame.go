package models

import (
	"errors"
	"qnhd/pkg/logging"
	"gorm.io/gorm"
)

type UserAvatarFrame struct {
	UId           uint64 `json:"uid"`
	AvatarFrameId uint64 `json:"avatar_frame_id"`
	CreatedAt     string `json:"created_at" gorm:"default:null;"`
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
	if err = db.Where("uid= ?", uid).First(&UserAvatarFrame{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user_avatar_frame = UserAvatarFrame{UId: uid, AvatarFrameId: aid}
		err = db.Select("uid", "avatar_frame_id").Create(&user_avatar_frame).Error
		return
	}
	db.First(&user_avatar_frame, "uid = ?", uid)
	err = db.Model(&user_avatar_frame).Where("uid", uid).Update("avatar_frame_id", aid).Error
	return
}

func GetUserAvatarFrameAddr(id uint64) (addr string) {
    // NOTE: 这种写法会导致Log输出错误信息
    // * 从逻辑上来说没问题，但是由于很多人没有设置头像框，导致很多人陷入'NOT FOUND' 的情况，在Go层面是Error，会直接输出到Log中
    // * 建议在 sql 层面解决
	// err := db.Model(&AvatarFrame{}).Select("avatar_frame.addr").Joins("JOIN qnhd.user_avatar_frame ON avatar_frame.id = qnhd.user_avatar_frame.avatar_frame_id").Where("user_avatar_frame.uid = ? AND avatar_frame.hidden = ?", id, false).First(&addr).Error
    
	err := db.Model(&AvatarFrame{}).Select("coalesce(avatar_frame.addr, '')").Joins("JOIN qnhd.user_avatar_frame ON avatar_frame.id = qnhd.user_avatar_frame.avatar_frame_id").Where("user_avatar_frame.uid = ? AND avatar_frame.hidden = ?", id, false).First(&addr).Error
	if err != nil {
        // 使用 coalesce 之后，不会因为没有设置头像框而出现error，则可以认为是错误
		logging.Error("Get User Avatar Frame Addr by Uid (%v) Error: %v" ,id ,err)
		addr = "Error"
	}
	return
}
