package models

type Blocked struct {
	Model
	Uid uint64 `json:"uid"`
}

func GetBlocked(maps interface{}) (bans []Blocked) {
	db.Where(maps).Find(&bans)

	return
}

func AddBlockedByUid(uid uint64) bool {
	db.Select("uid").Create(&Blocked{Uid: uid})

	return true
}

func AddBlockedByEmail(email string) bool {
	var user User
	db.Where("email = ?", email).First(&user)
	return AddBlockedByUid(user.Uid)
}

func DeleteBlockedByUid(uid uint64) bool {
	db.Where("uid = ?", uid).Delete(&Blocked{})
	return true
}

func DeleteBlockedByEmail(email string) bool {
	var user User
	db.Where("email = ?", email).First(&user)
	return AddBlockedByUid(user.Uid)
}

func IfBlockedByEmail(email string) bool {
	var user User
	db.Where("email = ?", email).Find(&user)
	return IfBlockedByUid(user.Uid)
}

func IfBlockedByUid(uid uint64) bool {
	var ban Blocked
	var cnt int
	db.Model(&ban).Where("uid = ?", uid).Count(&cnt)
	db.Where("uid = ?", uid).Last(&ban)
	return cnt > 0 && ban.DeletedAt == ""
}
