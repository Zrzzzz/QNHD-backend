package models

type PostType struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement;" json:"id"`
	Shortname string `json:"shortname"`
	Name      string `json:"name"`
	Ord       int    `json:"-" gorm:"default:null"`
	Hidden    bool   `json:"hidden" gorm:"default:false"`
}

func IsValidPostType(t int) bool {
	types, err := GetPostTypes()
	if err != nil {
		return false
	}
	for _, ty := range types {
		if t == int(ty.Id) {
			return true
		}
	}
	return false
}

func GetPostTypes() ([]PostType, error) {
	var ret []PostType
	err := db.Where("hidden = false").Order("ord DESC").Order("id").Find(&ret).Error
	return ret, err
}

func AddPostType(short, name string) error {
	return db.Create(&PostType{Shortname: short, Name: name}).Error
}
