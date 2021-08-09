package models

import "qnhd/pkg/setting"

type Admin struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;defualt:null" json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CheckAdmin(name string, password string) bool {
	if (name == setting.AdminName && password == setting.AdminPass) {
		return true
	}
	var admin Admin
	db.Select("id").Where(Admin{Name: name, Password: password}).First(&admin)
	return admin.Id > 0
}

func GetAdmins(maps interface{}) (admins []Admin) {
	db.Where(maps).Find(&admins)
	return
}

func AddAdmins(name string, password string) bool {
	db.Create(&Admin{Name: name, Password: password})
	return true
}

func EditAdmins(name string, password string) bool {
	db.Model(&Admin{Name: name}).Update("password", password)
	return true
}

func DeleteAdmins(name string) bool {
	db.Where("name = ?", name).Delete(&Admin{})
	return true
}

func (*Admin) TableName() string {
	return "admins"
}
