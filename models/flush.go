package models

// 删除记录
func FlushOldTagLog() error {
	return db.Where("created_at <= CURRENT_TIMESTAMP + '-2 day'").Delete(&LogTag{}).Error
}

// 清理已读楼层
// func FlushOldReadFloor() error {
// 	return db.Where()
// }
// 清理已读回复
// 清理已读点赞
// 清理已读通知
