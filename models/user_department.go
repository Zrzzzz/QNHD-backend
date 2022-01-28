package models

type UserDepartment struct {
	Uid          uint64 `json:"uid"`
	DepartmentId uint64 `json:"department_id"`
}

// 同时如果已经有部门了会删除之前的
func AddUserToDepartment(uid, departmentId uint64) error {
	var ud UserDepartment
	if err := db.Where("uid = ? AND department_id = ?", uid, departmentId).Find(&ud).Error; err != nil {
		return err
	}
	if ud.Uid > 0 {
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

func (UserDepartment) TableName() string {
	return "user_department"
}
