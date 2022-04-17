package models

import (
	"errors"
	"math"
	"qnhd/pkg/logging"
	"qnhd/pkg/template"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UnreadNoticeResponse struct {
	Notice
	IsRead bool `json:"is_read" gorm:"-"`
}

type LogUnreadNotice struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Uid      uint64 `json:"uid"`
	NoticeId uint64 `json:"notice_id"`
	Args     string `json:"args"`
	IsRead   bool   `json:"is_read" gorm:"default:false"`
	PubAt    string `json:"pub_at" gorm:"default:null"`
}

type userResult struct {
	Id     uint64 `json:"id"`
	Number string `json:"number"`
}

type noticeResult struct {
	LogUnreadNotice
	Notice
}

// 获取未读的所有notice
func GetUnreadNotices(c *gin.Context, uid uint64) ([]UnreadNoticeResponse, error) {
	var (
		logs []noticeResult
		ret  = []UnreadNoticeResponse{}
	)
	p := db.Model(&LogUnreadNotice{}).Where("uid = ? AND pub_at < ?", uid, gorm.Expr("CURRENT_TIMESTAMP"))
	if err := db.Debug().Unscoped().Table("(?) as p", p).
		Select("p.*, n.*").
		Joins("JOIN qnhd.notice as n ON n.id = p.notice_id").
		Order("p.id DESC").
		Find(&logs).Error; err != nil {
		return ret, err
	}
	for _, log := range logs {
		var resp = UnreadNoticeResponse{
			Notice: log.Notice,
		}
		resp.IsRead = log.IsRead
		// 模板进行替换
		resp.Content, _ = template.GeneTemplateString(log.Content, log.Args)
		ret = append(ret, resp)
	}
	return ret, nil
}

// 通知所有用户
func addUnreadNoticeToAllUser(noticeId uint64, pubAt string) error {
	var (
		notice  Notice
		users   []userResult
		numbers []string
	)
	if err := db.Where("id = ?", noticeId).Find(&notice).Error; err != nil {
		return err
	}
	// 查询所有用户id
	if err := db.Model(&User{}).Select("id", "number").Where("is_user = true AND active = true").Find(&users).Error; err != nil {
		return err
	}
	var logs []LogUnreadNotice
	for _, u := range users {
		logs = append(logs, LogUnreadNotice{Uid: u.Id, NoticeId: notice.Id, PubAt: pubAt})
		numbers = append(numbers, u.Number)
	}
	// 一次插入2个参数，只要少于65535就ok，经测试250效率较高
	insertCount := 250
	for i := 0; i < int(math.Ceil(float64(len(logs))/float64(insertCount))); i++ {
		min := (i + 1) * insertCount
		if len(logs) < min {
			min = len(logs)
		}
		db.Create(logs[i*insertCount : min])
	}
	if err := twtservice.NotifyNotice(notice.Sender, notice.Title, numbers...); err != nil {
		logging.Error(err.Error())
	}
	return nil
}

// 模板通知用户
func addUnreadNoticeToUser(uid []uint64, data map[string]interface{}) error {
	var notice Notice
	if err := db.Where("symbol = ?", data["symbol"].(string)).Find(&notice).Error; err != nil {
		return err
	}
	var logs []LogUnreadNotice
	var uidStrs []string
	for _, u := range uid {
		logs = append(logs, LogUnreadNotice{
			Uid:      u,
			NoticeId: notice.Id,
			Args:     data["args"].(string),
		})
		uidStrs = append(uidStrs, util.AsStrU(u))
	}
	insertCount := 250
	for i := 0; i < int(math.Ceil(float64(len(logs))/float64(insertCount))); i++ {
		min := (i + 1) * insertCount
		if len(logs) < min {
			min = len(logs)
		}
		db.Create(logs[i*insertCount : min])
	}
	if err := twtservice.NotifyNotice(notice.Sender, notice.Title, uidStrs...); err != nil {
		logging.Error(err.Error())
	}
	return nil
}

// 已读通知
func ReadNotice(uid, noticeId uint64) error {
	return db.Where("uid = ? AND notice_id = ?", uid, noticeId).Delete(&LogUnreadNotice{}).Error
}

// 删除通知记录
func DeleteMessageNotices(uid string, ids []string) error {
	return db.Where("uid = ? AND id IN (?)", uid, ids).Delete(&LogUnreadNotice{}).Error
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
