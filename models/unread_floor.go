package models

import (
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

type LogUnreadFloor struct {
	Uid       uint64 `json:"uid"`
	FloorId   uint64 `json:"floor_id"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

type UnreadFloorResponse struct {
	Type    int    `json:"type"`
	IsRead  bool   `json:"is_read"`
	ToFloor *Floor `json:"to_floor"`
	Post    Post   `json:"post"`
	Floor   Floor  `json:"floor"`
}

func GetUnreadFloors(c *gin.Context, uid string) ([]UnreadFloorResponse, error) {
	var (
		ret       = []UnreadFloorResponse{}
		logFloors []LogUnreadFloor
		floors    []Floor
		err       error
	)

	// 先筛选出未读记录
	logs := db.Model(&LogUnreadFloor{}).Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("created_at DESC")
	// 找到楼层
	if err = db.Table("(?) as a", logs).
		Unscoped().
		Select("f.*").
		Joins("JOIN qnhd.floor as f ON a.floor_id = f.id").
		Find(&floors).
		Where("f.deleted_at IS NULL").
		Error; err != nil {
		return ret, err
	}
	if err := logs.Find(&logFloors).Error; err != nil {
		return ret, err
	}
	// 对每个楼层分析
	for _, f := range floors {
		var r = UnreadFloorResponse{Floor: f}
		for _, log := range logFloors {
			if log.FloorId == f.Id {
				r.IsRead = log.IsRead
			}
		}
		// 搜索floor
		if f.SubTo > 0 {
			tof, e := GetFloor(util.AsStrU(f.ReplyTo))
			if e != nil {
				err = e
				break
			}
			r.Type = 1
			r.ToFloor = &tof
		} else {
			r.Type = 0
		}
		// 搜索帖子
		p, e := GetPost(util.AsStrU(f.PostId))
		if e != nil {
			err = e
			break
		}
		r.Post = p
		ret = append(ret, r)
	}

	return ret, err
}

// 添加评论通知
func addUnreadFloor(uid, floorId uint64) error {
	return db.Create(&LogUnreadFloor{
		Uid:     uid,
		FloorId: floorId,
	}).Error
}

// 已读评论
func ReadFloor(uid, floorId uint64) error {
	return db.Model(&LogUnreadFloor{}).
		Where("uid = ? AND floor_id = ?", uid, floorId).
		Update("is_read", true).Error
}

// 是否评论已读
func IsReadFloor(uid, floorId uint64) bool {
	var log LogUnreadFloor
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(log).Error; err != nil {
		return false
	}
	return log.IsRead
}
