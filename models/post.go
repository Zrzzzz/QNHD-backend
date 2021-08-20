package models

import (
	"errors"

	"gorm.io/gorm"
)

type Post struct {
	Model
	Uid        uint64 `json:"uid"`
	Content    string `json:"content"`
	PictureUrl string `json:"picture_url"`
	UpdatedAt  string `json:"updated_at" gorm:"null;"`
}

func GetPost(id string) (Post, error) {
	var post Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, nil
		}
		return post, err
	}
	return post, nil
}

func GetPosts(overNum, limit int, content string) ([]Post, error) {
	var posts []Post
	if err := db.Where("content LIKE ?", "%"+content+"%").Offset(overNum).Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func AddPosts(maps map[string]interface{}) (uint64, error) {
	var post = &Post{
		Uid:        maps["uid"].(uint64),
		Content:    maps["content"].(string),
		PictureUrl: maps["picture_url"].(string),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("uid", "content", "picture_url").Create(post).Error; err != nil {
			return err
		}
		tags, ok := maps["tags"].([]string)
		if ok {
			if err := AddPostWithTag(post.Id, tags); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return post.Id, nil
}

func DeletePostsUser(id, uid string) (uint64, error) {
	var post = Post{}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND uid = ?", id, uid).Delete(&post).Error; err != nil {
			return err
		}
		if err := DeleteTagInPost(id); err != nil {
			return err
		}
		if _, err := DeleteFloorsInPost(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return post.Id, nil
}

func DeletePostsAdmin(id string) (uint64, error) {
	var post = Post{}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&post).Error; err != nil {
			return err
		}
		if err := DeleteTagInPost(id); err != nil {
			return err
		}
		if _, err := DeleteFloorsInPost(id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return post.Id, nil
}

func (Post) TableName() string {
	return "posts"
}
