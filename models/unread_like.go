package models

import (
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

type LogUnreadLike struct {
	// 通知归属
	Uid       uint64   `json:"uid"`
	Type      LikeType `json:"type"`
	Id        uint64   `json:"id"`
	CreatedAt string   `json:"created_at" gorm:"default:null;"`
}

type UnreadLikeResponse struct {
	// 0为帖子 1位floor
	Type  int          `json:"type"`
	Post  PostResponse `json:"post"`
	Floor Floor        `json:"floor"`
}

type LikeType int

const (
	LIKE_POST LikeType = iota
	LIKE_FLOOR
)

func GetUnreadLikes(c *gin.Context, uid string) ([]UnreadLikeResponse, error) {
	var (
		ret  = []UnreadLikeResponse{}
		logs []LogUnreadLike
	)
	// 找到log
	if err := db.Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("created_at DESC").Find(&logs).Error; err != nil {
		return ret, err
	}
	// 逐个找floor
	for _, log := range logs {
		if log.Type == LIKE_POST {
			p, _ := GetPostResponse(util.AsStrU(log.Id))
			if p.Id > 0 {
				r := UnreadLikeResponse{Type: int(log.Type), Post: p}
				ret = append(ret, r)
			}
		} else if log.Type == LIKE_FLOOR {
			r := UnreadLikeResponse{Type: int(log.Type)}
			db.Where("id = ?", log.Id).Find(&r.Floor)
			if r.Floor.Id > 0 {
				ret = append(ret, r)
			}
		}
	}
	return ret, nil
}

func addUnreadLike(to uint64, likeType LikeType, id uint64) error {
	log := LogUnreadLike{Uid: to, Type: likeType, Id: id}
	return db.FirstOrCreate(&log, log).Error
}

func ReadLike(uid uint64, likeType LikeType, id uint64) error {
	return db.Where("uid = ? AND type = ? AND id = ?", uid, likeType, id).Delete(&LogUnreadLike{}).Error
}
