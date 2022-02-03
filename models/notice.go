package models

import "gorm.io/gorm"

type Notice struct {
	Model
	Sender  string `json:"sender"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Url     string `json:"url"`
	Read    int    `json:"read" gorm:"-"`
}

func GetNotices() ([]Notice, error) {
	var notices []Notice
	if err := db.Find(&notices).Order("id").Error; err != nil {
		return nil, err
	}
	return notices, nil
}

func AddNotice(data map[string]interface{}) (uint64, error) {
	var notice = Notice{
		Sender:  data["sender"].(string),
		Title:   data["title"].(string),
		Content: data["content"].(string),
		Url:     data["url"].(string),
	}
	// 创建通知
	if err := db.Create(&notice).Error; err != nil {
		return 0, err
	}
	// 对所有用户通知
	if err := addUnreadNoticeToAllUser(notice.Id); err != nil {
		return 0, err
	}
	return notice.Id, nil
}

func EditNotice(id uint64, data map[string]interface{}) error {
	if err := db.Model(&Notice{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func DeleteNotice(id uint64) (uint64, error) {
	var notice Notice
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Where("id = ?", id).Delete(&notice).Error; err != nil {
			return err
		}
		return db.Where("notice_id = ?", id).Delete(&LogUnreadNotice{}).Error
	})
	return notice.Id, err
}
