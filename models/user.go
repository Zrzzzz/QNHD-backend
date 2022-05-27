package models

import (
	"errors"
	"fmt"
	"math"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	giterrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	Uid                  uint64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;default:null;"`
	Nickname             string `json:"nickname" gorm:"default:''"`
	Number               string `json:"-" gorm:"default:''"`
	Password             string `json:"-" gorm:"column:password;"`
	PhoneNumber          string `json:"phone_number"`
	IsSuper              bool   `json:"is_super" gorm:"default:false;column:super_admin"`
	IsSchAdmin           bool   `json:"is_sch_admin" gorm:"default:false;column:school_department_admin"`
	IsStuAdmin           bool   `json:"is_stu_admin" gorm:"default:false;column:student_admin"`
	IsStuDistributeAdmin bool   `json:"is_stu_dis_admin" gorm:"default:false;column:school_distribute_admin"`
	IsUser               bool   `json:"is_user" gorm:"default:false;"`
	Active               bool   `json:"active" gorm:"default:true"`
	CreatedAt            string `json:"-" gorm:"autoCreateTime;default:null;"`
}

type NewUserData struct {
	Nickname     string `json:"nickname"`
	Password     string `json:"password" gorm:"column:password;"`
	PhoneNumber  string `json:"phone_number"`
	IsSuper      bool   `json:"is_super"`
	IsSchAdmin   bool   `json:"is_sch_admin"`
	IsStuAdmin   bool   `json:"is_stu_admin"`
	DepartmentId int    `json:"department_id"`
}

type UserRight struct {
	Super              bool
	SchAdmin           bool
	SchDistributeAdmin bool
	StuAdmin           bool
}

func RequireRight(uid string, right UserRight) bool {
	// 检查权限
	user, _ := GetUser(map[string]interface{}{"id": uid})
	var b = false
	if right.Super {
		b = b || user.IsSuper
	}
	if right.SchAdmin {
		b = b || user.IsSchAdmin
	}
	if right.StuAdmin {
		b = b || user.IsStuAdmin
	}
	if right.SchDistributeAdmin {
		b = b || user.IsStuDistributeAdmin
	}
	return b
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
	if !user.IsSuper && !user.IsSchAdmin && !user.IsStuAdmin {
		return fmt.Errorf("无管理员权限")
	}
	return nil
}

func RequireUser(uid string) error {
	user, err := GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		return err
	}
	if !user.IsUser {
		return fmt.Errorf("非用户身份")
	}
	return nil
}

func ExistUser(nickname, number string) (uint64, error) {
	var user User
	if err := db.Where(User{Nickname: nickname, Number: number, IsUser: true}).Order("id").Find(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func GetCommonUsers(c *gin.Context, maps map[string]interface{}) ([]User, error) {
	var (
		users     []User
		d         = db.Where("is_user = true")
		isBlocked = maps["is_blocked"].(string)
		IsBanned  = maps["is_banned"].(string)
	)
	if isBlocked == "1" {
		var blocks []uint64
		if err := db.Model(&Blocked{}).Select("uid").Distinct("uid").Where("expired_at > ?", gorm.Expr("CURRENT_TIMESTAMP")).Find(&blocks).Error; err != nil {
			return users, err
		}
		d = d.Where("id IN (?)", blocks)
	}
	if IsBanned == "1" {
		d = d.Where("active = false")
	}
	if err := d.Scopes(util.Paginate(c)).Order("id").Find(&users).Error; err != nil {
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
	users := db.Model(&User{}).Where("nickname like ? AND is_super = false AND is_user = false", "%"+name+"%")
	if err := db.
		Table("(?) as a", users).
		Select("a.*", "qd.name", "qd.introduction").
		Joins("LEFT JOIN qnhd.user_department as ud ON a.id = ud.uid").
		Joins("LEFT JOIN qnhd.department as qd ON ud.department_id = qd.id").
		Scopes(util.Paginate(c)).
		Order("id").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetUsersInDepartment(departmentId uint64) ([]User, error) {
	var users []User
	ud := db.Model(&UserDepartment{}).Where("department_id = ?", departmentId)
	err := db.Table("(?) as a", ud).
		Select("u.*").
		Joins("JOIN qnhd.user as u ON a.uid = u.id").
		Order("id").
		Find(&users).
		Error
	return users, err
}

func GetUser(maps map[string]interface{}) (User, error) {
	var u User
	if err := db.Where(maps).First(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func AddUser(nickname, number, password, phoneNumber string, isUser bool) (uint64, error) {
	var user = User{
		Nickname:    nickname,
		Number:      number,
		Password:    password,
		PhoneNumber: phoneNumber,
		IsUser:      isUser,
	}
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.Uid, nil
}

func AddUsers(users []NewUserData) error {
	// 先检查是否合理
	var err error
	for i, u := range users {
		if !(u.IsSchAdmin && u.IsStuAdmin && u.IsSuper) {
			err = giterrors.Wrap(err, fmt.Sprintf("line %v: 权限分配不正确", i))
		}
	}
	if err != nil {
		return err
	}
	var newUsers []User
	for _, u := range users {
		new := User{
			Nickname:    u.Nickname,
			Password:    u.Password,
			PhoneNumber: u.PhoneNumber,
			IsSuper:     u.IsSuper,
			IsSchAdmin:  u.IsSchAdmin,
			IsStuAdmin:  u.IsSchAdmin,
			IsUser:      false,
		}
		newUsers = append(newUsers, new)
	}
	// 一次插入2个参数，只要少于65535就ok
	insertCount := 250
	for i := 0; i < int(math.Ceil(float64(len(newUsers))/float64(insertCount))); i++ {
		min := (i + 1) * insertCount
		if len(newUsers) < min {
			min = len(newUsers)
		}
		if e := db.Create(newUsers[i*insertCount : min]).Error; e != nil {
			err = giterrors.Wrap(err, e.Error())
		}
	}
	// 看是否需要创建部门
	for i, new := range newUsers {
		if users[i].DepartmentId > 0 {
			if e := AddUserToDepartment(new.Uid, uint64(users[i].DepartmentId)); e != nil {
				err = giterrors.Wrap(err, e.Error())
			}
		}
	}
	return err
}

// 修改用户属性
func EditUser(uid string, maps map[string]interface{}) error {
	return db.Model(&User{}).Where("id = ?", uid).Updates(maps).Error
}

// 修改密码，要求原密码
func EditUserPasswd(uid string, rawPasswd, newPasswd string) error {
	var user User
	if err := db.Where("id = ? AND password = ?", uid, rawPasswd).First(&user).Error; err != nil {
		return err
	}
	return db.Model(&user).Update("password", newPasswd).Error
}

// 删除用户
func DeleteUser(uid uint64) error {
	return db.Where("id = ?", uid).Delete(&User{}).Error
}
