package models

import (
	"errors"
	"fmt"
	"qnhd/pkg/logging"
	"qnhd/pkg/upload"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostCampusType int

const (
	CAMPUS_NONE PostCampusType = iota
	CAMPUS_OLD
	CAMPUS_NEW
)

type PostType int

const (
	POST_HOLE PostType = iota
	POST_SCHOOL
	POST_ALL
)

type Post struct {
	Model
	Uid uint64 `json:"-" gorm:"column:uid"`

	// 帖子分类
	Type         PostType       `json:"type"`
	DepartmentId uint64         `json:"-" gorm:"column:department_id;default:0"`
	Campus       PostCampusType `json:"campus"`
	Solved       int            `json:"solved" gorm:"defalut:0"`

	// 帖子内容
	Title   string `json:"title"`
	Content string `json:"content"`

	// 各种数量
	FavCount  uint64 `json:"fav_count" gorm:"defalut:0"`
	LikeCount uint64 `json:"like_count" gorm:"defalut:0"`
	DisCount  uint64 `json:"-" gorm:"defalut:0"`

	// 评分
	Rating uint64 `json:"rating" gorm:"default:0"`

	UpdatedAt string `json:"-" gorm:"default:null;"`
}

type LogPostFav struct {
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}
type LogPostLike struct {
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}
type LogPostDis struct {
	Uid    uint64 `json:"uid"`
	PostId uint64 `json:"post_id"`
}

type PostResponse struct {
	Post
	Tag          *Tag            `json:"tag"`
	Floors       []FloorResponse `json:"floors"`
	CommentCount int             `json:"comment_count"`
	IsLike       bool            `json:"is_like"`
	IsDis        bool            `json:"is_dis"`
	IsFav        bool            `json:"is_fav"`
	IsOwner      bool            `json:"is_owner"`
	ImageUrls    []string        `json:"image_urls"`
	Department   *Department     `json:"department"`
}

func (p *Post) geneResponse(uid string) (PostResponse, error) {
	var pr PostResponse

	frs, err := GetShortFloorResponsesInPost(util.AsStrU(p.Id), uid)
	if err != nil {
		return pr, err
	}
	imgs, err := GetImageInPost(p.Id)
	if err != nil {
		return pr, err
	}
	pr = PostResponse{
		Post:         *p,
		Floors:       frs,
		CommentCount: getCommentCount(p.Id),
		IsLike:       IsLikePostByUid(uid, util.AsStrU(p.Id)),
		IsDis:        IsDisPostByUid(uid, util.AsStrU(p.Id)),
		IsFav:        IsFavPostByUid(uid, util.AsStrU(p.Id)),
		IsOwner:      IsOwnPostByUid(uid, util.AsStrU(p.Id)),
		ImageUrls:    imgs,
	}

	if p.DepartmentId > 0 {
		d, err := GetDepartment(p.DepartmentId)
		if err != nil {
			return pr, err
		}
		pr.Department = &d
	}
	tag, _ := GetTagInPost(util.AsStrU(p.Id))
	if tag != nil {
		pr.Tag = tag
	}

	return pr, nil
}

func transPostsToResponses(posts *[]Post, uid string) ([]PostResponse, error) {
	var prs = []PostResponse{}
	for _, p := range *posts {
		pr, err := p.geneResponse(uid)
		if err != nil {
			return prs, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func GetPost(postId string) (Post, error) {
	var post Post
	err := db.Where("id = ?", postId).First(&post).Error
	return post, err
}

func GetPostResponseAndVisit(postId string, uid string) (PostResponse, error) {
	var post Post
	var pr PostResponse
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		return pr, err
	}
	if _, err := addVisitHistory(uid, postId); err != nil {
		return pr, err
	}
	if err := addTagLogInPost(util.AsUint(postId), TAG_VISIT); err != nil {
		return pr, err
	}
	return post.geneResponse(uid)
}

func GetPostResponses(c *gin.Context, uid string, maps map[string]interface{}) ([]PostResponse, error) {
	var posts []Post
	content := maps["content"].(string)
	postType := maps["type"].(PostType)
	departmentId := maps["department_id"].(string)
	solved := maps["solved"].(string)
	tagId := maps["tag_id"].(string)

	var d = db.Scopes(util.Paginate(c)).Where("CONCAT(title,content) LIKE ?", "%"+content+"%").Order("created_at DESC")
	// 校区 不为全部时加上区分
	if postType != POST_ALL {
		d = d.Where("type = ?", postType)
	}
	// 如果有部门要加上
	if departmentId != "" {
		d = d.Where("department_id = ?", departmentId)
	}
	// 如果要加上是否解决的字段
	if solved != "" {
		d = d.Where("solved = ?", solved)
	}
	// 如果需要搜索标签
	if tagId != "" {
		// 搜索相关帖子
		var tagIds = []uint64{}
		// 不需要处理错误，空的返回也行
		db.Model(&PostTag{}).Select("post_id").Where("tag_id = ?", tagId).Find(&tagIds)
		// 然后加上条件
		d = d.Where("id IN (?)", tagIds)
		// 添加搜索记录
		addTagLog(util.AsUint(tagId), TAG_VISIT)
	}

	if err := d.Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponses(&posts, uid)
}

func GetUserPostResponses(c *gin.Context, uid string) ([]PostResponse, error) {
	var posts []Post
	if err := db.Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("id DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponses(&posts, uid)
}

func GetFavPostResponses(c *gin.Context, uid string) ([]PostResponse, error) {
	var posts []Post
	if err := db.Joins("JOIN log_post_fav ON posts.id = log_post_fav.post_id").Scopes(util.Paginate(c)).Order("id DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponses(&posts, uid)
}

func GetHistoryPostResponses(c *gin.Context, uid string) ([]PostResponse, error) {
	var posts []Post
	var ids []string
	if err := db.Model(&LogVisitHistory{}).Where("uid = ?", uid).Distinct("post_id").Scopes(util.Paginate(c)).Scan(&ids).Error; err != nil {
		return nil, err
	}

	if err := db.Where("id IN (?)", ids).Scopes(util.Paginate(c)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponses(&posts, uid)
}

func AddPost(maps map[string]interface{}) (uint64, error) {
	var err error
	var post = &Post{
		Type:    maps["type"].(PostType),
		Uid:     maps["uid"].(uint64),
		Campus:  maps["campus"].(PostCampusType),
		Title:   maps["title"].(string),
		Content: maps["content"].(string),
	}
	if post.Type == POST_HOLE {
		imgs, img_ok := maps["image_urls"].([]string)
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Select("type", "uid", "campus", "title", "content").Create(post).Error; err != nil {
				return err
			}
			if img_ok {
				if err := AddImageInPost(tx, post.Id, imgs); err != nil {
					return err
				}
			}
			// 如果有tag_id
			tagId, ok := maps["tag_id"].(string)
			if ok {
				if err := AddPostWithTag(post.Id, tagId); err != nil {
					return err
				}
				// 对帖子的tag增加记录
				addTagLog(util.AsUint(tagId), TAG_ADDPOST)
			}
			return nil
		})
		if err != nil {
			upload.DeleteImageUrls(imgs)
			return 0, err
		}
	} else if post.Type == POST_SCHOOL {
		// 先对department_id进行查找，不存在要报错
		departId := maps["department_id"].(uint64)
		if err = db.Where("id = ?", departId).First(&Department{}).Error; err != nil {
			return 0, err
		}
		post.DepartmentId = departId
		imgs, img_ok := maps["image_urls"].([]string)
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Select("type", "uid", "campus", "title", "content", "department_id").Create(post).Error; err != nil {
				return err
			}

			if img_ok {
				if err := AddImageInPost(tx, post.Id, imgs); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			upload.DeleteImageUrls(imgs)
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("invalid post type")
	}

	return post.Id, nil
}

func EditPostSolved(postId string, rating string) error {
	err := db.Model(&Post{}).Where("id = ?", postId).Updates(map[string]interface{}{
		"solved": 1,
		"rating": rating,
	}).Error
	return err
}

func DeletePostsUser(id, uid string) (uint64, error) {
	var post = Post{}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND uid = ?", id, uid).First(&post).Error; err != nil {
			return err
		}
		return deletePost(tx, &post)
	})
	if err != nil {
		return 0, err
	}
	return post.Id, nil
}

func DeletePostsAdmin(uid, postId string) (uint64, error) {
	// 首先判断是否有权限
	var post, _ = GetPost(postId)
	// 如果能删，要么是超管 要么是湖底帖且是湖底管理员
	// 如果不是超管
	if RequireRight(uid, UserRight{Super: true}) != nil {
		return 0, fmt.Errorf("无权删除")
	}
	// 湖底帖且是湖底管理员
	if !(post.Type == POST_HOLE && RequireRight(uid, UserRight{StuAdmin: true}) == nil) {
		return 0, fmt.Errorf("无权删除")
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", uid).First(&post).Error; err != nil {
			return err
		}
		return deletePost(tx, &post)
	})
	if err != nil {
		return 0, err
	}

	return post.Id, nil
}

func deletePost(tx *gorm.DB, post *Post) error {
	if err := tx.Delete(&post).Error; err != nil {
		return err
	}
	if err := DeleteTagInPost(tx, post.Id); err != nil {
		return err
	}
	if err := DeleteFloorsInPost(tx, post.Id); err != nil {
		return err
	}
	if err := DeleteImageInPost(tx, post.Id); err != nil {
		return err
	}
	return nil
}

func FavPost(postId string, uid string) (uint64, error) {
	var log LogPostFav

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid > 0 {
		return 0, fmt.Errorf("已收藏")
	}

	log.Uid = util.AsUint(uid)
	log.PostId = util.AsUint(postId)
	if err := db.Select("uid", "post_id").Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新收藏数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&post).Update("fav_count", post.FavCount+1).Error; err != nil {
		return 0, err
	}
	return post.FavCount, nil
}

func UnfavPost(postId string, uid string) (uint64, error) {
	var log LogPostFav

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid == 0 {
		return 0, fmt.Errorf("未收藏")
	}

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Delete(&log).Error; err != nil {
		return 0, err
	}

	// 更新收藏数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&post).Update("fav_count", post.FavCount-1).Error; err != nil {
		return 0, err
	}
	return post.FavCount, nil
}

func LikePost(postId string, uid string) (uint64, error) {
	var log LogPostLike

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		return 0, err
	}

	if log.Uid > 0 {
		return 0, fmt.Errorf("已点赞")
	}
	log.Uid = util.AsUint(uid)
	log.PostId = util.AsUint(postId)
	if err := db.Select("uid", "post_id").Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新点赞数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&post).Update("like_count", post.LikeCount+1).Error; err != nil {
		return 0, err
	}
	UnDisPost(postId, uid)
	return post.LikeCount, nil
}

func UnLikePost(postId string, uid string) (uint64, error) {
	var log LogPostLike

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		return 0, err
	}

	if log.Uid == 0 {
		return 0, fmt.Errorf("未点赞")
	}

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Delete(&log).Error; err != nil {
		return 0, err
	}

	// 更新点赞数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&post).Update("like_count", post.LikeCount-1).Error; err != nil {
		return 0, err
	}
	return post.LikeCount, nil
}

func DisPost(postId string, uid string) (uint64, error) {
	var log LogPostDis

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid > 0 {
		return 0, fmt.Errorf("已点踩")
	}
	log.Uid = util.AsUint(uid)
	log.PostId = util.AsUint(postId)
	if err := db.Select("uid", "post_id").Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新点踩数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&post).Update("dis_count", post.DisCount+1).Error; err != nil {
		return 0, err
	}
	UnLikePost(postId, uid)
	return post.DisCount, nil
}

func UnDisPost(postId string, uid string) (uint64, error) {
	var log LogPostDis

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid == 0 {
		return 0, fmt.Errorf("未点踩")
	}

	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Delete(&log).Error; err != nil {
		return 0, err
	}

	// 更新楼的点踩数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&post).Update("dis_count", post.DisCount-1).Error; err != nil {
		return 0, err
	}
	return post.DisCount, nil
}

func IsLikePostByUid(uid, postId string) bool {
	var log LogPostLike
	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Uid > 0
}

func IsDisPostByUid(uid, postId string) bool {
	var log LogPostDis
	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Uid > 0
}

func IsFavPostByUid(uid, postId string) bool {
	var log LogPostFav
	if err := db.Where("uid = ? AND post_id = ?", uid, postId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Uid > 0
}

func IsOwnPostByUid(uid, postId string) bool {
	var post, err = GetPost(postId)
	if err != nil {
		return false
	}
	return util.AsStrU(post.Uid) == uid
}

func (LogPostFav) TableName() string {
	return "log_post_fav"
}

func (LogPostLike) TableName() string {
	return "log_post_like"
}

func (LogPostDis) TableName() string {
	return "log_post_dis"
}

func (Post) TableName() string {
	return "posts"
}
