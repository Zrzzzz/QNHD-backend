package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	Uid          uint64 `gorm:"primaryKey;autoIncrement;default:null;" json:"uid"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	RegisteredAt string `json:"register_at" gorm:"autoCreateTime;default:null;"`
	Status       int8   `json:"status" gorm:"default:null;"`
}

func CheckUser(email string, password string) (uint64, error) {
	var user User
	if err := db.Select("uid").Where(User{Email: email, Password: password}).First(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func ExistUser(email string) (uint64, error) {
	var user User
	if err := db.Where(User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return user.Uid, nil
}

func GetUsers(maps interface{}) ([]User, error) {
	var users []User
	if err := db.Where(maps).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func AddUser(email string, password string) (uint64, error) {
	var user = User{
		Email:    email,
		Password: password,
	}
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func EditUser(email string, data interface{}) error {
	if err := db.Model(&User{}).Where("email = ?", email).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(email string) error {
	if err := db.Model(&User{Email: email}).Update("status", 0).Error; err != nil {
		return err
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
