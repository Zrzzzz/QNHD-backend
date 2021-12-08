package models

import (
	"errors"

	"gorm.io/gorm"
)

type Department struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement;null;" json:"id"`
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

func AddDepartment(maps map[string]interface{}) (uint64, error) {
	var depart = Department{
		Name:         maps["name"].(string),
		Introduction: maps["introduction"].(string),
	}
	if err := db.Create(&depart).Error; err != nil {
		return 0, err
	}
	return depart.Id, nil
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
