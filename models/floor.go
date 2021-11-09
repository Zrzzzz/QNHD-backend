package models

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Floor struct {
	Model
	Uid         uint64 `json:"uid"`
	PostId      uint64 `json:"post_id"`
	Content     string `json:"content"`
	Nickname    string `json:"nickname"`
	ReplyTo     uint64 `json:"reply_to" `
	ReplyToName string `json:"reply_to_name"`
	LikeCount   uint64 `json:"like_count"`
	DisCount    uint64 `json:"-"`
}

type LogFloorLike struct {
	Model
	Uid     uint64 `json:"uid"`
	FloorId uint64 `json:"floor_id"`
}

type LogFloorDis struct {
	Model
	Uid     uint64 `json:"uid"`
	FloorId uint64 `json:"floor_id"`
}

const OWNER_NAME = "Owner"

var FLOOR_NAME = []string{
	"Angus", "Bertram", "Conrad", "Devin", "Emmanuel", "Fitzgerald", "Gregary", "Herbert", "Ingram", "Joyce", "Kelly", "Leo", "Morton", "Nathaniel", "Orville", "Payne", "Quintion", "Regan", "Sean", "Tracy", "Uriah", "Valentine", "Walker", "Xavier", "Yves", "Zachary",
}

func GetFloorInPostShort(postId string) ([]Floor, error) {
	var floors []Floor
	if err := db.Where("post_id = ?", postId).Order("created_at").Limit(5).Find(&floors).Error; err != nil {
		return nil, err
	}
	return floors, nil
}

func GetFloorInPost(base int, pageSize int, postId string) ([]Floor, error) {
	var floors []Floor
	if err := db.Where("post_id = ?", postId).Order("created_at").Offset(base).Limit(pageSize).Find(&floors).Error; err != nil {
		return nil, err
	}
	return floors, nil
}

func GetFloorByUid(uid string) ([]Floor, error) {
	var floors []Floor
	if err := db.Where("uid = ?", uid).Order("created_at").Find(&floors).Error; err != nil {
		return nil, err
	}
	return floors, nil
}

func GetFloor(id string) (Floor, error) {
	var floor Floor
	if err := db.Where("id = ?", id).First(&floor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return floor, nil
		}
		return floor, err
	}
	return floor, nil
}

func AddFloor(maps map[string]interface{}) (uint64, error) {
	// TODO: 添加锁
	var post Post
	var nickname string
	uid := maps["uid"].(uint64)
	postId := maps["postId"].(uint64)
	// 先找到post主人
	if err := db.First(&post, postId).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}

	if post.Uid == uid {
		nickname = OWNER_NAME
	} else {
		// 还有可能已经发过言
		var floor Floor
		if err := db.Where("uid = ? AND post_id = ?", uid, postId).First(&floor).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, err
			}
		}
		if floor.Id > 0 {
			nickname = floor.Nickname
		} else {
			var cnt int64
			// 除去owner
			if err := db.Table("floors").Where("post_id = ? AND uid <> ?", postId, post.Uid).Distinct("uid").Count(&cnt).Error; err != nil {
				return 0, err
			}
			nickname = FLOOR_NAME[cnt]
		}
	}
	var newFloor = Floor{
		Uid:      uid,
		PostId:   postId,
		Content:  maps["content"].(string),
		Nickname: nickname,
	}
	if err := db.Select("uid", "post_id", "content", "nickname").Create(&newFloor).Error; err != nil {
		return 0, err
	}
	return newFloor.Id, nil
}

func ReplyFloor(maps map[string]interface{}) (uint64, error) {
	// TODO: 添加锁
	var post Post
	var nickname string
	uid := maps["uid"].(uint64)
	postId := maps["postId"].(uint64)
	// 先找到post主人
	if err := db.First(&post, postId).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}

	if post.Uid == uid {
		nickname = OWNER_NAME
	} else {
		// 还有可能已经发过言
		var floor Floor
		if err := db.Where("uid = ? AND post_id = ?", uid, postId).First(&floor).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, err
			}
		}
		if floor.Id > 0 {
			nickname = floor.Nickname
		} else {
			var cnt int64
			// 除去owner
			if err := db.Table("floors").Where("post_id = ? AND uid <> ?", postId, post.Uid).Distinct("uid").Count(&cnt).Error; err != nil {
				return 0, err
			}
			nickname = FLOOR_NAME[cnt]
		}
	}

	floorId := maps["replyToFloor"].(uint64)
	var floor Floor
	if err := db.First(&floor, floorId).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}

	var newFloor = Floor{
		Uid:         uid,
		PostId:      postId,
		Content:     maps["content"].(string),
		Nickname:    nickname,
		ReplyTo:     floor.Uid,
		ReplyToName: floor.Nickname,
	}
	if err := db.Select("uid", "post_id", "content", "nickname", "reply_to", "reply_to_name").Create(&newFloor).Error; err != nil {
		return 0, err
	}

	return newFloor.Id, nil
}

func DeleteFloorByAdmin(id string) (uint64, error) {
	var floor = Floor{}
	if err := db.Where("id = ?", id).Delete(&floor).Error; err != nil {
		return 0, err
	}
	return floor.Id, nil
}

func DeleteFloorByUser(postId, uid, floorId string) (uint64, error) {
	var floor = Floor{}
	if err := db.Where("post_id = ? AND uid = ? AND id = ?", postId, uid, floorId).Delete(&floor).Error; err != nil {
		return 0, err
	}
	return floor.Id, nil
}

func DeleteFloorsInPost(postId string) (uint64, error) {
	var floor Floor
	if err := db.Where("post_id = ?", postId).Delete(&floor).Error; err != nil {
		return 0, err
	}
	return floor.Id, nil
}

/* 点赞或者取消点赞楼层 */
func LikeFloor(floorId string, uid string) error {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	floorIdint, _ := strconv.ParseUint(floorId, 10, 64)

	var exist = false
	var log = LogFloorLike{Uid: uidint, FloorId: floorIdint}

	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	if log.Id > 0 {
		return fmt.Errorf("已被点赞")
	}

	if err := db.Unscoped().Where(log).Find(&log).Error; err != nil {
		return err
	}

	exist = log.Id > 0
	if exist {
		if err := db.Unscoped().Model(&log).Update("deleted_at", gorm.Expr("NULL")).Error; err != nil {
			return err
		}
	} else {
		if err := db.Select("uid", "floor_id").Create(&log).Error; err != nil {
			return err
		}
	}
	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err := db.Model(&floor).Update("likes", floor.LikeCount+1).Error; err != nil {
		return err
	}

	return nil
}

func UnlikeFloor(floorId string, uid string) error {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	floorIdint, _ := strconv.ParseUint(floorId, 10, 64)

	var exist = false
	var log = LogFloorLike{Uid: uidint, FloorId: floorIdint}
	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	if log.Id == 0 {
		return fmt.Errorf("未被点赞")
	}

	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	exist = log.Id > 0
	if exist {
		if err := db.Delete(&log).Error; err != nil {
			return err
		}
	}

	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err := db.Model(&floor).Update("likes", floor.LikeCount-1).Error; err != nil {
		return err
	}

	return nil
}

/* 点赞或者取消点赞楼层 */
func DisFloor(floorId string, uid string) error {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	floorIdint, _ := strconv.ParseUint(floorId, 10, 64)

	var exist = false
	var log = LogFloorDis{Uid: uidint, FloorId: floorIdint}

	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	if log.Id > 0 {
		return fmt.Errorf("已被点踩")
	}

	if err := db.Unscoped().Where(log).Find(&log).Error; err != nil {
		return err
	}

	exist = log.Id > 0
	if exist {
		if err := db.Unscoped().Model(&log).Update("deleted_at", gorm.Expr("NULL")).Error; err != nil {
			return err
		}
	} else {
		if err := db.Select("uid", "floor_id").Create(&log).Error; err != nil {
			return err
		}
	}
	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err := db.Model(&floor).Update("dis_count", floor.DisCount+1).Error; err != nil {
		return err
	}

	return nil
}

func UndisFloor(floorId string, uid string) error {
	uidint, _ := strconv.ParseUint(uid, 10, 64)
	floorIdint, _ := strconv.ParseUint(floorId, 10, 64)

	var exist = false
	var log = LogFloorDis{Uid: uidint, FloorId: floorIdint}
	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	if log.Id == 0 {
		return fmt.Errorf("未被点踩")
	}

	if err := db.Where(log).Find(&log).Error; err != nil {
		return err
	}
	exist = log.Id > 0
	if exist {
		if err := db.Delete(&log).Error; err != nil {
			return err
		}
	}

	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if err := db.Model(&floor).Update("dis_count", floor.DisCount-1).Error; err != nil {
		return err
	}

	return nil
}

func (LogFloorLike) TableName() string {
	return "log_floor_like"
}
func (LogFloorDis) TableName() string {
	return "log_floor_dis"
}

func (Floor) TableName() string {
	return "floors"
}
