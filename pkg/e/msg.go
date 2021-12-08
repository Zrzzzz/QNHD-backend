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

	ERROR_BANNED_USER:      "用户已被封禁",
	ERROR_NOT_BANNED_USER:  "用户未被封禁",
	ERROR_BLOCKED_USER:     "用户已被禁言",
	ERROR_NOT_BLOCKED_USER: "用户未被禁言",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_GENERATE_TOKEN:           "Token生成失败",
	ERROR_AUTH:                     "账号密码错误",
	ERROR_RIGHT:                    "无权访问",

	ERROR_SEND_EMAIL:               "发送邮件失败",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:   "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:  "检查图片失败",
	ERROR_UPLOAD_SAVE_IMAGE_FORMAT: "图片格式或大小不合规范",
	ERROR_EMAIL_CODE_CHECK:         "邮箱验证码错误",

	ERROR_DATABASE: "数据库错误，请上报管理员",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
