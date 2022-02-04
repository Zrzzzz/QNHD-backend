package models

import (
	"gorm.io/gorm"
)

type PostReplyImage struct {
	PostReplyId uint64         `json:"post_reply_id"`
	ImageUrl    string         `json:"image_url"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func GetImageInPostReply(replyId uint64) ([]string, error) {
	var imageUrls = []string{}
	if err := db.Model(&PostReplyImage{}).Select("image_url").Where("post_reply_id = ?", replyId).Find(&imageUrls).Error; err != nil {
		return imageUrls, err
	}
	return imageUrls, nil
}

func AddImageInPostReply(tx *gorm.DB, replyId uint64, imageUrls []string) error {
	if tx == nil {
		tx = db
	}
	var pis = []PostReplyImage{}
	if len(imageUrls) == 0 {
		return nil
	}
	for _, url := range imageUrls {
		pis = append(pis, PostReplyImage{
			PostReplyId: replyId,
			ImageUrl:    url,
		})
	}
	err := tx.Create(&pis).Error
	return err
}

func DeleteImageInPostReply(tx *gorm.DB, replyId uint64) error {
	if tx == nil {
		tx = db
	}
	return tx.Where("post_reply_id = ?", replyId).Delete(&PostReplyImage{}).Error
}
