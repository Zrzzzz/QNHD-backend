package models

import "strconv"

type VisitHistory struct {
	Model
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}

func GetVisitHistory(maps interface{}) ([]VisitHistory, error) {
	var bans []VisitHistory
	// TODO: 加入分页
	if err := db.Where(maps).Find(&bans).Error; err != nil {
		return bans, err
	}
	return bans, nil
}

func AddVisitHistory(uid string, postId string) (uint64, error) {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	pidint, _ := strconv.ParseUint(postId, 10, 64)
	var ps = VisitHistory{Uid: uidint, PostId: pidint}
	if err := db.Select("post_id", "uid").Create(&ps).Error; err != nil {
		return 0, err
	}

	return ps.Id, nil
}

func (VisitHistory) TableName() string {
	return "log_visit_history"
}
