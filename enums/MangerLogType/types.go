package ManagerLogType

type Enum int

const (
	USER_BAN Enum = iota
	USER_UNBAN

	USER_BLOCK
	USER_UNBLOCK

	POST_DELETE
	FLOOR_DELETE

	POST_ETAG
	POST_UNETAG

	POST_TOP
	POST_UNTOP

	POST_REPLY
	POST_DEPARTMENT_TRANSFER
	POST_TPYE_TRANSFER
	POST_DEPARTMENT_DISTRIBUTE

	USER_ADD
	USER_PERMISSION_CHANGE

	NOTICE_NEW
	NOTICE_DELETE
	NOTICE_EDIT

	USER_DETAIL

	TAG_POINT_ADD
	TAG_POINT_CLEAR
	TAG_DELETE
)
