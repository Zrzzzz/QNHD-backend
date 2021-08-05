package models

type User struct {
	UID        uint64 `gorm:"primaryKey;autoIncrement;" json:"uid"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	RegisterAt string `json:"register_at" gorm:"autoCreateTime"`
}

func GetUsers(maps interface{}) (users []User) {
	db.Where(maps).Find(&users)
	return
}

func GetUsersAll(maps interface{}) (cnt int) {
	db.Model(&User{}).Where(maps).Count(&cnt)
	return
}

func AddUser(email string, password string) bool {
	db.Create(&User{
		Email:    email,
		Password: password,
	})
	return true
}

func ExistUserByEmail(email string) bool {
	var user User
	db.Select("UID").Where("email = ?", email).First(&user)
	return user.UID > 0
}

func (User) TableName() string {
	return "users"
}
