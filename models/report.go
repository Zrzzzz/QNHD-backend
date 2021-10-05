package models

type Report struct {
	Model
	Uid     uint64 `json:"uid"`
	PostId  uint64 `json:"post_id"`
	FloorId uint64 `json:"floor_id"`
	Reason  string `json:"reason"`
}

func GetReports() ([]Report, error) {
	var reports []Report
	if err := db.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func (Report) TableName() string {
	return "reports"
}
