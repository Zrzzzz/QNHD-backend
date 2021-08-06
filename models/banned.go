package models

import (
	"log"
)

type Banned struct {
	Model
	Uid uint64 `json:"uid"`
}

func GetBanned(maps interface{}) (bans []Banned) {
	db.Where(maps).Find(&bans)

	return
}

func AddBannedByUid(uid uint64) bool {
	db.Select("uid").Create(&Banned{Uid: uid})
	db.Model(&User{}).Where("uid = ?", uid).Update("status", 0)
	return true
}

func AddBannedByEmail(email string) bool {
	var user User
	db.Where("email = ?", email).First(&user)
	log.Println(user)
	return AddBannedByUid(user.Uid)
}

func DeleteBannedByUid(uid uint64) bool {
	db.Where("uid = ?", uid).Delete(&Banned{})
	return true
}

func DeleteBannedByEmail(email string) bool {
	var user User
	db.Where("email = ?", email).First(&user)
	log.Println(user)
	return AddBannedByUid(user.Uid)
}

func IfBannedByEmail(email string) bool {
	var user User
	db.Where("email = ?", email).Find(&user)
	return IfBannedByUid(user.Uid)
}

func IfBannedByUid(uid uint64) bool {
	var ban Banned
	var cnt int
	db.Model(&ban).Where("uid = ?", uid).Count(&cnt)
	db.Where("uid = ?", uid).Last(&ban)
	return cnt > 0 && ban.DeletedAt == ""
}
