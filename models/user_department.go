package models

type UserDepartment struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement;null;" json:"id"`
	Uid          uint64 `json:"uid"`
	DepartmentId uint64 `json:"department_id"`
}

// 同时如果已经有部门了会删除之前的
func AddUserToDepartment(uid, departmentId uint64) error {
	var ud UserDepartment
	if err := db.Where("uid = ? AND department_id = ?", uid, departmentId).Find(&ud).Error; err != nil {
		return err
	}
	if ud.Id > 0 {
		ud.DepartmentId = departmentId
		if err := db.Model(&ud).Updates(ud).Error; err != nil {
			return err
		}
	} else {
		if err := db.Create(&UserDepartment{
			Uid:          uid,
			DepartmentId: departmentId,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func IsUserInDepartment(uid, departmentId string) (bool, error) {
	var ret UserDepartment
	if err := db.Where("uid = ? AND department_id = ?", uid, departmentId).Find(&ret).Error; err != nil {
		return false, err
	}
	return ret.Id > 0, nil
}

func (UserDepartment) TableName() string {
	return "user_department"
}
