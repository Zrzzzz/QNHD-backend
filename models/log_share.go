package models

import (
	"qnhd/enums/ShareLogType"
	"qnhd/enums/UserLevelOperationType"
	"qnhd/pkg/util"
)

type LogShare struct {
	Uid       uint64            `json:"uid"`
	ObjectId  uint64            `json:"post_id"`
	Type      ShareLogType.Enum `json:"type"`
	CreatedAt string            `json:"created_at" gorm:"default:null;"`
}

func AddShareLog(uid string, objectId uint64, t ShareLogType.Enum) error {
	if err := db.Create(&LogShare{
		Uid:      util.AsUint(uid),
		ObjectId: objectId,
		Type:     t,
	}).Error; err != nil {
		return err
	}
	EditUserLevel(uid, UserLevelOperationType.SHARE_POST)
	return nil
}
