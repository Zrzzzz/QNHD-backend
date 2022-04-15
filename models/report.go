package models

import (
	"fmt"

	"gorm.io/gorm"
)

type ReportType int

const (
	ReportTypePost ReportType = iota
	ReportTypeFloor
)

type Report struct {
	Model
	Uid       uint64 `json:"uid"`
	Type      int    `json:"type"`
	PostId    uint64 `json:"post_id"`
	FloorId   uint64 `json:"floor_id"`
	Reason    string `json:"reason"`
	Solved    bool   `json:"solved" gorm:"-"`
	IsDeleted bool   `json:"is_deleted" gorm:"-"`
}

func GetReports(rType ReportType) ([]Report, error) {
	var reports []Report
	if err := db.Unscoped().Where("type = ? AND solved = false", rType).Order("created_at DESC").Find(&reports).Error; err != nil {
		return nil, err
	}
	for i := range reports {
		reports[i].IsDeleted = reports[i].DeletedAt.Valid
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

func SolveReports(t string, id string) error {
	if t == "1" {
		return db.Model(&Report{}).Where("type = ? AND post_id = ?", t, id).Update("solved", true).Error
	} else if t == "2" {
		return db.Model(&Report{}).Where("type = ? AND floor_id = ?", t, id).Update("solved", true).Error
	} else {
		return fmt.Errorf("举报类型错误")
	}
}

// 删除举报
func deleteReports(tx *gorm.DB, maps map[string]interface{}) error {
	if tx == nil {
		tx = db
	}
	return tx.Where(maps).Delete(&Report{}).Error
}

// 恢复举报
func recoverReports(tx *gorm.DB, maps map[string]interface{}) error {
	if tx == nil {
		tx = db
	}
	return tx.Unscoped().Model(&Report{}).Where(maps).Update("deleted_at", gorm.Expr("NULL")).Error
}
