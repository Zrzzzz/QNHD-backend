package models

import (
	"errors"
	"time"

	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

// 禁言
type Blocked struct {
	Model
	Uid       uint64 `json:"uid"`
	Doer      string `json:"doer"`
	Reason    string `json:"reason"`
	ExpiredAt string `json:"expired_at"`
	LastTime  uint8  `json:"last_time"`
}

type BlockedDetail struct {
	Starttime string
	Overtime  string
	Remain    uint64
}

func GetBlocked(maps interface{}) ([]Blocked, error) {
	var blocked []Blocked
	if err := db.Unscoped().Where(maps).Order("id DESC").Find(&blocked).Error; err != nil {
		return nil, err
	}
	return blocked, nil
}

func AddBlockedByUid(uid uint64, doer string, reason string, last uint8) (uint64, error) {
	expired_at := time.Now().Add(time.Hour * 24 * time.Duration(last)).Format("2006-01-02 15:04:05")
	var blocked = Blocked{Uid: uid, Doer: doer, Reason: reason, ExpiredAt: expired_at, LastTime: last}
	if err := db.Select("Uid", "Doer", "Reason", "ExpiredAt", "LastTime").Create(&blocked).Error; err != nil {
		return 0, err
	}

	return blocked.Id, nil
}

func DeleteBlockedByUid(uid uint64) (uint64, error) {
	var blocked = Blocked{}
	if err := db.Where("uid = ?", uid).First(&blocked).Error; err != nil {
		return 0, err
	}
	if err := db.Delete(&blocked).Error; err != nil {
		return 0, err
	}
	return blocked.Id, nil
}

func IsBlockedByUid(uid uint64) bool {
	var block Blocked
	if err := db.Where("uid = ? AND expired_at < ?", uid, gorm.Expr("CURRENT_TIMESTAMP")).Last(&block).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return true
}

func IsBlockedByUidDetailed(uid uint64) (bool, *BlockedDetail, error) {
	var ban Blocked
	if err := db.Where("uid = ?", uid).Last(&ban).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	if ban.Uid > 0 {
		var nowtime, overtime carbon.Carbon
		nowtime = carbon.Now()
		overtime = carbon.Parse(ban.ExpiredAt, "Asia/Shanghai")

		remain := uint64(overtime.Timestamp() - nowtime.Timestamp())
		return true, &BlockedDetail{
			Starttime: ban.CreatedAt,
			Overtime:  ban.ExpiredAt,
			Remain:    remain,
		}, nil
	}
	return false, nil, nil
}
