package UserLevelOperationType

var point = map[Enum]int{
	VISIT_POST: 3,
	ADD_POST:   2,
	ADD_FLOOR:  1,
	SHARE_POST: 3,
	// 帖子被加精
	POST_RECOMMENDED: 100,
	// 举报受理
	REPORT_PASSED: 5,

	POST_DELETED:  -10,
	FLOOR_DELETED: -8,
	BLOCKED_1:     -1,
	BLOCKED_3:     -3,
	BLOCKED_7:     -7,
	BLOCKED_14:    -14,
	BLOCKED_30:    -30,
}

func (code Enum) GetPoint() int {
	return point[code]
}
