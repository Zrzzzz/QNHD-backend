package NoticeType

var templateArgs = map[Enum][]string{
	FLOOR_REPORT_SOLVE:       {"post", "floor"},
	POST_REPORT_SOLVE:        {"post"},
	POST_VALUED:              {"post"},
	BEEN_BLOCKED:             {"reason", "day"},
	POST_DELETED:             {"post"},
	FLOOR_DELETED:            {"post", "floor"},
	POST_TYPE_TRANSFER:       {"from_type", "post", "to_type"},
	POST_DEPARTMENT_TRANSFER: {"post", "department"},
}

func (code Enum) GetArgs() []string {
	return templateArgs[code]
}
