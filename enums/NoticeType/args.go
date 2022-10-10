package NoticeType

var templateArgs = map[Enum][]string{
	FLOOR_REPORT_SOLVE:        {"post", "floor"},
	POST_REPORT_SOLVE:         {"post"},
	POST_VALUED:               {"post"},
	FLOOR_VALUED:              {"floor"},
	BEEN_BLOCKED:              {"reason", "day"},
	POST_DELETED:              {"post"},
	FLOOR_DELETED:             {"post", "floor"},
	POST_TYPE_TRANSFER:        {"from_type", "post", "to_type"},
	POST_DEPARTMENT_TRANSFER:  {"post", "department"},
	POST_DELETED_WITH_REASON:  {"post", "reason"},
	FLOOR_DELETED_WITH_REASON: {"post", "floor", "reason"},
}

func (code Enum) GetArgs() []string {
	return templateArgs[code]
}
