package ManagerLogType

var msgSymbol = map[Enum]string{
	USER_BAN:   "user_ban",
	USER_UNBAN: "user_unban",

	USER_BLOCK:   "user_block",
	USER_UNBLOCK: "user_unblock",

	POST_DELETE:  "post_delete",
	FLOOR_DELETE: "floor_delete",

	POST_ETAG:   "post_etag",
	POST_UNETAG: "post_unetag",

	POST_TOP:   "post_top",
	POST_UNTOP: "post_untop",

	POST_REPLY:                 "post_reply",
	POST_REPLY_MODIFY:          "post_reply_modify",
	POST_DEPARTMENT_TRANSFER:   "post_department_transfer",
	POST_DEPARTMENT_DISTRIBUTE: "post_department_distribute",
	POST_TPYE_TRANSFER:         "post_type_transfer",
	POST_RETURN:                "post_return",
	POST_EDIT_COMMENTABLE:      "post_edit_commentable",

	FLOOR_EDIT_COMMENTABLE: "floor_edit_commentable",

	USER_ADD:               "user_add",
	USER_PERMISSION_CHANGE: "user_permission_change",

	NOTICE_NEW:    "notice_new",
	NOTICE_DELETE: "notice_delete",
	NOTICE_EDIT:   "notice_edit",

	USER_DETAIL:         "user_detail",
	USER_NICKNAME_RESET: "user_nickname_reset",
	USER_AVATAR_RESET:   "user_avatar_reset",

	TAG_POINT_ADD:   "tag_point_add",
	TAG_POINT_CLEAR: "tag_point_clear",
	TAG_DELETE:      "tag_delete",
}

func (code Enum) GetSymbol() string {
	return msgSymbol[code]
}
