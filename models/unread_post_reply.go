package models

import (
	"qnhd/pkg/logging"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

type LogUnreadPostReply struct {
	Uid       uint64 `json:"uid"`
	ReplyId   uint64 `json:"reply_id"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

type UnreadReplyResponse struct {
	Post  Post              `json:"post"`
	Reply PostReplyResponse `json:"reply"`
}

func GetUnreadPostReplys(c *gin.Context, uid string) ([]UnreadReplyResponse, error) {
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

// 已读通知
func ReadNotice(uid, noticeId uint64) error {
	return db.Where("uid = ? AND notice_id = ?", uid, noticeId).Delete(&LogUnreadNotice{}).Error
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

// 是否回复已读
func IsReadPostReply(uid, replyId uint64) bool {
	var log LogUnreadPostReply
	if err := db.Where("uid = ? AND reply_id = ?", uid, replyId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead
}
