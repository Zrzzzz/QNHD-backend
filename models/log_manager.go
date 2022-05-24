package models

import (
	ManagerLogType "qnhd/enums/MangerLogType"
)

type LogManager struct {
	// 通知归属
	Uid       uint64 `json:"uid"`
	ObjectId  uint64 `json:"object_id"`
	Type      string `json:"type"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

// 添加管理员日志
func addManagerLog(uid uint64, objectId uint64, t ManagerLogType.Enum) error {
	return db.Create(&LogManager{
		Uid:      uid,
		ObjectId: objectId,
		Type:     t.GetSymbol(),
		Detail:   "",
	}).Error
}

// 添加带信息的管理员日志
func addManagerLogWithDetail(uid uint64, objectId uint64, t ManagerLogType.Enum, detail string) error {
	return db.Create(&LogManager{
		Uid:      uid,
		ObjectId: objectId,
		Type:     t.GetSymbol(),
		Detail:   detail,
	}).Error
}

func AddManagerLogWithDetail(uid uint64, objectId uint64, t ManagerLogType.Enum, detail string) error {
	return db.Create(&LogManager{
		Uid:      uid,
		ObjectId: objectId,
		Type:     t.GetSymbol(),
		Detail:   detail,
	}).Error
}
