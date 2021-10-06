package models

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Post struct {
	Model
	Uid        uint64 `json:"uid"`
	Content    string `json:"content"`
	PictureUrl string `json:"picture_url"`
	Favs       uint64 `json:"favs"`
	UpdatedAt  string `json:"updated_at" gorm:"null;"`
}

type LogPostFav struct {
	Model
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}

func GetPost(id string, uid string) (Post, error) {
	var post Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, nil
		}
		return post, err
	}
	if _, err := AddVisitHistory(uid, id); err != nil {
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

func AddPost(maps map[string]interface{}) (uint64, error) {
	var post = &Post{
		Uid:        maps["uid"].(uint64),
		Content:    maps["content"].(string),
		PictureUrl: maps["picture_url"].(string),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("uid", "content", "picture_url").Create(post).Error; err != nil {
			return err
		}
		tagId, ok := maps["tag_id"].(string)
		if ok {
			if err := AddPostWithTag(post.Id, tagId); err != nil {
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

func FavPost(postId string, uid string) error {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	postIdint, _ := strconv.ParseUint(postId, 10, 64)

	var exist = false
	var log = LogPostFav{Uid: uidint, PostId: postIdint}

	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	if log.Id > 0 {
		return fmt.Errorf("已收藏")
	}

	if err := db.Unscoped().Where(log).Find(&log).Error; err != nil {
		return err
	}

	exist = log.Id > 0
	if exist {
		if err := db.Unscoped().Model(&log).Update("deleted_at", gorm.Expr("NULL")).Error; err != nil {
			return err
		}
	} else {
		if err := db.Select("uid", "post_id").Create(&log).Error; err != nil {
			return err
		}
	}
	// 更新楼的likes
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err := db.Model(&post).Update("favs", post.Favs+1).Error; err != nil {
		return err
	}
	return nil
}

func UnfavPost(postId string, uid string) error {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	postIdint, _ := strconv.ParseUint(postId, 10, 64)

	var exist = false
	var log = LogPostFav{Uid: uidint, PostId: postIdint}
	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	if log.Id == 0 {
		return fmt.Errorf("未收藏")
	}

	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	exist = log.Id > 0
	if exist {
		if err := db.Delete(&log).Error; err != nil {
			return err
		}
	}

	// 更新楼的likes
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err := db.Model(&post).Update("favs", post.Favs-1).Error; err != nil {
		return err
	}

	return nil
}

func GetFavPostsById(uid string) ([]Post, error) {
	var posts []Post
	if err := db.Where("uid = ?", uid).Find(&posts).Error; err != nil {
		return posts, err
	}
	return posts, nil
}

func (LogPostFav) TableName() string {
	return "log_post_fav"
}

func (Post) TableName() string {
	return "posts"
}
