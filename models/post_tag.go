package models

import (
	"errors"

	"gorm.io/gorm"
)

type PostTag struct {
	PostId uint64 `json:"post_id"`
	TagId  uint64 `json:"tag_id"`
}

func GetTagInPost(postId string) (*Tag, error) {
	var tag *Tag
	if err := db.Joins("JOIN qnhd.post_tag as pt ON qnhd.tag.id = pt.tag_id").Where("post_id = ?", postId).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return tag, err
		}
	}
	return tag, nil
}

func AddPostWithTag(tx *gorm.DB, postId uint64, tagId uint64) error {
	if tx == nil {
		tx = db
	}
	// 先查询是否有tag
	var tag Tag
	if err := tx.Where("id = ?", tagId).First(&tag).Error; err != nil {
		return err
	}
	err := tx.Create(&PostTag{
		PostId: postId,
		TagId:  tagId,
	}).Error
	return err
}

func DeleteTagInPost(tx *gorm.DB, postId uint64) error {
	if tx == nil {
		tx = db
	}
	err := tx.Where("post_id = ?", postId).Delete(&PostTag{}).Error
	return err
}
