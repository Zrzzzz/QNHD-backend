package models

import (
	"errors"
	"qnhd/pkg/upload"

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

func DeleteImageInPost(tx *gorm.DB, postId uint64) error {
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
