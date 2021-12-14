package models

import (
	"errors"

	"gorm.io/gorm"
)

type Department struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement;default:null;" json:"id"`
	Name         string `json:"name"`
	Introduction string `json:"introduction"`
}

func GetDepartments(name string) ([]Department, error) {
	var departs []Department
	if err := db.Where("name LIKE ?", "%"+name+"%").Find(&departs).Error; err != nil {
		return nil, err
	}
	return departs, nil
}

func GetDepartment(id uint64) (Department, error) {
	var depart Department
	err := db.Where("id = ?", id).First(&depart).Error
	return depart, err
}

func GetDepartmentHasUser(uid uint64) (Department, error) {
	var depart Department
	err := db.Joins("JOIN user_department AS ud ON ud.department_id = departments.id AND ud.uid = ?", uid).First(&depart).Error
	return depart, err
}

func AddDepartment(maps map[string]interface{}) (uint64, error) {
	var depart = Department{
		Name:         maps["name"].(string),
		Introduction: maps["introduction"].(string),
	}
	err := db.Create(&depart).Error
	return depart.Id, err
}

func EditDepartment(id, introduction string) error {
	if err := db.Where("id = ?", id).Updates(map[string]interface{}{
		"introduction": introduction,
	}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteDepartment(id string) (uint64, error) {
	var depart Department
	if err := db.Where("id = ?", id).Delete(&depart).Error; err != nil {
		return 0, err
	}
	return depart.Id, nil
}

func ExistDepartmentByName(name string) (bool, error) {
	var depart Department
	if err := db.Where("name = ?", name).First(&depart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return depart.Id > 0, nil
}

func (Department) TableName() string {
	return "departments"
}
