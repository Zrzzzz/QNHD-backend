package models

import (
	"errors"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogUnreadFloor struct {
	Uid       uint64 `json:"uid"`
	FloorId   uint64 `json:"floor_id"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

type LogUnreadNotice struct {
	Uid      uint64 `json:"uid"`
	NoticeId uint64 `json:"floor_id"`
}

type LogUnreadPostReply struct {
	Uid       uint64 `json:"uid"`
	ReplyId   uint64 `json:"reply_id"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

type UnreadFloorResponse struct {
	Type    int    `json:"type"`
	ToFloor *Floor `json:"to_floor"`
	Post    Post   `json:"post"`
	Floor   Floor  `json:"floor"`
}

type UnreadReplyResponse struct {
	Post  Post              `json:"post"`
	Reply PostReplyResponse `json:"reply"`
}

func GetMessageFloors(c *gin.Context, uid string) ([]UnreadFloorResponse, error) {
	var (
		ret    = []UnreadFloorResponse{}
		floors []Floor
		err    error
	)

	// 先筛选出未读记录
	logs := db.Model(&LogUnreadFloor{}).Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("is_read, created_at DESC")
	// 找到楼层
	if err = db.Table("(?) as a", logs).
		Unscoped().
		Select("f.*").
		Joins("JOIN qnhd.floor as f ON a.floor_id = f.id").
		Find(&floors).
		Where("f.deleted_at IS NULL").
		Error; err != nil {
		return ret, err
	}
	// 对每个楼层分析
	for _, f := range floors {
		var r = UnreadFloorResponse{Floor: f}
		// 搜索floor
		if f.SubTo > 0 {
			tof, e := GetFloor(util.AsStrU(f.ReplyTo))
			if e != nil {
				err = e
				break
			}
			r.Type = 1
			r.ToFloor = &tof
		} else {
			r.Type = 0
		}
		// 搜索帖子
		p, e := GetPost(util.AsStrU(f.PostId))
		if e != nil {
			err = e
			break
		}
		r.Post = p
		ret = append(ret, r)
	}
	return ret, err
}

func GetMessagePostReplys(c *gin.Context, uid string) ([]UnreadReplyResponse, error) {
	var (
		err    error
		replys []PostReply
		ret    = []UnreadReplyResponse{}
	)
	// 先筛选出未读记录
	logs := db.Model(&LogUnreadPostReply{}).Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("is_read, created_at DESC")
	// 找到回复
	if err = db.Table("(?) as a", logs).
		Unscoped().
		Select("pr.*").
		Joins("JOIN qnhd.post_reply as pr ON a.reply_id = pr.id").
		Find(&replys).
		Where("pr.deleted_at IS NULL").
		Error; err != nil {
		return ret, err
	}
	// 再生成返回数据
	for _, r := range replys {
		rp, e := r.geneResponse()
		if e != nil {
			logging.Error(e.Error())
		}
		var u = UnreadReplyResponse{Reply: rp}
		// 搜索帖子
		p, e := GetPost(util.AsStrU(rp.PostId))
		if e != nil {
			err = e
			break
		}
		u.Post = p
		ret = append(ret, u)
	}
	return ret, err
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
		Update("is_read", 1).Error
}

// 添加回复通知
func AddUnreadPostReply(postId, replyId uint64) error {
	var uid uint64
	if err := db.Model(&Post{}).Select("uid").Where("id = ?", postId).Find(&uid).Error; err != nil {
		return err
	}
	return db.Create(&LogUnreadPostReply{
		Uid:     uid,
		ReplyId: replyId,
		IsRead:  false,
	}).Error
}

// 已读回复
func ReadPostReply(uid, replyId uint64) error {
	return db.Model(&LogUnreadPostReply{}).
		Where("uid = ? AND reply_id = ?", uid, replyId).
		Update("is_read", 1).Error
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
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead
}

// 是否回复已读
func IsReadPostReply(uid, replyId uint64) bool {
	var log LogUnreadPostReply
	if err := db.Where("uid = ? AND reply_id = ?", uid, replyId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead
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
	if err := db.Model(&LogUnreadFloor{}).Where("uid = ? AND is_read = true", uid).Count(&fcnt).Error; err != nil {
		return ret, err
	}
	// 获取回复未读数
	if err := db.Model(&LogUnreadPostReply{}).Where("uid = ? AND is_read = true", uid).Count(&rcnt).Error; err != nil {
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
