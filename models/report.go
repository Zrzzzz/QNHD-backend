package models

type Report struct {
	Model
	Uid    uint64 `json:"uid"`
	PostId string `json:"post_id"`
	Reason string `json:"reason"`
}

func GetReports() ([]Report, error) {
	var reports []Report
	if err := db.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}
