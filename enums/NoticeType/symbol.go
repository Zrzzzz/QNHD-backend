package NoticeType

var msgSymbol = map[Enum]string{
	FLOOR_REPORT_SOLVE:       "floor_report_solve",
	POST_REPORT_SOLVE:        "post_report_solve",
	POST_VALUED:              "post_valued",
	BEEN_BLOCKED:             "been_blocked",
	POST_DELETED:             "post_deleted",
	FLOOR_DELETED:            "floor_deleted",
	POST_TYPE_TRANSFER:       "post_type_transfer",
	POST_DEPARTMENT_TRANSFER: "post_department_transfer",
}

func (code Enum) GetSymbol() string {
	return msgSymbol[code]
}
