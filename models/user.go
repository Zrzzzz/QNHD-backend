package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Uid       uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;default:null;"`
	Number    string `json:"number"`
	Password  string `json:"-"`
	Super     int    `json:"super"`
	SchAdmin  int    `json:"sch_admin"`
	StuAdmin  int    `json:"stu_admin"`
	Active    int    `json:"active" gorm:"default:1"`
	CreatedAt string `json:"created_at" gorm:"autoCreateTime;default:null;"`
}
type UserRight struct {
	Super    bool
	SchAdmin bool
	StuAdmin bool
}

// demand uid has admin right that ur param is true
func AdminRightDemand(uid string, ur UserRight) (bool, error) {
	// 检查权限
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		return false, err
	}
	var b = false
	if ur.Super {
		b = b || user.Super == 1
	}
	if ur.SchAdmin {
		b = b || user.SchAdmin == 1
	}
	if ur.StuAdmin {
		b = b || user.StuAdmin == 1
	}
	return b, nil
}

func UserRightDemand(uid string) (bool, error) {
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		return false, err
	}
	return user.Super == 0 && user.StuAdmin == 0 && user.SchAdmin == 0, nil
}

func CheckUser(number string, password string) (uint64, error) {
	var user User
	if err := db.Where(User{Number: number, Password: password}).First(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func ExistUser(number string) (uint64, error) {
	var user User
	if err := db.Where(User{Number: number}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return user.Uid, nil
}

func GetCommonUsers(uid string, overnum, pageSize int) ([]User, error) {
	var users []User
	if err := db.Where("id like ? AND super = 0 AND sch_admin = 0 AND stu_admin = 0", "%"+uid+"%").Offset(overnum).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsers(uid string, overnum, pageSize int) ([]User, error) {
	var users []User
	if err := db.Where("id like ? AND Super <> 1", "%"+uid+"%").Offset(overnum).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(maps map[string]interface{}) (User, error) {
	var u User
	if err := db.Where(maps).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func AddUser(number, password string) (uint64, error) {
	var user = User{
		Number:   number,
		Password: password,
	}
	if err := db.Select("number", "password").Create(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func EditUser(uid string, maps map[string]interface{}) error {
	if err := db.Model(&User{}).Where("id = ? AND (super = 1 OR stu_admin = 1 OR sch_admin = 1)", uid).Updates(maps).Error; err != nil {
		return err
	}
	return nil
}

func (User) TableName() string {
	return "roles"
}
