package models

import (
	"errors"
	"qnhd/pkg/util"

	"gorm.io/gorm"
)

type PostTag struct {
	Id     uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	PostId uint64 `json:"post_id" `
	TagId  uint64 `json:"tag_id" `
}

func GetTagInPost(postId string) (Tag, error) {
	var tag Tag
	if err := db.Joins("JOIN post_tag ON tags.id = post_tag.tag_id").Where("post_id = ?", postId).Find(&tag).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return tag, err
		} else {
			return tag, nil
		}
	}
	return tag, nil
}

func AddPostWithTag(postId uint64, tagId string) error {
	// 先查询是否有tag
	var tag Tag
	if err := db.Where("id = ?", tagId).First(&tag).Error; err != nil {
		return err
	}
	err := db.Create(&PostTag{
		PostId: postId,
		TagId:  util.AsUint(tagId),
	}).Error
	return err
}

func DeleteTagInPost(tx *gorm.DB, postId string) error {
	if tx == nil {
		tx = db
	}
	err := tx.Where("post_id = ?", postId).Delete(&PostTag{}).Error
	return err
}

func (PostTag) TableName() string {
	return "post_tag"
}
