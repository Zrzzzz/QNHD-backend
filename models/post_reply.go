package models

type PostReplyType int

const (
	PostReplyFromUser PostReplyType = iota
	PostReplyFromSchool
)

type PostReply struct {
	Model
	PostId  uint64        `json:"post_id"`
	From    PostReplyType `json:"from"`
	Content string        `json:"content"`
}

// 获取单个回复
func GetPostReply(replyId string) (PostReply, error) {
	var pr PostReply
	err := db.Where("id = ?", replyId).Find(&pr).Error
	return pr, err
}

// 获取帖子内的回复记录
func GetPostReplys(postId string) ([]PostReply, error) {
	var prs = []PostReply{}
	err := db.Where("post_id = ?", postId).Find(&prs).Error
	return prs, err
}

// 添加帖子的回复
func AddPostReply(maps map[string]interface{}) (uint64, error) {
	var pr = PostReply{
		PostId:  maps["post_id"].(uint64),
		From:    maps["from"].(PostReplyType),
		Content: maps["content"].(string),
	}
	err := db.Create(&pr).Error
	return pr.Id, err
}

func (PostReply) TableName() string {
	return "post_reply"
}
