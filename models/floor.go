package models

type Floor struct {
	Model
	Uid     uint64 `json:"uid"`
	PostId  string `json:"post_id" gorm:"index"`
	Content string `json:"content"`
}

func GetFloorInPostShort(postId string) (floors []Floor) {
	db.Where("post_id = ?", postId).Order("created_at").Limit(10).Find(&floors)
	return
}

func GetFloorInPost(overNum int, pageSize int, postId string) (floors []Floor) {
	db.Where("post_id = ?", postId).Order("created_at").Offset(overNum).Limit(pageSize).Find(&floors)
	return
}

func GetFloorByUid(uid string) (floors []Floor) {
	db.Where("uid = ?", uid).Order("created_at").Find(&floors)
	return
}

func GetFloor(id string) (floor Floor) {
	db.Where("id = ?", id).First(&floor)
	return
}

func DeleteFloor(id string) bool {
	db.Where("id = ?", id).Delete(&Floor{})
	return true
}

func (Floor) TableName() string {
	return "floors"
}
