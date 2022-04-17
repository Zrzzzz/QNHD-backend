package models

import (
	"gorm.io/gorm"
)

type Notice struct {
	Model
	Sender  string `json:"sender"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Symbol  string `json:"symbol"`
}

func GetNoticeTemplates() ([]Notice, error) {
	var notices []Notice
	err := db.Order("id").Find(&notices).Error
	return notices, err
}

// 向所有用户添加通知
func AddNoticeToAllUsers(data map[string]interface{}) error {
	data["symbol"] = "public"
	id, err := AddNoticeTemplate(data)
	if err != nil {
		return err
	}
	// 对所有用户通知
	if err := addUnreadNoticeToAllUser(id, data["pub_at"].(string)); err != nil {
		return err
	}
	return nil
}

func AddNoticeTemplate(data map[string]interface{}) (uint64, error) {
	var notice = Notice{
		Sender:  data["sender"].(string),
		Title:   data["title"].(string),
		Content: data["content"].(string),
		Symbol:  data["symbol"].(string),
	}
	// 创建模板
	err := db.Create(&notice).Error
	return notice.Id, err
}

func EditNoticeTemplate(id uint64, data map[string]interface{}) error {
	if err := db.Where("id = ?", id).Updates(&Notice{
		Sender:  data["sender"].(string),
		Title:   data["title"].(string),
		Content: data["content"].(string),
	}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNoticeTemplate(id uint64) (uint64, error) {
	var notice Notice
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Where("id = ?", id).Delete(&notice).Error; err != nil {
			return err
		}
		return db.Where("notice_id = ?", id).Delete(&LogUnreadNotice{}).Error
	})
	return notice.Id, err
}
