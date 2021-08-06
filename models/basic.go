package models

type Model struct {
	Id       uint64 `gorm:"primary_key" json:"id"`
	CreatedAt string `json:"create_at"`
	DeletedAt string `json:"delete_at" gorm:"null;"`
}
