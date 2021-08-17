package models

type Report struct {
	Model
	Uid    uint64 `json:"uid"`
	PostId string `json:"post_id"`
	Reason string `json:"reason"`
}

func GetReports() (reports []Report) {
	db.Find(&reports)
	return
}
