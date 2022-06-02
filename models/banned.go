package models

import (
	"errors"
	ManagerLogType "qnhd/enums/MangerLogType"

	"gorm.io/gorm"
)

// 封号
type Banned struct {
	Model
	Uid    uint64 `json:"uid"`
	Doer   uint64 `json:"doer"`
	Reason string `json:"reason"`
}

func GetBanned(maps interface{}) ([]Banned, error) {
	var bans []Banned
	if err := db.Where(maps).Order("id DESC").Find(&bans).Error; err != nil {
		return bans, err
	}
	return bans, nil
}

func AddBannedByUid(uid uint64, doer uint64, reason string) (uint64, error) {
	var ban = Banned{Uid: uid, Doer: doer, Reason: reason}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ban).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{}).Where("id = ?", uid).Update("active", false).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	addManagerLog(doer, uid, ManagerLogType.USER_BAN)

	return ban.Id, nil
}

func DeleteBannedByUid(uid uint64) (uint64, error) {
	var ban Banned
	if err := db.Where("uid = ?", uid).First(&ban).Error; err != nil {
		return 0, err
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&ban).Error; err != nil {
			return err
		}
		if err := tx.Model(&User{}).Where("id = ?", uid).Update("active", true).Error; err != nil {
			return err
		}
		return nil
	})

	addManagerLog(ban.Doer, uid, ManagerLogType.USER_UNBAN)

	return ban.Id, err
}

func IsBannedByUid(uid uint64) bool {
	var ban Banned
	if err := db.Where("uid = ?", uid).Last(&ban).Error; err != nil {
		return !errors.Is(err, gorm.ErrRecordNotFound)
	}
	return true
}
