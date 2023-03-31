package models

import (
	"errors"
	"qnhd/pkg/logging"

	"gorm.io/gorm"
)

type AvatarFrame struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Addr      string `json:"addr"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
	Comment   string `json:"comment"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Hidden    string `json:"hidden"`
}

// GetAllAvatarFrames 获取所有 AvatarFrame
func GetAllAvatarFrames(f int) (avatar_frame_list []AvatarFrame, err error) {
	d := db
	if f == 1 {
		// 1 是必须要 hidden
		d = d.Where("hidden = ?", false)
	}
	err = d.Order("id").Find(&avatar_frame_list).Error
	return
}

// GetAddrById 通过 id 获取整个 AvatarFrame
func GetAddrById(id uint64) (avatar_frame AvatarFrame, err error) {
	err = db.Where("hidden = ?", false).First(&avatar_frame, id).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return avatar_frame, nil
	// }
	return avatar_frame, err
}

// AddNewAvatarFrame 添加一条数据
func AddNewAvatarFrame(addr string, comment string, t string, n string) (avatar_frame AvatarFrame, err error) {
	avatar_frame = AvatarFrame{Addr: addr, Comment: comment, Type: t, Name: n}
	err = db.Select("addr", "comment", "type", "name").Create(&avatar_frame).Error
	return
}

// UpdateAvatarFrame 更新存储地址
func UpdateAvatarFrame(id uint64, addr string, comment string, t string, n string, h string) (avatar_frame AvatarFrame, err error) {
	err = db.First(&avatar_frame, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return avatar_frame, err
	}
	avatar_frame.Comment = comment
	avatar_frame.Addr = addr
	avatar_frame.Type = t
	avatar_frame.Name = n
	avatar_frame.Hidden = h
	err = db.Save(&avatar_frame).Error
	return
}

// GetAddrByType 通过 type 获取 avatar_frame_list
func GetAddrByType(t string) (avatar_frame_list []AvatarFrame, err error) {
	if err = db.Where("type = ? AND hidden = ?", t, false).Order("id").Find(&avatar_frame_list).Error; err != nil {
		logging.Error("通过 type 获取 avatar_frame_list 错误")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
	}
	return
}
