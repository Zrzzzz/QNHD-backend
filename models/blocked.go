package models

import (
	"time"

	"github.com/uniplaces/carbon"
)

type Blocked struct {
	Model
	Uid       uint64 `json:"uid" gorm:"index"`
	ExpiredAt string `json:"expired_at"`
	LastTime  uint8  `json:"last_time"`
}

type BlockedDetail struct {
	Starttime string
	Overtime  string
	Remain    uint64
}

func GetBlocked(maps interface{}) (blocked []Blocked) {
	db.Unscoped().Where(maps).Find(&blocked)

	return
}

func AddBlockedByUid(uid uint64, last uint8) bool {
	expired_at := time.Now().Add(time.Hour * 24 * time.Duration(last)).Format("2006-01-02 15:04:05")
	db.Select("Uid", "ExpiredAt", "LastTime").Create(&Blocked{Uid: uid, ExpiredAt: expired_at, LastTime: last})

	return true
}

func DeleteBlockedByUid(uid uint64) bool {
	db.Where("uid = ?", uid).Delete(&Blocked{})
	return true
}

func IfBlockedByUid(uid uint64) bool {
	var ban Blocked
	db.Where("uid = ?", uid).Last(&ban)

	return ban.Uid > 0
}

func IfBlockedByUidDetailed(uid uint64) (bool, *BlockedDetail) {
	var ban Blocked
	db.Where("uid = ?", uid).Last(&ban)
	if ban.Uid > 0 {
		var nowtime, overtime *carbon.Carbon
		nowtime = carbon.Now()
		overtime, _ = carbon.Parse(carbon.RFC3339Format, ban.ExpiredAt, "Asia/Shanghai")
		remain := uint64(overtime.Timestamp() - nowtime.Timestamp())
		return true, &BlockedDetail{
			Starttime: ban.CreatedAt,
			Overtime:  ban.ExpiredAt,
			Remain:    remain,
		}
	}
	return false, nil
}
