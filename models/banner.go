package models

import (
	"errors"

	"gorm.io/gorm"
)

type Banner struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

func GetNewestBanner() (Banner, error) {
	var banner Banner
	err := db.Last(&banner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return banner, nil
	}
	return banner, err
}

func AddNewBanner(content string) error {
	return db.Select("content").Create(&Banner{Content: content}).Error
}
