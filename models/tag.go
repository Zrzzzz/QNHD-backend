package models

import (
	"errors"
	"fmt"
	"math/rand"
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/enums/TagPointType"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/segment"
	"qnhd/pkg/util"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Id     uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Uid    uint64 `json:"-"`
	Name   string `json:"name"`
	Tokens string `json:"-"`
}

type LogTag struct {
	TagId     uint64            `json:"tag_id"`
	Point     TagPointType.Enum `json:"point"`
	CreatedAt string            `json:"created_at" gorm:"default:null;"`
}

type HotTagResult struct {
	TagId int    `json:"tag_id"`
	Point int    `json:"point"`
	Name  string `json:"name"`
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	if err := db.Where("name = ?", name).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return tag.Id > 0, nil
}

func GetTags(name string) ([]Tag, error) {
	var tags []Tag
	var d = db.Model(&Tag{})
	if name != "" {
		d = db.Select("p.*", "ts_rank(p.tokens, q) as score").
			Table("(?) as p, plainto_tsquery(?) as q", d, segment.Cut(name, " ")).
			Where("q @@ p.tokens").Order("score DESC")
	}
	if err := d.Order("id").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTag(tagId string) (Tag, error) {
	var tag Tag
	err := db.Where("id = ?", tagId).Find(&tag).Error
	return tag, err
}

func GetRecommendTag(lastId int) (HotTagResult, error) {
	var tag HotTagResult
	tags, err := GetHotTags(10)
	if err != nil {
		return tag, err
	}
	if len(tags) == 0 {
		var t Tag
		db.Last(&t)
		tag.TagId = int(t.Id)
		tag.Point = 0
		tag.Name = t.Name
		return tag, nil
	}
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(tags))
	for idx == int(lastId) {
		idx = rand.Intn(len(tags))
	}
	return tags[idx], nil
}

// 获取24小时内高赞tag
func GetHotTags(cnt int) ([]HotTagResult, error) {
	var results []HotTagResult
	logs := db.Model(&LogTag{}).Where("created_at > CURRENT_TIMESTAMP + '-1 day'")
	if err := db.Table("(?) as a", logs).
		Joins("JOIN qnhd.tag ON qnhd.tag.id = tag_id").
		Select("tag_id", "sum(point) as point", "name").
		Group("tag_id, name").
		Limit(cnt).
		Order("point desc").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func AddTag(name, uid string) (uint64, error) {
	var tag = Tag{Name: filter.Filter(name), Uid: util.AsUint(uid)}
	if err := db.Select("name", "uid").Create(&tag).Error; err != nil {
		return 0, err
	}
	if err := flushTagTokens(tag.Id, tag.Name); err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func DeleteTagAdmin(uid string, id uint64) (uint64, error) {
	var tag Tag
	var err error
	if err = db.Where("id = ?", id).Find(&tag).Error; err != nil {
		return 0, err
	}
	if tag.Id > 0 {
		err = deleteTag(id)
	}
	addManagerLogWithDetail(util.AsUint(uid), id, ManagerLogType.TAG_DELETE,
		fmt.Sprintf("name: %s, creator: %v", tag.Name, tag.Uid))
	return tag.Id, err
}

func DeleteTag(id uint64, uid string) (uint64, error) {
	var tag Tag
	var err error
	if err = db.Where("id = ? AND uid = ?", id, uid).Find(&tag).Error; err != nil {
		return 0, err
	}
	if tag.Id > 0 {
		err = deleteTag(id)
	}
	return tag.Id, err
}

func deleteTag(id uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除下面的关联帖子
		if err := tx.Where("tag_id = ?", id).Delete(&PostTag{}).Error; err != nil {
			return err
		}
		// 删除tag记录
		if err := tx.Where("tag_id = ?", id).Delete(&LogTag{}).Error; err != nil {
			return err
		}
		// 删除tag
		return tx.Where("id = ?", id).Delete(&Tag{}).Error
	})
}

func addTagLogInPost(postId uint64, point TagPointType.Enum) error {
	var pt PostTag
	if err := db.Where("post_id = ?", postId).Find(&pt).Error; err != nil {
		return err
	}
	if pt.TagId != 0 {
		addTagLog(pt.TagId, point)
	}
	return nil
}

// 增加Tag访问记录
func addTagLog(id uint64, point TagPointType.Enum) {
	var log = LogTag{TagId: id, Point: point}
	if err := db.Create(&log).Error; err != nil {
		logging.Error("add tag log error: %v", log)
	}
}

// 给tag加热度
func AddTagLog(uid string, id uint64, point int64) error {
	var log = LogTag{TagId: id, Point: TagPointType.Enum(point)}
	addManagerLogWithDetail(util.AsUint(uid), id, ManagerLogType.TAG_POINT_ADD, fmt.Sprintf("add: %d", point))
	return db.Create(&log).Error
}

// 清空tag热度
func ClearTagLog(uid string, id uint64) error {
	addManagerLog(util.AsUint(uid), id, ManagerLogType.TAG_POINT_CLEAR)
	return db.Where("tag_id = ?", id).Delete(&LogTag{}).Error
}
