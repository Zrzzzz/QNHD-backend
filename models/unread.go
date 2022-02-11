package models

type MessageCount struct {
	Floor  int `json:"floor"`
	Reply  int `json:"reply"`
	Notice int `json:"notice"`
}

// 获取总未读数
func GetMessageCount(uid string) (MessageCount, error) {
	var ret = MessageCount{}
	// 楼层未读 回复未读 通知未读
	var fcnt, rcnt, ncnt int64
	// 获取楼层未读数
	if err := db.Model(&LogUnreadFloor{}).Where("uid = ? AND is_read = true", uid).Count(&fcnt).Error; err != nil {
		return ret, err
	}
	// 获取回复未读数
	if err := db.Model(&LogUnreadPostReply{}).Where("uid = ? AND is_read = true", uid).Count(&rcnt).Error; err != nil {
		return ret, err
	}
	// 获取通知未读数
	if err := db.Model(&LogUnreadNotice{}).Where("uid = ?", uid).Count(&ncnt).Error; err != nil {
		return ret, err
	}
	ret.Floor = int(fcnt)
	ret.Reply = int(rcnt)
	ret.Notice = int(ncnt)
	return ret, nil
}

// 全部已读
func ReadAllMessage(uid uint64) error {
	if err := db.Model(&LogUnreadFloor{}).Where("uid = ?", uid).Update("is_read", true).Error; err != nil {
		return err
	}
	if err := db.Model(&LogUnreadPostReply{}).Where("uid = ?", uid).Update("is_read", true).Error; err != nil {
		return err
	}
	if err := db.Where("uid = ?", uid).Delete(&LogUnreadLike{}).Error; err != nil {
		return err
	}
	return db.Where("uid = ?", uid).Delete(&LogUnreadNotice{}).Error
}
