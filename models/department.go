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
	if err := db.Where("name LIKE ?", "%"+name+"%").Order("id DESC").Find(&departs).Error; err != nil {
		return nil, err
	}
	return departs, nil
}

func GetDepartment(id uint64) (Department, error) {
	var depart Department
	err := db.Where("id = ?", id).First(&depart).Error
	return depart, err
}

// 获取帖子所在部门
func GetDepartmentByPostId(id uint64) (Department, error) {
	// 首先判断是否存在部门
	var departId int
	var depart Department
	if err := db.Model(&Post{}).Select("department_id").Where("id = ?", id).Find(&departId).Error; err != nil {
		return depart, err
	}
	// 使用主键查询
	err := db.First(&depart, departId).Error
	return depart, err
}

// 获取用户所在部门
func GetDepartmentByUid(uid uint64) (Department, error) {
	var depart Department
	err := db.Joins("JOIN qnhd.user_department AS ud ON ud.department_id = qnhd.department.id AND ud.uid = ?", uid).First(&depart).Error
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
	if err := db.Model(&Department{}).Where("id = ?", id).Updates(map[string]interface{}{
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

// 是否有部门已存在此名字
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

// 是否为部门管理员
func IsDepartmentHasUser(uid, departmentId uint64) bool {
	var cnt int64
	if err := db.Model(&UserDepartment{}).Select("Count(*)").Where("uid = ? AND department_id = ?", uid, departmentId).Count(&cnt).Error; err != nil {
		return false
	}
	return cnt > 0
}
