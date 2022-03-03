package models

import (
	"qnhd/pkg/filter"

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
	err := db.Where("post_id = ?", postId).Find(&prs).Order("id").Error
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
		Content: filter.Filter(maps["content"].(string)),
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
		images []PostReplyImage
	)
	err := ttx.Transaction(func(tx *gorm.DB) error {
		// 获取所有image
		logs := tx.Model(&PostReply{}).Where("post_id = ?", postId)
		if err := tx.Table("(?) as a", logs).
			Select("b.*").
			Joins("JOIN qnhd.post_reply_image as b ON a.id = b.post_reply_id").
			Find(&images).
			Error; err != nil {
			return err
		}
		if len(images) == 0 {
			return nil
		}
		// 删除reply
		if err := tx.Where("post_id = ?", postId).Delete(&PostReply{}).Error; err != nil {
			return err
		}
		// 删除记录
		if err := tx.Where("post_id = ?", postId).Delete(&LogUnreadPostReply{}).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
