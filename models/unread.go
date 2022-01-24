package models

import (
	"errors"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogUnreadFloor struct {
	Id      uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Uid     uint64 `json:"uid"`
	FloorId uint64 `json:"floor_id"`
	IsRead  int    `json:"is_read"`
}

func (LogUnreadFloor) TableName() string {
	return "log_unread_floor"
}

type LogUnreadNotice struct {
	Uid      uint64 `json:"uid"`
	NoticeId uint64 `json:"floor_id"`
}

func (LogUnreadNotice) TableName() string {
	return "log_unread_notice"
}

type LogUnreadPostReply struct {
	Id      uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Uid     uint64 `json:"uid"`
	ReplyId uint64 `json:"floor_id"`
	IsRead  int    `json:"is_read"`
}

func (LogUnreadPostReply) TableName() string {
	return "log_unread_post_reply"
}

func GetMessageFloors(c *gin.Context, uid string) ([]LogUnreadFloor, error) {
	var logs []LogUnreadFloor
	// 未读的优先，按照时间
	err := db.Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("`is_read`, `id` DESC").Find(&logs).Error
	return logs, err
}

func GetMessagePostReplys(c *gin.Context, uid string) ([]LogUnreadPostReply, error) {
	var logs []LogUnreadPostReply
	err := db.Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("`is_read`, `id` DESC").Find(&logs).Error
	return logs, err
}

// 通知所有用户
func addUnreadNoticeToAllUser(noticeId uint64) error {
	var userIds []uint64
	// 查询所有用户id
	if err := db.Model(&User{}).Select("id").Where("is_user = 1 AND active = 1").Find(&userIds).Error; err != nil {
		return err
	}
	var logs = []LogUnreadNotice{}
	for _, id := range userIds {
		logs = append(logs, LogUnreadNotice{id, noticeId})
	}
	err := db.Create(logs).Error
	return err
}

// 已读通知
func ReadNotice(uid, noticeId uint64) error {
	return db.Where("uid = ? AND notice_id = ?", uid, noticeId).Delete(&LogUnreadNotice{}).Error
}

// 添加评论通知
func addUnreadFloor(uid, floorId uint64) error {
	return db.Create(&LogUnreadFloor{
		Uid:     uid,
		FloorId: floorId,
	}).Error
}

// 已读评论
func ReadFloor(uid, floorId uint64) error {
	return db.Model(&LogUnreadFloor{}).
		Where("uid = ? AND floor_id = ?", uid, floorId).
		Update("`is_read`", 1).Error
}

// 添加回复通知
func AddUnreadPostReply(postId, replyId uint64) error {
	var uid uint64
	if err := db.Model(&Post{}).Select("uid").Where("id = ?", postId).Find(&uid).Error; err != nil {
		return err
	}
	return db.Select("uid", "reply_id").Create(&LogUnreadPostReply{
		Uid:     uid,
		ReplyId: replyId,
	}).Error
}

// 已读回复
func ReadPostReply(uid, replyId uint64) error {
	return db.Model(&LogUnreadPostReply{}).
		Where("uid = ? AND reply_id = ?", uid, replyId).
		Update("`is_read`", 1).Error
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

// 是否评论已读
func IsReadFloor(uid, floorId uint64) bool {
	var log LogUnreadFloor
	if err := db.Where("uid = ? AND floorId = ?", uid, floorId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead == 1
}

// 是否回复已读
func IsReadPostReply(uid, replyId uint64) bool {
	var log LogUnreadPostReply
	if err := db.Where("uid = ? AND reply_id = ?", uid, replyId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead == 1
}

type MessageCount struct {
	Floor  int `json:"floor"`
	Reply  int `json:"reply"`
	Notice int `json:"notice"`
}

// 获取总未读数
func GetMessageCount(uid string) (MessageCount, error) {
	var ret = MessageCount{}
	// 楼层未读 回复未读 通知未读
	var fcnt, rcnt, ncnt int64
	// 获取楼层未读数
	if err := db.Model(&LogUnreadFloor{}).Where("uid = ? AND is_read = 0", uid).Count(&fcnt).Error; err != nil {
		return ret, err
	}
	// 获取回复未读数
	if err := db.Model(&LogUnreadPostReply{}).Where("uid = ? AND is_read = 0", uid).Count(&rcnt).Error; err != nil {
		return ret, err
	}
	// 获取通知未读数
	if err := db.Model(&LogUnreadNotice{}).Where("uid = ?", uid).Count(&ncnt).Error; err != nil {
		return ret, err
	}
	ret.Floor = int(fcnt)
	ret.Reply = int(rcnt)
	ret.Notice = int(ncnt)
	return ret, nil
}
