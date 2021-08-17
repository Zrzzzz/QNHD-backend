package models

import "strconv"

type PostTag struct {
	Id     uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	PostId uint64 `json:"post_id" gorm:"index"`
	TagId  uint64 `json:"tag_id" gorm:"index"`
}

func GetTagsInPost(postId string) (tags []Tag) {
	db.Joins("JOIN post_tag ON tags.id = post_tag.tag_id").Where("post_id = ?", postId).Find(&tags)
	return
}

func AddPostWithTag(postId uint64, tags []string) {
	addDb := db.Select("PostId", "TagId")
	for _, t := range tags {
		intt, _ := strconv.ParseUint(t, 10, 64)
		addDb.Create(&PostTag{
			PostId: postId,
			TagId:  intt,
		})
	}
}

func DeleteTagInPost(postId string) {
	db.Where("post_id = ?", postId).Delete(&PostTag{})
}

func (PostTag) TableName() string {
	return "post_tag"
}
