package models

import (
	"qnhd/pkg/setting"
)

type Admin struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;defualt:null" json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CheckAdmin(name string, password string) (bool, error) {
	if name == setting.AppSetting.AdminName {
		return true, nil
	}
	var admin Admin
	if err := db.Select("id").Where(Admin{Name: name, Password: password}).First(&admin).Error; err != nil {
		return false, err
	}
	return true, nil
}

func ExistAdmin(name string) (bool, error) {
	if name == setting.AppSetting.AdminName {
		return true, nil
	}
	var admin Admin
	if err := db.Select("id").Where(Admin{Name: name}).First(&admin).Error; err != nil {
		return false, err
	}
	return true, nil
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

func EditAdmins(name string, password string) (bool, error) {
	if err := db.Model(&Admin{}).Where("name = ?", name).Update("password", password).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteAdmins(name string) (uint64, error) {
	var admin = Admin{}
	if err := db.Where("name = ?", name).Delete(&admin).Error; err != nil {
		return 0, err
	}
	return admin.Id, nil
}

func (*Admin) TableName() string {
	return "admins"
}
