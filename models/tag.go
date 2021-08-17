package models

type Tag struct {
	Id   uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Name string `json:"name"`
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Where("name = ?", name).First(&tag)
	return tag.Id > 0
}

func GetTags(name string) (tags []Tag) {
	db.Where("name LIKE ?", "%"+name+"%").Find(&tags)
	return
}

func AddTags(name string) bool {
	db.Select("name").Create(&Tag{Name: name})
	return true
}

func DeleteTags(id uint64) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}

func (Tag) TableName() string {
	return "tags"
}
