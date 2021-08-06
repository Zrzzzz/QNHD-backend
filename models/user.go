package models

type User struct {
	Uid          uint64 `gorm:"primaryKey;autoIncrement;default:null;" json:"uid"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	RegisteredAt string `json:"register_at" gorm:"autoCreateTime;default:null;"`
	Status       int8   `json:"status" gorm:"default:null;"`
}

func GetUsers(maps interface{}) (users []User) {
	db.Where(maps).Find(&users)
	return
}

func AddUser(email string, password string) bool {
	db.Create(&User{
		Email:    email,
		Password: password,
	})
	return true
}

func EditUser(email string, data interface{}) bool {
	db.Model(&User{}).Where("email = ?", email).Updates(data)
	return true
}

func DeleteUser(email string) bool {
	db.Model(&User{Email: email}).Update("status", 0)
	return true
}

func ExistUserByEmail(email string) bool {
	var user User
	db.Select("uid").Where("email = ?", email).First(&user)
	return user.Uid > 0
}

func ValidUser(email string, password string) bool {
	var user User
	db.Where("email = ?", email).First(&user)
	if user.Email == email {
		return user.Password == password
	}
	return false
}

func (User) TableName() string {
	return "users"
}
