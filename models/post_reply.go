package models

import (
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/enums/PostReplyType"
	"qnhd/pkg/filter"
	"qnhd/pkg/util"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PostReply struct {
	Model
	PostId  uint64             `json:"post_id"`
	Sender  PostReplyType.Enum `json:"sender"`
	Content string             `json:"content"`
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
	err := db.Where("post_id = ?", postId).Order("id").Find(&prs).Error
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
	sender := maps["sender"].(PostReplyType.Enum)
	var pr = PostReply{
		PostId:  maps["post_id"].(uint64),
		Sender:  sender,
		Content: filter.CommonFilter.Filter(maps["content"].(string)),
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
	if sender == PostReplyType.SCHOOL {
		uid := maps["uid"].(string)
		addManagerLog(util.AsUint(uid), pr.Id, ManagerLogType.POST_REPLY)
	}
	return pr.Id, err
}

func EditPostReply(maps map[string]interface{}) error {
	sender := maps["sender"].(PostReplyType.Enum)
	var pr = PostReply{
		Sender:  sender,
		Content: filter.CommonFilter.Filter(maps["content"].(string)),
	}
	urls := maps["urls"].([]string)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&PostReply{}).Where("id = ?", maps["reply_id"].(uint64)).Updates(pr).Error; err != nil {
			return err
		}
		if err := DeleteImageInPostReply(tx, pr.Id); err != nil {
			return err
		}
		if len(urls) != 0 {
			if err := AddImageInPostReply(tx, pr.Id, urls); err != nil {
				return err
			}
		}
		return nil
	})
	if sender == PostReplyType.SCHOOL {
		uid := maps["uid"].(string)
		addManagerLog(util.AsUint(uid), pr.Id, ManagerLogType.POST_REPLY_MODIFY)
	}
	return err
}

func DeletePostReply(id string) error {
	return db.Where("id = ?").Delete(&PostReply{}).Error
}

// 删除帖子内的回复记录
func DeletePostReplysInPost(ttx *gorm.DB, postId uint64) error {
	if ttx == nil {
		ttx = db
	}
	var (
		replies []PostReply
	)
	err := ttx.Transaction(func(tx *gorm.DB) error {
		logs := tx.Model(&PostReply{}).Where("post_id = ?", postId)
		if err := logs.Find(&replies).Error; err != nil {
			return err
		}
		// 删除reply
		if err := tx.Where("post_id = ?", postId).Delete(&PostReply{}).Error; err != nil {
			return err
		}
		for _, r := range replies {
			// 删除记录
			if err := tx.Where("reply_id = ?", r.Id).Delete(&LogUnreadPostReply{}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// 恢复帖子内的回复记录
func RecoverPostReplysInPost(tx *gorm.DB, postId uint64) error {
	if tx == nil {
		tx = db
	}

	// 删除reply
	return tx.Unscoped().Model(&PostReply{}).Where("post_id = ?", postId).Update("deleted_at", gorm.Expr("NULL")).Error
}
