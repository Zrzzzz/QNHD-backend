package models

import (
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
)

type PostReplyExcelData struct {
	CreateMon  string `json:"create_mon"`
	CreateAt   string `json:"created_at"`
	TransAt    string `json:"trans_at"`
	IsReply    bool   `json:"is_reply"`
	Reply      string `json:"reply"`
	ReplyAt    string `json:"reply_at"`
	Department string `json:"department"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	UserNumber string `json:"user_number"`
	Username   string `json:"username"`
	Rating     int    `json:"rating"`
}

func ExportPostReplyExcel(startTime string, endTime string, c *gin.Context) ([]PostReplyExcelData, int64, error) {
	var posts []Post
	var cnt int64
	d := db.Model(&Post{}).Where("type = 1").Where("created_at >= ?", startTime).Order("created_at DESC")
	if endTime != "" {
		d.Where("created_at <= ?", endTime)
	}
	// 全部数量
	if err := d.Count(&cnt).Error; err != nil {
		return nil, cnt, err
	}
	// 找到提问问题
	if err := d.Scopes(util.Paginate(c)).Find(&posts).Error; err != nil {
		return nil, cnt, err
	}
	var ret = []PostReplyExcelData{}
	for _, p := range posts {
		var c = carbon.Parse(p.CreatedAt, "Asia/Shanghai")
		var r = PostReplyExcelData{
			CreateMon: c.Format("Y-m", "Asia/Shanghai"),
			CreateAt:  c.Format("Y-m-d H:i:s", "Asia/Shanghai"),
			IsReply:   p.Solved == 2 || p.Solved == 1,
			Title:     p.Title,
			Content:   p.Content,
			Rating:    int(p.Rating),
		}

		if p.Solved == 2 || p.Solved == 1 {
			prs, _ := GetPostReplys(util.AsStrU(p.Id))
			if len(prs) > 0 {
				r.ReplyAt = carbon.Parse(prs[0].CreatedAt, "Asia/Shanghai").Format("Y-m-d H:i:s", "Asia/Shanghai")
				r.Reply = prs[0].Content
			}
		}
		department, _ := GetDepartment(p.DepartmentId)
		r.Department = department.Name
		u, _ := GetUser(map[string]interface{}{"id": p.Uid})
		r.UserNumber = u.Number
		r.Username = u.Realname
		ret = append(ret, r)
	}
	return ret, cnt, nil
}
