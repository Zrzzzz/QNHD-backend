package models

import (
	"qnhd/pkg/upload"

	"github.com/pkg/errors"
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
	Sender  PostReplyType `json:"sender"`
	Content string        `json:"content"`
}

func (PostReply) TableName() string {
	return "post_reply"
}

type PostReplyResponse struct {
	PostReply
	ImageUrls []string `json:"image_urls"`
}

// 转换为response
func (p *PostReply) geneResponse() (PostReplyResponse, error) {
	var prr = PostReplyResponse{PostReply: *p}
	err := db.Model(&PostReplyImage{}).Select("image_url").Where("post_reply_id = ?", prr.Id).Find(&prr.ImageUrls).Error
	return prr, err
}

// 获取单个回复
func GetPostReply(replyId string) (PostReply, error) {
	var pr PostReply
	err := db.Where("id = ?", replyId).Find(&pr).Error
	return pr, err
}

// 获取单个带图片回复
func GetPostReplyResponse(replyId string) (PostReplyResponse, error) {
	var prr PostReplyResponse
	pr, err := GetPostReply(replyId)
	if err != nil {
		return prr, err
	}
	return pr.geneResponse()
}

// 获取帖子内的回复记录
func GetPostReplys(postId string) ([]PostReply, error) {
	var prs = []PostReply{}
	err := db.Where("post_id = ?", postId).Find(&prs).Error
	return prs, err
}

// 获取带图片回复
func GetPostReplyResponses(postId string) ([]PostReplyResponse, error) {
	var rets []PostReplyResponse
	var err error
	prs, err := GetPostReplys(postId)
	if err != nil {
		return rets, err
	}
	for _, pr := range prs {
		ret, e := pr.geneResponse()
		if e != nil {
			err = errors.Wrap(err, e.Error())
		} else {
			rets = append(rets, ret)
		}
	}
	return rets, err
}

// 添加帖子的回复
func AddPostReply(maps map[string]interface{}) (uint64, error) {
	var pr = PostReply{
		PostId:  maps["post_id"].(uint64),
		Sender:  maps["sender"].(PostReplyType),
		Content: maps["content"].(string),
	}
	urls := maps["urls"].([]string)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&pr).Error; err != nil {
			return err
		}
		if len(urls) != 0 {
			if err := AddImageInPostReply(tx, pr.Id, urls); err != nil {
				return err
			}
		}
		return nil
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
