package models

import (
	"errors"
	"qnhd/pkg/util"

	"gorm.io/gorm"
)

type Tag struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Uid  uint64 `json:"-"`
	Name string `json:"name"`
}

type LogTag struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	TagId     uint64 `json:"tag_id"`
	CreatedAt string `json:"create_at"`
}

type HotTagResult struct {
	TagId int    `json:"tag_id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	if err := db.Where("name = ?", name).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return tag.Id > 0, nil
}

func GetTags(name string) ([]Tag, error) {
	var tags []Tag
	if err := db.Where("name LIKE ?", "%"+name+"%").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// 获取24小时内高赞tag
func GetHotTags() ([]HotTagResult, error) {
	var results []HotTagResult
	if err := db.Model(&LogTag{}).
		Select("tag_id", "count(*) as count", "name").
		Where("created_at BETWEEN CONCAT(DATE_SUB(CURDATE(),INTERVAL 1 DAY), \" 08:00:00\") AND CONCAT(CURDATE(), \" 08:00:00\")").
		Group("tag_id").
		Limit(5).
		Order("count desc").
		Joins("LEFT JOIN tags ON tags.id = tag_id").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func AddTag(name, uid string) (uint64, error) {
	var tag = Tag{Name: name, Uid: util.AsUint(uid)}
	if err := db.Select("name", "uid").Create(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func DeleteTagAdmin(id uint64) (uint64, error) {
	var tag Tag
	if err := db.Where("id = ?", id).Delete(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func DeleteTag(id uint64, uid string) (uint64, error) {
	var tag Tag
	if err := db.Where("id = ? AND uid = ?", id, uid).Delete(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func AddTagLogInPost(postId uint64) error {
	var pt PostTag
	if err := db.Where("post_id = ?", postId).Find(&pt).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else {
			return nil
		}
	}
	return AddTagLog(pt.TagId)
}

// 增加Tag访问记录
func AddTagLog(id uint64) error {
	var tag = LogTag{TagId: id}
	if err := db.Select("tag_id").Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

// 删除24小时之前的记录
func FlushOldTagLog() error {
	if err := db.Where("HOUR(TIMEDIFF(NOW(), created_at)) >= 24;").Delete(&LogTag{}).Error; err != nil {
		return err
	}
	return nil
}

func (LogTag) TableName() string {
	return "log_tag"
}

func (Tag) TableName() string {
	return "tags"
}
