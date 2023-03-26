package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)
const (
	ERROR_EXIST_EMAIL = 10001 + iota
	ERROR_NOT_EXIST_EMAIL
	ERROR_NOT_EXIST_ARTICLE
	ERROR_EXIST_USER
	ERROR_NOT_EXIST_USER
	ERROR_EXIST_TAG
	ERROR_NOT_EXIST_TAG
	ERROR_EXIST_DEPARTMENT
	ERROR_NOT_EXIST_DEPARTMENT
	ERROR_POST_TYPE
)

const (
	ERROR_BANNED_USER = 10101 + iota
	ERROR_NOT_BANNED_USER
	ERROR_BLOCKED_USER
	ERROR_BLOCKED_USER_DAY
	ERROR_NOT_BLOCKED_USER
	ERROR_POST_COUNT_LIMITED
	ERROR_FLOOR_COUNT_LIMITED
)

const (
	ERROR_AUTH_CHECK_TOKEN_FAIL = 20001 + iota
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT
	ERROR_GENERATE_TOKEN
	ERROR_AUTH
	ERROR_RIGHT
)

const (
	ERROR_SEND_EMAIL = 30001 + iota
	ERROR_SAVE_FILE
	ERROR_SERVER
)

const (
	ERROR_DATABASE = 40001 + iota
)
