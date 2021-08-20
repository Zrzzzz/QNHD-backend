package models

import (
	"errors"

	"gorm.io/gorm"
)

type Tag struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Name string `json:"name"`
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

func (Tag) TableName() string {
	return "tags"
}
