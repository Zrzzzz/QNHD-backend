package models

import (
	"errors"
	"math/rand"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Uid  uint64 `json:"-"`
	Name string `json:"name"`
}

type LogTag struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement;" json:"id"`
	TagId     uint64    `json:"tag_id"`
	Point     TAG_POINT `json:"point"`
	CreatedAt string    `json:"create_at"`
}

type HotTagResult struct {
	TagId int    `json:"tag_id"`
	Point int    `json:"point"`
	Name  string `json:"name"`
}

type TAG_POINT uint64

const (
	TAG_SEARCH TAG_POINT = iota + 1
	TAG_VISIT
	TAG_ADDFLOOR
	TAG_ADDPOST
)

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
	if err := db.Where("name LIKE ?", "%"+name+"%").Find(&tags).Error; err != nil {
		return nil, err
	}
	// 如果有name，对搜索到的加入记录，仅匹配精确搜索
	for _, t := range tags {
		if t.Name == "name" {
			addTagLog(t.Id, TAG_SEARCH)
		}
	}
	return tags, nil
}

func GetRecommendTag() (HotTagResult, error) {
	var tag HotTagResult
	tags, err := GetHotTags()
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
	return tags[rand.Intn(len(tags))], nil
}

// 获取24小时内高赞tag
func GetHotTags() ([]HotTagResult, error) {
	var results []HotTagResult
	if err := db.Model(&LogTag{}).
		Select("tag_id", "sum(point) as point", "name").
		Where("created_at > DATE_SUB(NOW(),INTERVAL 1 DAY)").
		Group("tag_id").
		Limit(5).
		Order("point desc").
		Joins("LEFT JOIN tags ON tags.id = tag_id").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func AddTag(name, uid string) (uint64, error) {
	var tag = Tag{Name: name, Uid: util.AsUint(uid)}
	if err := db.Select("name", "uid").Create(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func DeleteTagAdmin(id uint64) (uint64, error) {
	var tag Tag
	if err := db.Where("id = ?", id).Delete(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func DeleteTag(id uint64, uid string) (uint64, error) {
	var tag Tag
	if err := db.Where("id = ? AND uid = ?", id, uid).Delete(&tag).Error; err != nil {
		return 0, err
	}
	return tag.Id, nil
}

func addTagLogInPost(postId uint64, point TAG_POINT) error {
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
func addTagLog(id uint64, point TAG_POINT) {
	var log = LogTag{TagId: id, Point: point}
	if err := db.Select("tag_id").Create(&log).Error; err != nil {
		logging.Error("add tag log error: %v", log)
	}
}

func (LogTag) TableName() string {
	return "log_tag"
}

func (Tag) TableName() string {
	return "tags"
}
