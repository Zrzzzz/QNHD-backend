package models

type LogPostDepartmentTransfer struct {
	Doer      uint64 `json:"doer"`
	PostId    uint64 `json:"post_id"`
	Raw       uint64 `json:"raw"`
	New       uint64 `json:"new"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

type LogPostTypeTransfer struct {
	Doer      uint64 `json:"doer"`
	PostId    uint64 `json:"post_id"`
	Raw       int    `json:"raw"`
	New       int    `json:"new"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

func AddPostDepartmentTransferLog(uid, postId, rawDepartment, newDepartment uint64) error {
	return db.Create(&LogPostDepartmentTransfer{
		Doer:   uid,
		PostId: postId,
		Raw:    rawDepartment,
		New:    newDepartment,
	}).Error
}

func AddPostTypeTransferLog(uid, postId uint64, rawType, newType int) error {
	return db.Create(&LogPostTypeTransfer{
		Doer:   uid,
		PostId: postId,
		Raw:    rawType,
		New:    newType,
	}).Error
}
