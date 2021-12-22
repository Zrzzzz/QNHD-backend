package models

import (
	"errors"
	"fmt"
	"qnhd/pkg/logging"
	"qnhd/pkg/upload"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Floor struct {
	Model
	Uid         uint64 `json:"uid"`
	PostId      uint64 `json:"post_id"`
	Content     string `json:"content"`
	Nickname    string `json:"nickname"`
	ImageURL    string `json:"image_url"`
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

// var FLOOR_NAME = []string{
// 	"Angus", "Bertram", "Conrad", "Devin", "Emmanuel", "Fitzgerald", "Gregary", "Herbert", "Ingram", "Joyce", "Kelly", "Leo", "Morton", "Nathaniel", "Orville", "Payne", "Quintion", "Regan", "Sean", "Tracy", "Uriah", "Valentine", "Walker", "Xavier", "Yves", "Zachary",
// }
const FLOOR_NAME = "测试"

// 根据id返回
func GetFloor(floorId string) (Floor, error) {
	var floor Floor
	err := db.Where("id = ?", floorId).First(&floor).Error
	return floor, err
}

// 缩略返回帖子内楼层，即返回5条
func GetShortFloorsInPost(postId string) ([]Floor, error) {
	var floors []Floor
	if err := db.Where("post_id = ? AND reply_to IS NULL", postId).Order("created_at").Limit(5).Find(&floors).Error; err != nil {
		return nil, err
	}
	return floors, nil
}

// 分页返回帖子里的楼层
func GetFloorsInPost(c *gin.Context, postId string) ([]Floor, error) {
	var floors []Floor
	if err := db.Where("post_id = ? AND reply_to IS NULL", postId).Order("created_at").Scopes(util.Paginate(c)).Find(&floors).Error; err != nil {
		return nil, err
	}
	return floors, nil
}

// 缩略返回楼层内的回复，即返回5条
func GetFloorShortReplys(floorId string) ([]Floor, error) {
	var floors []Floor
	err := db.Where("reply_to = ?", floorId).Order("created_at").Limit(5).Find(&floors).Error
	return floors, err
}

// 分页返回楼层内的回复
func GetFloorReplys(c *gin.Context, floorId string) ([]Floor, error) {
	var floors []Floor
	err := db.Where("reply_to = ?", floorId).Order("created_at").Scopes(util.Paginate(c)).Find(&floors).Error
	return floors, err
}

func GetFloorByUid(uid string) ([]Floor, error) {
	var floors []Floor
	if err := db.Where("uid = ?", uid).Order("created_at").Find(&floors).Error; err != nil {
		return nil, err
	}
	return floors, nil
}

func AddFloor(maps map[string]interface{}) (uint64, error) {
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
			// nickname = FLOOR_NAME[cnt]
			nickname = fmt.Sprintf("%v%d", FLOOR_NAME, cnt)
		}
	}
	var newFloor = Floor{
		Uid:      uid,
		PostId:   postId,
		Content:  maps["content"].(string),
		Nickname: nickname,
		ImageURL: maps["image_url"].(string),
	}
	if err := db.Select("uid", "post_id", "content", "nickname", "image_url").Create(&newFloor).Error; err != nil {
		return 0, err
	}
	return newFloor.Id, nil
}

func ReplyFloor(maps map[string]interface{}) (uint64, error) {
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
			// nickname = FLOOR_NAME[cnt]
			nickname = fmt.Sprintf("%v%d", FLOOR_NAME, cnt)
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
		ImageURL:    maps["image_url"].(string),
		ReplyTo:     floor.Uid,
		ReplyToName: floor.Nickname,
	}
	if err := db.Select("uid", "post_id", "content", "nickname", "image_url", "reply_to", "reply_to_name").Create(&newFloor).Error; err != nil {
		return 0, err
	}

	return newFloor.Id, nil
}

func DeleteFloorByAdmin(id string) (uint64, error) {
	var floor = Floor{}
	if err := db.Where("id = ?", id).First(&floor).Error; err != nil {
		return 0, err
	}
	// 删除本地文件
	if err := upload.DeleteImageUrls([]string{floor.ImageURL}); err != nil {
		logging.Error(err.Error())
		return 0, err
	}
	if err := db.Delete(&floor).Error; err != nil {
		return 0, err
	}
	return floor.Id, nil
}

func DeleteFloorByUser(uid, floorId string) (uint64, error) {
	var floor = Floor{}
	if err := db.Where("uid = ? AND id = ?", uid, floorId).First(&floor).Error; err != nil {
		return 0, err
	}
	// 删除本地文件
	if err := upload.DeleteImageUrls([]string{floor.ImageURL}); err != nil {
		logging.Error(err.Error())
		return 0, err
	}
	if err := db.Delete(&floor).Error; err != nil {
		return 0, err
	}
	return floor.Id, nil
}

func DeleteFloorsInPost(tx *gorm.DB, postId string) error {
	if tx == nil {
		tx = db
	}
	var floors []Floor
	var imgs = []string{}
	if err := tx.Where("post_id = ?", postId).Find(&floors).Error; err != nil {
		return err
	}
	for _, f := range floors {
		if f.ImageURL != "" {
			imgs = append(imgs, f.ImageURL)
		}
	}
	// 删除本地文件
	if err := upload.DeleteImageUrls(imgs); err != nil {
		logging.Error(err.Error())
		return err
	}
	if err := tx.Where("post_id = ?", postId).Delete(&Floor{}).Error; err != nil {
		return err
	}
	return nil
}

/* 点赞或者取消点赞楼层 */
func LikeFloor(floorId string, uid string) (uint64, error) {
	uidint := util.AsUint(uid)
	floorIdint := util.AsUint(floorId)

	var exist = false
	var log = LogFloorLike{Uid: uidint, FloorId: floorIdint}

	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Id > 0 {
		return 0, fmt.Errorf("已被点赞")
	}

	if err := db.Unscoped().Where(log).Find(&log).Error; err != nil {
		return 0, err
	}

	exist = log.Id > 0
	if exist {
		if err := db.Unscoped().Model(&log).Update("deleted_at", gorm.Expr("NULL")).Error; err != nil {
			return 0, err
		}
	} else {
		if err := db.Select("uid", "floor_id").Create(&log).Error; err != nil {
			return 0, err
		}
	}
	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("like_count", floor.LikeCount+1).Error; err != nil {
		return 0, err
	}
	if _, err := UndisFloor(floorId, uid); err != nil {
		return 0, err
	}
	return floor.LikeCount + 1, nil
}

func UnlikeFloor(floorId string, uid string) (uint64, error) {
	uidint := util.AsUint(uid)
	floorIdint := util.AsUint(floorId)
	var exist = false
	var log = LogFloorLike{Uid: uidint, FloorId: floorIdint}
	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Id == 0 {
		return 0, fmt.Errorf("未被点赞")
	}

	if err := db.Where(log).Find(&log).Error; err != nil {
		return 0, err
	}
	exist = log.Id > 0
	if exist {
		if err := db.Delete(&log).Error; err != nil {
			return 0, err
		}
	}

	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("like_count", floor.LikeCount-1).Error; err != nil {
		return 0, err
	}

	return floor.LikeCount - 1, nil
}

/* 点赞或者取消点赞楼层 */
func DisFloor(floorId string, uid string) (uint64, error) {
	uidint := util.AsUint(uid)
	floorIdint := util.AsUint(floorId)

	var exist = false
	var log = LogFloorDis{Uid: uidint, FloorId: floorIdint}

	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Id > 0 {
		return 0, fmt.Errorf("已被点踩")
	}

	if err := db.Unscoped().Where(log).Find(&log).Error; err != nil {
		return 0, err
	}

	exist = log.Id > 0
	if exist {
		if err := db.Unscoped().Model(&log).Update("deleted_at", gorm.Expr("NULL")).Error; err != nil {
			return 0, err
		}
	} else {
		if err := db.Select("uid", "floor_id").Create(&log).Error; err != nil {
			return 0, err
		}
	}
	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("dis_count", floor.DisCount+1).Error; err != nil {
		return 0, err
	}
	if _, err := UnlikeFloor(floorId, uid); err != nil {
		return 0, err
	}
	return floor.DisCount + 1, nil
}

func UndisFloor(floorId string, uid string) (uint64, error) {
	uidint := util.AsUint(uid)
	floorIdint := util.AsUint(floorId)

	var exist = false
	var log = LogFloorDis{Uid: uidint, FloorId: floorIdint}
	// 首先判断点没点过赞
	if err := db.Where(log).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Id == 0 {
		return 0, fmt.Errorf("未被点踩")
	}

	if err := db.Where(log).Find(&log).Error; err != nil {
		return 0, err
	}
	exist = log.Id > 0
	if exist {
		if err := db.Delete(&log).Error; err != nil {
			return 0, err
		}
	}

	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("dis_count", floor.DisCount-1).Error; err != nil {
		return 0, err
	}

	return floor.DisCount - 1, nil
}

func IsLikeFloorByUid(uid, floorId string) bool {
	var log LogFloorLike
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Id > 0
}

func IsDisFloorByUid(uid, floorId string) bool {
	var log LogFloorDis
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Id > 0
}

func IsOwnFloorByUid(uid, floorId string) bool {
	var Floor, err = GetFloor(floorId)
	if err != nil {
		return false
	}
	return fmt.Sprintf("%d", Floor.Uid) == uid
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
