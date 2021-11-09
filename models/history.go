package models

import "strconv"

type LogVisitHistory struct {
	Model
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}

func AddVisitHistory(uid string, postId string) (uint64, error) {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	pidint, _ := strconv.ParseUint(postId, 10, 64)
	var ps = LogVisitHistory{Uid: uidint, PostId: pidint}
	if err := db.Select("post_id", "uid").Create(&ps).Error; err != nil {
		return 0, err
	}

	return ps.Id, nil
}

func (LogVisitHistory) TableName() string {
	return "log_visit_history"
}
