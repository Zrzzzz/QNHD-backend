package models

import (
	"qnhd/pkg/util"
)

type LogVisitHistory struct {
	Model
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}

func AddVisitHistory(uid string, postId string) (uint64, error) {

	var ps = LogVisitHistory{Uid: util.AsUint(uid), PostId: util.AsUint(postId)}
	if err := db.Select("post_id", "uid").Create(&ps).Error; err != nil {
		return 0, err
	}

	return ps.Id, nil
}

func (LogVisitHistory) TableName() string {
	return "log_visit_history"
}
