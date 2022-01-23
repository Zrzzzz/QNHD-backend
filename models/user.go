package models

import (
	"errors"
	"fmt"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	Uid         uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;default:null;"`
	Number      string `json:"number"`
	Password    string `json:"-" gorm:"column:password;"`
	PhoneNumber string `json:"phone_number"`
	IsSuper     int    `json:"is_super"`
	IsSchAdmin  int    `json:"is_sch_admin"`
	IsStuAdmin  int    `json:"is_stu_admin"`
	IsUser      int    `json:"is_user"`
	Active      int    `json:"active" gorm:"default:1"`
	CreatedAt   string `json:"-" gorm:"autoCreateTime;default:null;"`
}

type UserRight struct {
	Super    bool
	SchAdmin bool
	StuAdmin bool
}

func RequireRight(uid string, right UserRight) error {
	// 检查权限
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		return err
	}
	var b = false
	if right.Super {
		b = b || user.IsSuper == 1
	}
	if right.SchAdmin {
		b = b || user.IsSchAdmin == 1
	}
	if right.StuAdmin {
		b = b || user.IsStuAdmin == 1
	}
	if !b {
		return fmt.Errorf("未赋予权限")
	}
	return nil
}

func RequireAdmin(uid string) error {
	// 检查权限
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("未注册用户")
		}
		return err
	}
	if user.IsUser != 0 {
		return fmt.Errorf("非管理员身份")
	}
	if user.IsSuper == 0 && user.IsSchAdmin == 0 && user.IsStuAdmin == 0 {
		return fmt.Errorf("无管理员权限")
	}
	return nil
}

func RequireUser(uid string) error {
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		return err
	}
	fmt.Println(user)
	if user.IsUser != 1 {
		return fmt.Errorf("非用户身份")
	}
	return nil
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

func GetCommonUsers(c *gin.Context, name string) ([]User, error) {
	var users []User
	if err := db.Where("number like ? AND is_user = 1", "%"+name+"%").Scopes(util.Paginate(c)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsers(c *gin.Context, name string) ([]User, error) {
	var users []User
	if err := db.Where("number like ? AND is_super = 0", "%"+name+"%").Scopes(util.Paginate(c)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type Manager struct {
	User
	Name         string `json:"department_name"`
	Introduction string `json:"department_introduction"`
}

func GetManagers(c *gin.Context, name string) ([]Manager, error) {
	var list []Manager
	users := db.Model(&User{}).Where("number like ? AND is_super = 0 AND is_user = 0", "%"+name+"%")
	if err := db.
		Table("(?) as a", users).
		Select("a.*, `departments`.`name`, `departments`.`introduction`").
		Joins("LEFT JOIN user_department ON a.id = user_department.uid").
		Joins("LEFT JOIN departments ON user_department.department_id = departments.id").
		Scopes(util.Paginate(c)).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetUser(maps map[string]interface{}) (User, error) {
	var u User
	if err := db.Where(maps).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func AddUser(number, password, phoneNumber string) (uint64, error) {
	var user = User{
		Number:      number,
		Password:    password,
		PhoneNumber: phoneNumber,
		IsUser:      0,
	}
	if err := db.Select("number", "password", "phone_number", "is_user").Create(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func EditUser(uid string, maps map[string]interface{}) error {
	return db.Model(&User{}).Where("id = ?", uid).Updates(maps).Error
}

func EditUserPasswd(uid string, rawPasswd, newPasswd string) error {
	var user User
	if err := db.Where("id = ? AND password = ?", uid, rawPasswd).First(&user).Error; err != nil {
		return err
	}
	return db.Model(&user).Update("password", newPasswd).Error
}

func (User) TableName() string {
	return "users"
}
