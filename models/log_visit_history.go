package models

import (
	"qnhd/enums/UserLevelOperationType"
	"qnhd/pkg/util"
)

type LogVisitHistory struct {
	Uid       uint64 `json:"uid"`
	PostId    uint64 `json:"post_id"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

func AddVisitHistory(uid string, postId string) error {
	var ps = LogVisitHistory{Uid: util.AsUint(uid), PostId: util.AsUint(postId)}
	if err := db.Select("post_id", "uid").Create(&ps).Error; err != nil {
		return err
	}
	EditUserLevel(uid, UserLevelOperationType.VISIT_POST)
	return nil
}
