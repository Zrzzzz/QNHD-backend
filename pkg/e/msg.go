package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_EXIST_EMAIL:          "该邮箱已存在",
	ERROR_NOT_EXIST_EMAIL:      "该邮箱不存在",
	ERROR_NOT_EXIST_ARTICLE:    "该文章不存在",
	ERROR_EXIST_USER:           "该用户已存在",
	ERROR_NOT_EXIST_USER:       "该用户不存在",
	ERROR_EXIST_TAG:            "该标签已存在",
	ERROR_NOT_EXIST_TAG:        "该标签不存在",
	ERROR_EXIST_DEPARTMENT:     "该部门已存在",
	ERROR_NOT_EXIST_DEPARTMENT: "该部门不存在",
	ERROR_POST_TYPE:            "帖子类型错误",

	ERROR_BANNED_USER:         "用户已被封禁",
	ERROR_NOT_BANNED_USER:     "用户未被封禁",
	ERROR_BLOCKED_USER:        "用户已被禁言",
	ERROR_BLOCKED_USER_DAY:    "禁言天数不在范围",
	ERROR_NOT_BLOCKED_USER:    "用户未被禁言",
	ERROR_POST_COUNT_LIMITED:  "发帖数量到达上限",
	ERROR_FLOOR_COUNT_LIMITED: "评论数量到达上限",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_GENERATE_TOKEN:           "Token生成失败",
	ERROR_AUTH:                     "账号密码错误",
	ERROR_RIGHT:                    "无权访问",

	ERROR_SEND_EMAIL: "发送邮件失败",
	ERROR_SAVE_FILE:  "保存文件失败",
	ERROR_SERVER:     "服务器错误",

	ERROR_DATABASE: "数据库错误，请上报管理员",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
