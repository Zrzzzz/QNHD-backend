package models

import (
	"errors"
	"qnhd/pkg/upload"

	"gorm.io/gorm"
)

type PostImage struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	PostId   uint64 `json:"post_id" `
	ImageUrl string `json:"image_url" `
}

func GetImageInPost(postId string) ([]string, error) {
	var imageUrls = []string{}
	var ret = []PostImage{}
	if err := db.Where("post_id = ?", postId).Find(&ret).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return imageUrls, nil
		} else {
			return imageUrls, err
		}
	}
	for _, r := range ret {
		imageUrls = append(imageUrls, r.ImageUrl)
	}
	return imageUrls, nil
}

func AddImageInPost(postId uint64, imageUrls []string) error {
	var pis = []PostImage{}
	if len(pis) == 0 {
		return nil
	}
	for _, url := range imageUrls {
		pis = append(pis, PostImage{
			PostId:   postId,
			ImageUrl: url,
		})
	}
	err := db.Create(&pis).Error
	return err
}

func DeleteImageInPost(tx *gorm.DB, postId string) error {
	if tx == nil {
		tx = db
	}
	// 先删除本地文件
	imageUrls, err := GetImageInPost(postId)
	if err != nil {
		return err
	}
	err = upload.DeleteImageUrls(imageUrls)
	if err != nil {
		return err
	}
	err = tx.Where("post_id = ?", postId).Delete(&PostImage{}).Error
	return err
}

func (PostImage) TableName() string {
	return "post_image"
}
