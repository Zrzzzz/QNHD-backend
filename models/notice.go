package models

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
	if err := db.Find(&notices).Error; err != nil {
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
	if err := db.Select("Content").Create(&notice).Error; err != nil {
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
	if err := db.Where("id = ?", id).Delete(&notice).Error; err != nil {
		return 0, err
	}
	return notice.Id, nil
}

func (Notice) TableName() string {
	return "notices"
}
