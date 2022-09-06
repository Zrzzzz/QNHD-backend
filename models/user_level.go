package models

import (
	"fmt"
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/enums/UserLevelOperationType"
)

func EditUserLevel(uid string, t UserLevelOperationType.Enum) error {
	switch t {
	// 加分
	case UserLevelOperationType.VISIT_POST:
		return addVisitExp(uid)
	case UserLevelOperationType.ADD_POST:
		return addPostExp(uid)
	case UserLevelOperationType.ADD_FLOOR:
		return addFloorExp(uid)
	case UserLevelOperationType.SHARE_POST:
		return sharePostExp(uid)
	case UserLevelOperationType.POST_RECOMMENDED,
		UserLevelOperationType.REPORT_PASSED:
		return ChangeUserExp(uid, t.GetPoint())
	// 扣分
	case UserLevelOperationType.FLOOR_DELETED, UserLevelOperationType.POST_DELETED:
		return deleteExp(uid, t)
	case UserLevelOperationType.BLOCKED_1,
		UserLevelOperationType.BLOCKED_3,
		UserLevelOperationType.BLOCKED_7,
		UserLevelOperationType.BLOCKED_14,
		UserLevelOperationType.BLOCKED_30:
		return blockExp(uid, t)
	}
	return nil
}

func ChangeUserExp(uid string, change int) error {
	var user User
	db.Where("id = ?", uid).Find(&user)
	return db.Model(&User{}).Where("id = ?", uid).Update("level_point", user.LevelPoint+change).Error
}

func addVisitExp(uid string) error {
	var logs []LogVisitHistory
	db.Where("uid = ? AND created_at >= CURRENT_DATE", uid).Find(&logs)
	if len(logs) > 1 {
		return nil
	}
	return ChangeUserExp(uid, UserLevelOperationType.VISIT_POST.GetPoint())
}

func addPostExp(uid string) error {
	var posts []Post
	db.Where("uid = ? AND created_at >= CURRENT_DATE", uid).Find(&posts)
	if len(posts) > 3 {
		return nil
	}
	return ChangeUserExp(uid, UserLevelOperationType.ADD_POST.GetPoint())
}

func addFloorExp(uid string) error {
	var floors []Floor
	db.Where("uid = ? AND created_at >= CURRENT_DATE", uid).Find(&floors)
	if len(floors) > 3 {
		return nil
	}
	return ChangeUserExp(uid, UserLevelOperationType.ADD_FLOOR.GetPoint())
}

func sharePostExp(uid string) error {
	var logs []LogShare
	db.Where("uid = ? AND created_at >= CURRENT_DATE", uid).Find(&logs)
	if len(logs) > 1 {
		return ChangeUserExp(uid, UserLevelOperationType.SHARE_POST.GetPoint())
	}
	return nil
}

func deleteExp(uid string, t UserLevelOperationType.Enum) error {
	// 如果三天内有被删的记录，加扣100%
	if t == UserLevelOperationType.POST_DELETED {
		// 找用户近三天的帖子
		var ids []uint64
		db.Model(&Post{}).Select("id").Where("uid = ? AND deleted_at IS NOT NULL", uid).Find(&ids)
		// 如果没有
		if len(ids) == 0 {
			return ChangeUserExp(uid, t.GetPoint())
		} else {
			var logs []LogManager
			db.Where("object_id IN (?) AND type = ?", logs, ManagerLogType.POST_DELETE.GetSymbol()).Find(&logs)
			if len(logs) != 0 {
				return ChangeUserExp(uid, t.GetPoint()*2)
			}
		}
	} else if t == UserLevelOperationType.FLOOR_DELETED {
		// 找用户近三天的帖子
		var ids []uint64
		db.Model(&Floor{}).Select("id").Where("uid = ? AND deleted_at IS NOT NULL", uid).Find(&ids)
		// 如果没有
		if len(ids) == 0 {
			return ChangeUserExp(uid, t.GetPoint())
		} else {
			var logs []LogManager
			db.Where("object_id IN (?) AND type = ?", logs, ManagerLogType.FLOOR_DELETE.GetSymbol()).Find(&logs)
			if len(logs) != 0 {
				return ChangeUserExp(uid, t.GetPoint()*2)
			}
		}
	}
	return nil
}

func blockExp(uid string, t UserLevelOperationType.Enum) error {
	// 找五天内的禁言记录
	var logs []LogManager
	db.Where("object_id = ? AND type = ?", uid, ManagerLogType.USER_BLOCK.GetSymbol()).Find(&logs)
	if len(logs) != 0 {
		return ChangeUserExp(uid, int(float64(t.GetPoint())*1.5))
	} else {
		return ChangeUserExp(uid, t.GetPoint())
	}
}

// 检查是否能发帖子
func AddPostCheck(uid string) error {
	var (
		point   int
		postcnt int64
	)
	db.Model(&User{}).Select("level_point").Where("id = ?", uid).Find(&point)
	db.Model(&Post{}).Select("uid = ? AND created_at >= CURRENT_DATE", uid).Count(&postcnt)
	linfo := GetLevelInfo(point)
	e := fmt.Errorf("发帖数量受限")
	if linfo.Level < -30 && postcnt >= 1 ||
		linfo.Level < -20 && postcnt >= 1 ||
		linfo.Level < -10 && postcnt >= 3 {
		return e
	}
	return nil
}

// 检查是否能发评论
func AddFloorCheck(uid string) error {
	var (
		point    int
		floorcnt int64
	)
	db.Model(&User{}).Select("level_point").Where("id = ?", uid).Find(&point)
	db.Model(&Floor{}).Select("uid = ? AND created_at >= CURRENT_DATE", uid).Count(&floorcnt)
	linfo := GetLevelInfo(point)
	e := fmt.Errorf("发帖数量受限")
	if linfo.Level < -30 && floorcnt >= 2 ||
		linfo.Level < -20 && floorcnt >= 5 ||
		linfo.Level < -10 && floorcnt >= 8 {
		return e
	}
	return nil
}
