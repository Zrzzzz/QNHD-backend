package models

type QnhdSetting struct {
	IOSLahei   bool `json:"ios_lahei"`
	FrontVisit bool `json:"front_visit"`
}

func (QnhdSetting) TableName() string {
	return "qnhd.setting"
}

func GetSetting() QnhdSetting {
	var q QnhdSetting
	db.First(&q)
	return q
}

func EditSetting(canVisit bool) error {
	return db.Model(&QnhdSetting{}).Where("true").Update("front_visit", canVisit).Error
}
