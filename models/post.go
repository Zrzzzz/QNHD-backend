package models

type Post struct {
	Model
	Uid        uint64 `json:"uid"`
	Content    string `json:"content"`
	PictureUrl string `json:"picture_url"`
	UpdatedAt  string `json:"updated_at" gorm:"null;"`
}

func GetPost(id string) (post Post) {
	db.Where("id = ?", id).First(&post)
	return
}

func GetPosts(overNum, limit int, content string) (posts []Post) {
	db.Where("content LIKE ?", "%"+content+"%").Offset(overNum).Limit(limit).Find(&posts)
	return
}

func AddPosts(maps map[string]interface{}) bool {
	var post = &Post{
		Uid:        maps["uid"].(uint64),
		Content:    maps["content"].(string),
		PictureUrl: maps["picture_url"].(string),
	}
	db.Select("uid", "content", "picture_url").Create(post)
	tags, ok := maps["tags"].([]string)
	if ok {
		AddPostWithTag(post.Id, tags)
	}
	return true
}

func DeletePostsUser(id, uid string) bool {
	db.Where("id = ? AND uid = ?", id, uid).Delete(&Post{})
	DeleteTagInPost(id)
	DeleteFloorsInPost(id)
	return true
}

func DeletePostsAdmin(id string) bool {
	db.Where("id = ?", id).Delete(&Post{})
	DeleteTagInPost(id)
	DeleteFloorsInPost(id)
	return true
}

func (Post) TableName() string {
	return "posts"
}
