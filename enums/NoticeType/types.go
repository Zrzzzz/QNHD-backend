package NoticeType

type Enum int

const (
	FLOOR_REPORT_SOLVE Enum = iota
	POST_REPORT_SOLVE
	POST_VALUED
	FLOOR_VALUED
	BEEN_BLOCKED
	POST_DELETED
	FLOOR_DELETED
	POST_TYPE_TRANSFER
	POST_DEPARTMENT_TRANSFER
	POST_DELETED_WITH_REASON
	FLOOR_DELETED_WITH_REASON
)
