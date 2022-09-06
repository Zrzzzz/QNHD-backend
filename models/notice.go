package models

import (
	"fmt"
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Notice struct {
	Model
	Sender  string `json:"sender"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Symbol  string `json:"symbol"`
}

const NOTICE_DEPARTMENT = "department_manager"

func GetNotices(c *gin.Context, departmentOnly bool) ([]Notice, error) {
	var notices []Notice
	d := db
	if departmentOnly {
		d = d.Where("symbol = ?", NOTICE_DEPARTMENT)
	} else {
		d = d.Where("symbol = 'public'")
	}
	err := d.Scopes(util.Paginate(c)).Order("created_at DESC").Find(&notices).Error
	return notices, err
}

// 向所有用户添加通知
func AddNoticeToAllUsers(uid string, data map[string]interface{}) error {
	data["symbol"] = "public"
	var user User
	db.Where("id = ?", uid).Find(&user)
	if user.IsSchAdmin {
		data["symbol"] = NOTICE_DEPARTMENT
	}
	id, err := AddNoticeTemplate(data)
	if err != nil {
		return err
	}

	// 对所有用户通知
	if err := addUnreadNoticeToAllUser(id, data["pub_at"].(string), user.IsSchAdmin); err != nil {
		return err
	}

	addManagerLog(util.AsUint(uid), id, ManagerLogType.NOTICE_NEW)
	return nil
}

func AddNoticeTemplate(data map[string]interface{}) (uint64, error) {
	var notice = Notice{
		Sender:  data["sender"].(string),
		Title:   data["title"].(string),
		Content: data["content"].(string),
		Symbol:  data["symbol"].(string),
	}
	// 创建模板
	err := db.Create(&notice).Error
	return notice.Id, err
}

func EditNoticeTemplate(uid string, id uint64, data map[string]interface{}) error {
	var (
		notice Notice
		user   User
	)
	db.Where("id = ?", id).Find(&notice)
	db.Where("id = ?", uid).Find(&user)
	if user.IsSchAdmin && notice.Symbol != NOTICE_DEPARTMENT {
		return fmt.Errorf("不能修改非部门公告")
	}
	if err := db.Where("id = ?", id).Updates(&Notice{
		Sender:  data["sender"].(string),
		Title:   data["title"].(string),
		Content: data["content"].(string),
	}).Error; err != nil {
		return err
	}
	addManagerLog(util.AsUint(uid), id, ManagerLogType.NOTICE_EDIT)
	return nil
}

func DeleteNoticeTemplate(uid string, id uint64) (uint64, error) {
	var (
		notice Notice
		user   User
	)
	db.Where("id = ?", id).Find(&notice)
	db.Where("id = ?", uid).Find(&user)
	if user.IsSchAdmin && notice.Symbol != NOTICE_DEPARTMENT {
		return 0, fmt.Errorf("不能删除非部门公告")
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&notice).Error; err != nil {
			return err
		}
		return tx.Where("notice_id = ?", id).Delete(&LogUnreadNotice{}).Error
	})
	addManagerLog(util.AsUint(uid), id, ManagerLogType.NOTICE_DELETE)
	return notice.Id, err
}
