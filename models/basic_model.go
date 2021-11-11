package models

import "gorm.io/gorm"

type Model struct {
	Id        uint64         `gorm:"primaryKey;autoIncrement;" json:"id"`
	CreatedAt string         `json:"create_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ListRes struct {
	List  interface{} `json:"list"`
	Total int         `json:"total" example:"1"`
}

type IdRes struct {
	Id int `json:"id"`
}
