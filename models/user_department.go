package models

type UserDepartment struct {
	Uid          uint64 `json:"uid"`
	DepartmentId uint64 `json:"department_id"`
}

// 同时如果已经有部门了会删除之前的
func AddUserToDepartment(uid, departmentId uint64) error {
	var ud UserDepartment
	if err := db.Where("uid = ?", uid).Find(&ud).Error; err != nil {
		return err
	}
	if ud.Uid > 0 {
		if err := db.Model(&ud).Where("uid = ?", ud.Uid).Update("department_id", departmentId).Error; err != nil {
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
