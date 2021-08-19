package models

type Notice struct {
	Model
	Content string `json:"content"`
}

func GetNotices() (notices []Notice) {
	db.Find(&notices)
	return
}

func AddNotices(data map[string]interface{}) bool {
	db.Select("Content").Create(&Notice{
		Content: data["content"].(string),
	})
	return true
}

func EditNotices(id uint64, data map[string]interface{}) bool {
	db.Model(&Notice{}).Where("id = ?", id).Updates(data)
	return true
}

func DeleteNotices(id uint64) bool {
	db.Where("id = ?", id).Delete(&Notice{})
	return true
}

func (Notice) TableName() string {
	return "notices"
}
