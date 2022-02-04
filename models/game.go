package models

import (
	"errors"

	"gorm.io/gorm"
)

type Game struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

func GetNewestGame() (Game, error) {
	var game Game
	err := db.Last(&game).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return game, nil
	}
	return game, err
}

func AddNewGame(content string) error {
	return db.Select("content").Create(&Game{Content: content}).Error
}
