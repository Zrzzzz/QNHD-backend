package models

import (
	"math"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogUnreadFloor struct {
	Uid       uint64         `json:"uid"`
	FloorId   uint64         `json:"floor_id"`
	IsRead    bool           `json:"is_read"`
	CreatedAt string         `json:"created_at" gorm:"default:null;"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type UnreadFloorResponse struct {
	Type    int            `json:"type"`
	IsRead  bool           `json:"is_read"`
	ToFloor *FloorResponse `json:"to_floor"`
	Post    PostResponse   `json:"post"`
	Floor   FloorResponse  `json:"floor"`
}

func GetUnreadFloors(c *gin.Context, uid string) ([]UnreadFloorResponse, error) {
	var (
		ret       = []UnreadFloorResponse{}
		logFloors []LogUnreadFloor
		floors    []Floor
		err       error
	)

	// 先筛选出未读记录
	if err = db.Model(&LogUnreadFloor{}).Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("created_at DESC").Find(&logFloors).Error; err != nil {
		return ret, err
	}
	for _, log := range logFloors {
		var floor Floor
		if err = db.Where("id = ?", log.FloorId).First(&floor).Error; err != nil {
			continue
		}
		floors = append(floors, floor)
	}

	// 对每个楼层分析
	for _, f := range floors {
		var r = UnreadFloorResponse{Floor: f.geneResponse(false, false)}
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
			tofr := tof.geneResponse(false, false)
			r.ToFloor = &tofr
		} else {
			r.Type = 0
		}
		// 搜索帖子
		p, e := GetPost(util.AsStrU(f.PostId))
		if e != nil {
			err = e
			break
		}
		r.Post = p.geneResponse(false)
		ret = append(ret, r)
	}
	if err != gorm.ErrRecordNotFound {
		return ret, err
	}
	return ret, nil
}

// 添加评论通知
func addUnreadFloor(floorId uint64, uids ...uint64) error {
	var logs []LogUnreadFloor
	for _, u := range uids {
		logs = append(logs, LogUnreadFloor{Uid: u, FloorId: floorId})
	}
	// 一次插入2个参数，只要少于65535就ok，经测试250效率较高
	insertCount := 250
	for i := 0; i < int(math.Ceil(float64(len(logs))/float64(insertCount))); i++ {
		min := (i + 1) * insertCount
		if len(logs) < min {
			min = len(logs)
		}
		db.Create(logs[i*insertCount : min])
	}
	return nil
}

// 已读评论
func ReadFloor(uid, floorId uint64) error {
	return db.Model(&LogUnreadFloor{}).
		Where("uid = ? AND floor_id = ?", uid, floorId).
		Update("is_read", true).Error
}

// 已读帖子内的评论
func ReadFloorInPost(uid, postId uint64) error {
	var floors []uint64
	if err := db.Model(&Floor{}).Select("id").Where("post_id = ?", postId).Find(&floors).Error; err != nil {
		return err
	}
	return db.Model(&LogUnreadFloor{}).
		Where("uid = ? AND floor_id IN (?)", uid, floors).
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
