package models

import (
	"qnhd/enums/NoticeType"
	"qnhd/pkg/template"
)

func addNoticeWithTemplate(t NoticeType.Enum, uid []uint64, args []string) error {
	if len(uid) == 0 {
		return nil
	}
	data := make(map[string]interface{})
	data["symbol"] = NoticeType.GetSymbol(t)
	list := NoticeType.GetArgs(t)
	data["args"] = template.GeneArgs(list, args)
	return addUnreadNoticeToUser(uid, data)
}
