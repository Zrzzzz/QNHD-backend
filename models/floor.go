package models

type Floor struct {
	Model
	Uid     uint64 `json:"uid"`
	PostId  uint64 `json:"post_id" gorm:"index"`
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

func AddFloor(maps map[string]interface{}) bool {
	db.Select("uid", "post_id", "content").Create(&Floor{
		Uid:     maps["uid"].(uint64),
		PostId:  maps["postId"].(uint64),
		Content: maps["content"].(string),
	})
	return true
}

func DeleteFloorByAdmin(id string) bool {
	db.Where("id = ?", id).Delete(&Floor{})
	return true
}

func DeleteFloorByUser(postId, uid, floorId string) bool {
	db.Where("post_id = ? AND uid = ? AND id = ?", postId, uid, floorId).Delete(&Floor{})
	return true
}

func DeleteFloorsInPost(postId string) bool {
	db.Where("post_id = ?", postId).Delete(&Floor{})
	return true
}

func (Floor) TableName() string {
	return "floors"
}
