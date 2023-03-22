package models

import (
	"errors"

	"gorm.io/gorm"
)

type AvatarFrame struct {
	Id        	uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Addr        string `json:"addr"`
	CreatedAt 	string `json:"created_at" gorm:"default:null;"`
	Comment     string `json:"comment"`
}

// GetAllAvatarFrames 获取所有 AvatarFrame
func GetAllAvatarFrames() (avatar_frame_list []AvatarFrame, err error){
	err = db.Find(&avatar_frame_list).Error
	return 
}
// GetAddrById 通过 id 获取整个 AvatarFrame
func GetAddrById(id uint64) (avatar_frame AvatarFrame, err error) {
	err = db.First(&avatar_frame, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return avatar_frame, nil
	}
	return avatar_frame, err
}

// AddNewAvatarFrame 添加一条数据
func AddNewAvatarFrame(addr string, comment string) (avatar_frame AvatarFrame, err error) {
	avatar_frame = AvatarFrame{Addr: addr, Comment: comment}
	err = db.Select("addr", "comment").Create(&avatar_frame).Error
	return
}

// UpdateAvatarFrame 更新存储地址
func UpdateAvatarFrame(id uint64, addr string, comment string) (avatar_frame AvatarFrame, err error){
	err = db.First(&avatar_frame, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return avatar_frame, err
	}
	avatar_frame.Comment = comment
	avatar_frame.Addr = addr
	err = db.Save(&avatar_frame).Error
	return 
}
