package models

type Model struct {
	ID       int `gorm:"primary_key" json:"id"`
	CreateAt int `json:"create_at"`
	DeleteAt int `json:"delete_at"`
}
