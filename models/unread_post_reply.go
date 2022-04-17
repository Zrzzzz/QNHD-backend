package models

import (
	"qnhd/pkg/logging"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"

	"github.com/gin-gonic/gin"
)

type LogUnreadPostReply struct {
	Uid       uint64 `json:"uid"`
	ReplyId   uint64 `json:"reply_id"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

type UnreadReplyResponse struct {
	IsRead bool              `json:"is_read"`
	Post   PostResponse      `json:"post"`
	Reply  PostReplyResponse `json:"reply"`
}

func GetUnreadPostReplys(c *gin.Context, uid string) ([]UnreadReplyResponse, error) {
	var (
		err    error
		replys []PostReply
		logPrs []LogUnreadPostReply
		ret    = []UnreadReplyResponse{}
	)
	// 先筛选出未读记录
	if err = db.Model(&LogUnreadPostReply{}).Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("created_at DESC").Find(&logPrs).Error; err != nil {
		return ret, err
	}
	for _, log := range logPrs {
		var r PostReply
		if err = db.Where("id = ?", log.ReplyId).First(&r).Error; err != nil {
			continue
		}
		replys = append(replys, r)
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
		u.Post = p.geneResponse()
		// 加上未读
		for _, l := range logPrs {
			if l.ReplyId == r.Id {
				u.IsRead = l.IsRead
			}
		}

		ret = append(ret, u)
	}
	return ret, err
}

// 添加回复通知
func AddUnreadPostReply(postId, replyId uint64) error {
	var (
		uid  uint64
		user User
		post Post
	)
	if err := db.Find(&post, postId).Error; err != nil {
		return err
	}
	if err := db.Model(&Post{}).Select("uid").Where("id = ?", postId).Find(&uid).Error; err != nil {
		return err
	}
	if err := db.Where("id = ?", uid).Find(&user).Error; err != nil {
		return err
	}
	if err := db.Create(&LogUnreadPostReply{
		Uid:     uid,
		ReplyId: replyId,
		IsRead:  false,
	}).Error; err != nil {
		return err
	}
	if err := twtservice.NotifyPostReply(post.Title, user.Number); err != nil {
		logging.Error(err.Error())
	}
	return nil
}

// 已读回复
func ReadPostReply(uid, replyId uint64) error {
	return db.Model(&LogUnreadPostReply{}).
		Where("uid = ? AND reply_id = ?", uid, replyId).
		Update("is_read", true).Error
}

// 是否回复已读
func IsReadPostReply(uid, replyId uint64) bool {
	var log LogUnreadPostReply
	if err := db.Where("uid = ? AND reply_id = ?", uid, replyId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead
}
