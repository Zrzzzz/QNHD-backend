package models

import (
	"errors"

	"gorm.io/gorm"
)

type Banned struct {
	Model
	Uid uint64 `json:"uid" `
}

func GetBanned(maps interface{}) ([]Banned, error) {
	var bans []Banned
	if err := db.Where(maps).Find(&bans).Error; err != nil {
		return bans, err
	}
	return bans, nil
}

func AddBannedByUid(uid uint64) (uint64, error) {
	var ban = Banned{Uid: uid}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("uid").Create(&ban).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{}).Where("uid = ?", uid).Update("status", 0).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return ban.Id, nil
}

func DeleteBannedByUid(uid uint64) (uint64, error) {
	var ban Banned
	if err := db.Where("uid = ?", uid).Delete(&ban).Error; err != nil {
		return 0, err
	}
	return ban.Id, nil
}

func IfBannedByEmail(email string) (bool, error) {
	var user User
	if err := db.Where("email = ?", email).Find(&user).Error; err != nil {
		return false, err
	}
	return IfBannedByUid(user.Uid)
}

func IfBannedByUid(uid uint64) (bool, error) {
	var ban Banned
	if err := db.Where("uid = ?", uid).Last(&ban).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
