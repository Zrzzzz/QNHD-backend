package models

import "gorm.io/gorm"

type ReportType int

const (
	ReportTypePost ReportType = iota
	ReportTypeFloor
)

type Report struct {
	Model
	Uid     uint64 `json:"uid"`
	Type    int    `json:"type"`
	PostId  uint64 `json:"post_id"`
	FloorId uint64 `json:"floor_id"`
	Reason  string `json:"reason"`
}

func GetReports(rType ReportType) ([]Report, error) {
	var reports []Report
	if err := db.Where("type = ?", rType).Find(&reports).Order("created_at DESC").Error; err != nil {
		return nil, err
	}
	return reports, nil
}

func AddReport(maps map[string]interface{}) error {
	var report = &Report{
		Uid:     maps["uid"].(uint64),
		Type:    maps["type"].(int),
		PostId:  maps["post_id"].(uint64),
		FloorId: maps["floor_id"].(uint64),
		Reason:  maps["reason"].(string),
	}
	err := db.Create(report).Error
	return err
}

func DeleteReport(id string) error {
	return deleteReports(nil, map[string]interface{}{"id": id})
}

// 删除举报
func deleteReports(tx *gorm.DB, maps map[string]interface{}) error {
	if tx == nil {
		tx = db
	}
	return tx.Where(maps).Delete(&Report{}).Error
}
