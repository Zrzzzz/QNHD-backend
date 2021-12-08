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
	var ret []PostImage
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
	err := db.Transaction(func(tx *gorm.DB) error {
		addDb := tx.Select("post_id", "image_url")
		for _, url := range imageUrls {
			if err := addDb.Create(&PostImage{
				PostId:   postId,
				ImageUrl: url,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
	if err != nil {
		return err
	}
	return nil
}

func (PostImage) TableName() string {
	return "post_image"
}
