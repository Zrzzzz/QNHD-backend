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
	Super       int    `json:"super"`
	SchAdmin    int    `json:"sch_admin"`
	StuAdmin    int    `json:"stu_admin"`
	IsUser      int    `json:"user"`
	Active      int    `json:"active" gorm:"default:1"`
	CreatedAt   string `json:"created_at" gorm:"autoCreateTime;default:null;"`
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
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return false, fmt.Errorf("未注册用户")
		}
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
	if !b {
		return b, fmt.Errorf("未赋予权限")
	}
	return b, nil
}

func UserRightDemand(uid string) (bool, error) {
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		return false, err
	}
	return user.IsUser == 1, nil
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
	if err := db.Where("number like ? AND user = 1", "%"+name+"%").Scopes(util.Paginate(c)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsers(c *gin.Context, name string) ([]User, error) {
	var users []User
	if err := db.Where("number like ? AND super = 0", "%"+name+"%").Scopes(util.Paginate(c)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetManagers(c *gin.Context, name string) ([]User, error) {
	var users []User
	if err := db.Where("number like ? AND super = 0 AND user = 0", "%"+name+"%").Scopes(util.Paginate(c)).Find(&users).Error; err != nil {
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

func AddUser(number, password, phoneNumber string) (uint64, error) {
	var user = User{
		Number:      number,
		Password:    password,
		PhoneNumber: phoneNumber,
		IsUser:      0,
	}
	if err := db.Select("number", "password", "phone_number", "user").Create(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func EditUser(uid string, maps map[string]interface{}) error {
	if err := db.Model(&User{}).Where("id = ? AND id <> 1", uid).Updates(maps).Error; err != nil {
		return err
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
