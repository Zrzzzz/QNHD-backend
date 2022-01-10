package models

type Admin struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;defualt:null" json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func GetAdmins(maps interface{}) ([]Admin, error) {
	var admins []Admin
	if err := db.Where(maps).Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func AddAdmins(name string, password string) (uint64, error) {
	var admin = Admin{Name: name, Password: password}
	if err := db.Create(&admin).Error; err != nil {
		return 0, err
	}
	return admin.Id, nil
}

func EditAdmins(name string, password string) error {
	if err := db.Model(&Admin{}).Where("name = ?", name).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAdmins(name string) (uint64, error) {
	var admin = Admin{}
	if err := db.Where("name = ?", name).First(&admin).Error; err != nil {
		return 0, err
	}
	if err := db.Delete(&admin).Error; err != nil {
		return 0, err
	}
	return admin.Id, nil
}

func (*Admin) TableName() string {
	return "admins"
}
