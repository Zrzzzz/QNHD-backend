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
	Raw       uint64 `json:"raw"`
	New       uint64 `json:"new"`
	CreatedAt string `json:"created_at" gorm:"default:null;"`
}

func AddPostDepartmentTransferLog(uid, postId, rawType, newType uint64) error {
	return db.Create(&LogPostDepartmentTransfer{
		Doer:   uid,
		PostId: postId,
		Raw:    rawType,
		New:    newType,
	}).Error
}

func AddPostTypeTransferLog(uid, postId, rawType, newType uint64) error {
	return db.Create(&LogPostDepartmentTransfer{
		Doer:   uid,
		PostId: postId,
		Raw:    rawType,
		New:    newType,
	}).Error
}
