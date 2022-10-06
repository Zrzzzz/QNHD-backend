package models

import (
	"errors"
	"fmt"
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/enums/NoticeType"
	"qnhd/enums/UserLevelOperationType"
	"qnhd/pkg/util"
	"time"

	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

// 禁言
type Blocked struct {
	Model
	Uid       uint64 `json:"uid"`
	Doer      uint64 `json:"doer"`
	Reason    string `json:"reason"`
	ExpiredAt string `json:"expired_at"`
	LastTime  uint8  `json:"last_time"`
}

type BlockedDetail struct {
	Starttime string
	Overtime  string
	Remain    int
}

func GetBlocked(maps interface{}) ([]Blocked, error) {
	var blocked []Blocked
	if err := db.Where(maps).Order("id DESC").Find(&blocked).Error; err != nil {
		return nil, err
	}
	return blocked, nil
}

func AddBlockedByUid(uid uint64, doer uint64, reason string, last uint8) (uint64, error) {
	expired_at := time.Now().Add(time.Hour * 24 * time.Duration(last)).Format("2006-01-02 15:04:05")
	var blocked = Blocked{Uid: uid, Doer: doer, Reason: reason, ExpiredAt: expired_at, LastTime: last}
	if err := db.Select("Uid", "Doer", "Reason", "ExpiredAt", "LastTime").Create(&blocked).Error; err != nil {
		return 0, err
	}

	addNoticeWithTemplate(NoticeType.BEEN_BLOCKED, []uint64{uid}, []string{reason, fmt.Sprintf("%d", last)})
	addManagerLogWithDetail(doer, uid, ManagerLogType.USER_BLOCK, fmt.Sprintf("reason: %s, day: %d", reason, last))
	switch last {
	case 1:
		EditUserLevel(util.AsStrU(uid), UserLevelOperationType.BLOCKED_1)
	case 3:
		EditUserLevel(util.AsStrU(uid), UserLevelOperationType.BLOCKED_3)
	case 7:
		EditUserLevel(util.AsStrU(uid), UserLevelOperationType.BLOCKED_7)
	case 14:
		EditUserLevel(util.AsStrU(uid), UserLevelOperationType.BLOCKED_14)
	case 30:
		EditUserLevel(util.AsStrU(uid), UserLevelOperationType.BLOCKED_30)
	}

	return blocked.Id, nil
}

func DeleteBlockedByUid(doer uint64, uid uint64) (uint64, error) {
	var blocked = Blocked{}
	if err := db.Where("uid = ?", uid).Last(&blocked).Error; err != nil {
		return 0, err
	}
	if err := db.Where("uid = ?", uid).Delete(&Blocked{}).Error; err != nil {
		return 0, err
	}
	addManagerLog(doer, uid, ManagerLogType.USER_UNBLOCK)

	return blocked.Id, nil
}

func IsBlockedByUid(uid uint64) bool {
	var block Blocked
	if err := db.Where("uid = ? AND expired_at > ?", uid, gorm.Expr("CURRENT_TIMESTAMP")).Last(&block).Error; err != nil {
		return false
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

		remain := int(overtime.Timestamp() - nowtime.Timestamp())
		return remain > 0, &BlockedDetail{
			Starttime: ban.CreatedAt,
			Overtime:  ban.ExpiredAt,
			Remain:    remain,
		}, nil
	}
	return false, nil, nil
}
