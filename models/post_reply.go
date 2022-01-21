package models

import (
	"qnhd/pkg/upload"

	"gorm.io/gorm"
)

type PostReplyType int

const (
	PostReplyFromUser PostReplyType = iota
	PostReplyFromSchool
)

type PostReply struct {
	Model
	PostId  uint64        `json:"post_id"`
	From    PostReplyType `json:"from"`
	Content string        `json:"content"`
}

// 获取单个回复
func GetPostReply(replyId string) (PostReply, error) {
	var pr PostReply
	err := db.Where("id = ?", replyId).Find(&pr).Error
	return pr, err
}

// 获取帖子内的回复记录
func GetPostReplys(postId string) ([]PostReply, error) {
	var prs = []PostReply{}
	err := db.Where("post_id = ?", postId).Find(&prs).Error
	return prs, err
}

// 添加帖子的回复
func AddPostReply(maps map[string]interface{}) (uint64, error) {
	var pr = PostReply{
		PostId:  maps["post_id"].(uint64),
		From:    maps["sender"].(PostReplyType),
		Content: maps["content"].(string),
	}
	urls := maps["urls"].([]string)
	err := db.Transaction(func(tx *gorm.DB) error {
		if len(urls) != 0 {
			if err := AddImageInPostReply(tx, maps["post_id"].(uint64), urls); err != nil {
				return err
			}
		}
		return tx.Create(&pr).Error
	})
	return pr.Id, err
}

func DeletePostReplysInPost(ttx *gorm.DB, postId uint64) error {
	if ttx == nil {
		ttx = db
	}
	var (
		replys   []PostReply
		replyIds []uint64
		urls     []string
	)
	err := ttx.Transaction(func(tx *gorm.DB) error {
		// TODO: 做连表
		// 获取所有reply
		if err := tx.Where("post_id = ?", postId).Find(&replys).Error; err != nil {
			return err
		}
		if len(replyIds) == 0 {
			return nil
		}
		for _, r := range replys {
			replyIds = append(replyIds, r.Id)
		}
		// 获取所有图片
		if err := tx.Model(&PostReplyImage{}).Select("image_url").Where("post_reply_id IN (?)", urls).Find(&urls).Error; err != nil {
			return err
		}
		// 删除本地图片
		if err := upload.DeleteImageUrls(urls); err != nil {
			return err
		}
		// 删除reply
		if err := tx.Where("post_id = ?", postId).Delete(&replys).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (PostReply) TableName() string {
	return "post_reply"
}
