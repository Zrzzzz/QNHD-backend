package models

type Banner struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	URL       string `json:"url"`
	Ord       int    `json:"ord"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

func GetBanners() ([]Banner, error) {
	var banner []Banner
	err := db.Order("ord DESC").Find(&banner).Error
	return banner, err
}

func AddBanner(maps map[string]interface{}) error {
	return db.Create(&Banner{
		Name:  maps["name"].(string),
		Title: maps["title"].(string),
		Image: maps["image"].(string),
		URL:   maps["url"].(string),
		Ord:   0,
	}).Error
}

func UpdateBannerOrder(id uint64, order int) error {
	return db.Model(&Banner{}).Where("id = ?", id).Update("ord", order).Error
}

func DeleteBanner(id uint64) error {
	return db.Where("id = ?", id).Delete(&Banner{}).Error
}
