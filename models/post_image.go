package models

import (
	"gorm.io/gorm"
)

type PostImage struct {
	PostId   uint64 `json:"post_id" `
	ImageUrl string `json:"image_url" `
}

func GetImageInPost(postId uint64) ([]string, error) {
	var imageUrls = []string{}
	var ret = []PostImage{}
	if err := db.Where("post_id = ?", postId).Find(&ret).Error; err != nil {
		return imageUrls, err
	}
	for _, r := range ret {
		imageUrls = append(imageUrls, r.ImageUrl)
	}
	return imageUrls, nil
}

func AddImageInPost(tx *gorm.DB, postId uint64, imageUrls []string) error {
	if tx == nil {
		tx = db
	}
	var pis = []PostImage{}
	if len(imageUrls) == 0 {
		return nil
	}
	for _, url := range imageUrls {
		pis = append(pis, PostImage{
			PostId:   postId,
			ImageUrl: url,
		})
	}
	err := tx.Create(&pis).Error
	return err
}
