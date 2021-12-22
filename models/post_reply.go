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

func GetPostReplys(postId string) ([]PostReply, error) {
	var prs = []PostReply{}
	err := db.Where("post_id = ?", postId).Find(&prs).Error
	return prs, err
}

func AddPostReply(maps map[string]interface{}) error {
	err := db.Create(&PostReply{
		PostId:  maps["post_id"].(uint64),
		From:    maps["from"].(PostReplyType),
		Content: maps["content"].(string),
	}).Error
	return err
}

func (PostReply) TableName() string {
	return "post_reply"
}
