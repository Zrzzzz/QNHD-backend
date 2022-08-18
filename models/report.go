package models

import (
	"fmt"
	"qnhd/enums/ReportType"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Report struct {
	Model
	Uid       uint64          `json:"uid"`
	Type      ReportType.Enum `json:"type"`
	PostId    uint64          `json:"post_id"`
	FloorId   uint64          `json:"floor_id"`
	Reason    string          `json:"reason"`
	Solved    bool            `json:"solved"`
	IsDeleted bool            `json:"is_deleted" gorm:"-"`
}

type PostReportResponse struct {
	P       PostResponse `json:"post"`
	Reports []Report     `json:"reports"`
}

type FloorReportResponse struct {
	F       FloorResponse `json:"floor"`
	Reports []Report      `json:"reports"`
}

func GetPostReports(c *gin.Context) ([]PostReportResponse, error) {
	var posts []Post
	var ret = []PostReportResponse{}
	ids := db.Model(&Report{}).Select("post_id", "count(*) as cnt").Where("type = ? AND solved = false", ReportType.POST).Group("post_id").Order("cnt DESC").Scopes(util.Paginate(c))
	d := db.Unscoped().Select("p.*").Table("(?) as a", ids).Joins("JOIN qnhd.post p ON p.id = a.post_id")

	if err := d.Find(&posts).Error; err != nil {
		return nil, err
	}
	for _, i := range posts {
		var r PostReportResponse
		r.P = i.geneResponse(true)
		r.Reports = getReports(ReportType.POST, i.Id)
		ret = append(ret, r)
	}
	return ret, nil
}

func GetFloorReports(c *gin.Context) ([]FloorReportResponse, error) {
	var floors []Floor
	var ret = []FloorReportResponse{}
	ids := db.Model(&Report{}).Select("floor_id", "count(*) as cnt").Where("type = ? AND solved = false", ReportType.FLOOR).Group("floor_id").Order("cnt DESC").Scopes(util.Paginate(c))
	d := db.Unscoped().Select("p.*").Table("(?) as a", ids).Joins("JOIN qnhd.floor p ON p.id = a.floor_id")

	if err := d.Find(&floors).Error; err != nil {
		return nil, err
	}
	for _, i := range floors {
		var r FloorReportResponse
		r.F = i.geneResponse(false, true)
		r.Reports = getReports(ReportType.FLOOR, i.Id)
		ret = append(ret, r)
	}
	return ret, nil
}

func getReports(rType ReportType.Enum, id uint64) (reports []Report) {
	if rType == ReportType.POST {
		db.Where("type = ? AND post_id = ? AND solved = false", rType, id).Order("created_at").Find(&reports)
	} else {
		db.Where("type = ? AND floor_id = ? AND solved = false", rType, id).Order("created_at").Find(&reports)
	}
	return
}

func AddReport(maps map[string]interface{}) error {
	uid := maps["uid"].(uint64)
	t := maps["type"].(int)
	postId := maps["post_id"].(uint64)
	floorId := maps["floor_id"].(uint64)
	var report Report
	db.Where("uid = ? AND type = ? AND post_id = ? AND floor_id = ?", uid, t, postId, floorId).Find(&report)
	if report.Id > 0 {
		return fmt.Errorf("不能多次举报哦")
	}
	report = Report{
		Uid:     uid,
		Type:    ReportType.Enum(t),
		PostId:  postId,
		FloorId: floorId,
		Reason:  maps["reason"].(string),
	}
	err := db.Create(&report).Error
	return err
}

func SolveReport(t string, id string) error {
	if t == "1" {
		return db.Model(&Report{}).Where("type = ? AND post_id = ?", t, id).Update("solved", true).Error
	} else if t == "2" {
		return db.Model(&Report{}).Where("type = ? AND floor_id = ?", t, id).Update("solved", true).Error
	} else {
		return fmt.Errorf("举报类型错误")
	}
}

// 删除举报
func deleteReports(tx *gorm.DB, query interface{}, args ...interface{}) error {
	if tx == nil {
		tx = db
	}
	return tx.Where(query, args...).Delete(&Report{}).Error
}

// 恢复举报
func recoverReports(tx *gorm.DB, query interface{}, args ...interface{}) error {
	if tx == nil {
		tx = db
	}
	return tx.Unscoped().Model(&Report{}).Where(query, args...).Update("deleted_at", gorm.Expr("NULL")).Error
}
