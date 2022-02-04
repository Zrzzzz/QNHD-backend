package models

import (
	"errors"

	"gorm.io/gorm"
)

type LogUnreadNotice struct {
	Uid      uint64 `json:"uid"`
	NoticeId uint64 `json:"floor_id"`
}

// 通知所有用户
func addUnreadNoticeToAllUser(noticeId uint64) error {
	var userIds []uint64
	// 查询所有用户id
	if err := db.Model(&User{}).Select("id").Where("is_user = true AND active = true").Find(&userIds).Error; err != nil {
		return err
	}
	var logs = []LogUnreadNotice{}
	for _, id := range userIds {
		logs = append(logs, LogUnreadNotice{id, noticeId})
	}
	err := db.Create(logs).Error
	return err
}

// 是否通知已读
func IsReadNotice(uid, noticeId uint64) bool {
	var log LogUnreadNotice
	err := db.Where("uid = ? AND notice_id = ?", uid, noticeId).First(log).Error
	if err != nil && errors.Is(gorm.ErrRecordNotFound, err) {
		return true
	}
	return false
}
