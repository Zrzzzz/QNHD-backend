package models

func GetPostCount(from, to string) (int64, error) {
	var cnt int64
	err := db.Model(&Post{}).Where("created_at > ? AND created_at < ?", from, to).Count(&cnt).Error
	return cnt, err
}

func GetFloorCount(from, to string) (int64, error) {
	var cnt int64
	err := db.Model(&Floor{}).Where("created_at > ? AND created_at < ?", from, to).Count(&cnt).Error
	return cnt, err
}

func GetVisitPostCount(from, to string) (int64, error) {
	var cnt int64
	err := db.Model(&LogVisitHistory{}).Where("created_at > ? AND created_at < ?", from, to).Count(&cnt).Error
	return cnt, err
}
