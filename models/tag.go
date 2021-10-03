package models

import (
	"errors"

	"gorm.io/gorm"
)

type Tag struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Name string `json:"name"`
}

type LogTag struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	TagId     uint64 `json:"tag_id"`
	CreatedAt string `json:"create_at"`
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
func GetHotTags() ([]Tag, error) {
	// var tags []Tag
	// var logTags []LogTag
	// if err := db.Where("created_at BETWEEN CONCAT(DATE_SUB(CURDATE(),INTERVAL 1 DAY), \" 08:00:00\") AND CONCAT(CURDATE(), \" 08:00:00\"").Group("")
	return []Tag{}, nil
}

func AddTags(name string) (uint64, error) {
	var tag = Tag{Name: name}
	if err := db.Select("name").Create(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func DeleteTags(id uint64) (uint64, error) {
	var tag = Tag{}
	if err := db.Where("id = ?", id).Delete(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
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
