package models

import (
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EpiInfo struct {
	Infoid       uint64         `json:"infoid" gorm:"column:INFOID"`
	Infotime     string         `json:"infotime" gorm:"column:INFOTIME"`
	SystemId     uint64         `json:"system_id" gorm:"column:SYSTEM_ID"`
	ClassId      uint64         `json:"class_id" gorm:"column:CLASS_ID"`
	CategoryId   string         `json:"category_id" gorm:"column:CATEGORY_ID"`
	Author       string         `json:"author" gorm:"column:AUTHOR"`
	Title        string         `json:"title" gorm:"column:TITLE"`
	PreTitle     string         `json:"pre_title" gorm:"column:PRE_TITLE"`
	ViceTitle    string         `json:"vice_title" gorm:"column:VICE_TITLE"`
	Keyword      string         `json:"keyword" gorm:"column:KEYWORD"`
	Rank         uint64         `json:"rank" gorm:"column:RANK"`
	FpPic        string         `json:"fp_pic" gorm:"column:FP_PIC"`
	Url          string         `json:"url" gorm:"column:URL"`
	InfoFrom     string         `json:"info_from" gorm:"column:INFO_FROM"`
	InfoProvider string         `json:"info_provider" gorm:"column:INFO_PROVIDER"`
	ReadCnt      uint64         `json:"read_cnt" gorm:"column:READ_CNT"`
	AuditMark    string         `json:"audit_mark" gorm:"column:AUDIT_MARK"`
	AuditUser    string         `json:"audit_user" gorm:"column:AUDIT_USER"`
	AuditTime    string         `json:"audit_time" gorm:"column:AUDIT_TIME"`
	Content      string         `json:"content" gorm:"column:CONTENT"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

func GetEpiInfos(c *gin.Context) ([]EpiInfo, int, error) {
	var cnt int64
	var infos = []EpiInfo{}
	d := db.Model(&EpiInfo{})

	d.Count(&cnt)
	if err := d.Scopes(util.Paginate(c)).Order(clause.OrderByColumn{Column: clause.Column{Name: "RANK"}, Desc: true}).Order(clause.OrderByColumn{Column: clause.Column{Name: "INFOTIME"}, Desc: true}).Find(&infos).Error; err != nil {
		return nil, 0, err
	}

	return infos, int(cnt), nil
}

func AddEpiInfoReadCount(id uint64) (uint64, error) {
	var cnt int64
	d := db.Model(&EpiInfo{}).Where(map[string]interface{}{
		"INFOID": id,
	})
	if err := d.Select("READ_CNT").Find(&cnt).Error; err != nil {
		return 0, err
	}
	d.Update("READ_CNT", cnt+1)

	return uint64(cnt + 1), nil
}

func (EpiInfo) TableName() string {
	return "qnhd.outside_epidemic_info"
}
