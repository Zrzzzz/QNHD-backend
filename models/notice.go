package models

type Notice struct {
	Model
	Content string `json:"content"`
}

func GetNotices() ([]Notice, error) {
	var notices []Notice
	if err := db.Find(&notices).Error; err != nil {
		return nil, err
	}
	return notices, nil
}

func AddNotices(data map[string]interface{}) (uint64, error) {
	var notice = Notice{
		Content: data["content"].(string),
	}
	if err := db.Select("Content").Create(&notice).Error; err != nil {
		return 0, err
	}
	return notice.Id, nil
}

func EditNotices(id uint64, data map[string]interface{}) error {

	if err := db.Model(&Notice{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func DeleteNotices(id uint64) (uint64, error) {
	var notice = Notice{}
	if err := db.Where("id = ?", id).Delete(&notice).Error; err != nil {
		return 0, err
	}
	return notice.Id, nil
}

func (Notice) TableName() string {
	return "notices"
}
