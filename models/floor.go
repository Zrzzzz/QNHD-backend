package models

type Floor struct {
	Model
	Uid         uint64 `json:"uid"`
	PostId      uint64 `json:"post_id" `
	Content     string `json:"content"`
	Nickname    string `json:"nickname"`
	ReplyTo     uint64 `json:"reply_to" `
	ReplyToName string `json:"reply_to_name"`
}

var FLOOR_NAME = []string{
	"Angus", "Bertram", "Conrad", "Devin", "Emmanuel", "Fitzgerald", "Gregary", "Herbert", "Ingram", "Joyce", "Kelly", "Leo", "Morton", "Nathaniel", "Orville", "Payne", "Quintion", "Regan", "Sean", "Tracy", "Uriah", "Valentine", "Walker", "Xavier", "Yves", "Zachary",
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
	var post Post
	var nickname string
	uid := maps["uid"].(uint64)
	postId := maps["postId"].(uint64)
	// 先找到post主人
	db.First(&post, postId)

	if post.Uid == uid {
		nickname = "Owner"
	} else {
		// 还有可能已经发过言
		var floor Floor
		db.Where("uid = ? AND post_id = ?", uid, postId).First(&floor)
		if floor.Id > 0 {
			nickname = floor.Nickname
		} else {
			var cnt int
			// 除去owner
			db.Table("floors").Where("post_id = ? AND uid <> ?", postId, post.Uid).Select("count(distinct(uid))").Count(&cnt)
			nickname = FLOOR_NAME[cnt]
		}
	}
	db.Select("uid", "post_id", "content", "nickname").Create(&Floor{
		Uid:      uid,
		PostId:   postId,
		Content:  maps["content"].(string),
		Nickname: nickname,
	})
	return true
}

func ReplyFloor(maps map[string]interface{}) bool {
	var post Post
	var nickname string
	uid := maps["uid"].(uint64)
	postId := maps["postId"].(uint64)
	// 先找到post主人
	db.First(&post, postId)

	if post.Uid == uid {
		nickname = "Owner"
	} else {
		// 还有可能已经发过言
		var floor Floor
		db.Where("uid = ? AND post_id = ?", uid, postId).First(&floor)
		if floor.Id > 0 {
			nickname = floor.Nickname
		} else {
			var cnt int
			// 除去owner
			db.Table("floors").Where("post_id = ? AND uid <> ?", postId, post.Uid).Select("count(distinct(uid))").Count(&cnt)
			nickname = FLOOR_NAME[cnt]
		}
	}

	floorId := maps["replyToFloor"].(uint64)
	var floor Floor
	db.First(&floor, floorId)

	db.Select("uid", "post_id", "content", "nickname", "reply_to", "reply_to_name").Create(&Floor{
		Uid:         uid,
		PostId:      postId,
		Content:     maps["content"].(string),
		Nickname:    nickname,
		ReplyTo:     floor.Uid,
		ReplyToName: floor.Nickname,
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
