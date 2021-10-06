package models

import (
	"errors"
	"strconv"

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
	err := db.Transaction(func(tx *gorm.DB) error {
		addDb := tx.Select("PostId", "TagId")
		intt, _ := strconv.ParseUint(tagId, 10, 64)
		if err := addDb.Create(&PostTag{
			PostId: postId,
			TagId:  intt,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteTagInPost(postId string) error {
	if err := db.Where("post_id = ?", postId).Delete(&PostTag{}).Error; err != nil {
		return err
	}
	return nil
}

func (PostTag) TableName() string {
	return "post_tag"
}
