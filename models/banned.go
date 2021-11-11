package models

import (
	"errors"

	"gorm.io/gorm"
)

type Banned struct {
	Model
	Uid    uint64 `json:"uid"`
	Doer   string `json:"doer"`
	Reason string `json:"reason"`
}

func GetBanned(maps interface{}) ([]Banned, error) {
	var bans []Banned
	if err := db.Where(maps).Find(&bans).Error; err != nil {
		return bans, err
	}
	return bans, nil
}

func AddBannedByUid(uid uint64, doer string, reason string) (uint64, error) {
	var ban = Banned{Uid: uid, Doer: doer, Reason: reason}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("uid", "reason").Create(&ban).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{}).Where("id = ?", uid).Update("active", 0).Error; err != nil {
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
