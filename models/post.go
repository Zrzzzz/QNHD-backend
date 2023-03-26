package models

import (
	"fmt"
	"qnhd/enums/LikeType"
	ManagerLogType "qnhd/enums/MangerLogType"
	"qnhd/enums/NoticeType"
	"qnhd/enums/PostCampusType"
	"qnhd/enums/PostEtagType"
	"qnhd/enums/PostSearchModeType"
	"qnhd/enums/PostSolveType"
	"qnhd/enums/PostValueModeType"
	"qnhd/enums/ReportType"
	"qnhd/enums/TagPointType"
	"qnhd/enums/UserLevelOperationType"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/segment"

	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
	giterrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

const POST_ALL = 0

type Post struct {
	Model
	Uid uint64 `json:"uid" gorm:"column:uid"`

	// 帖子分类
	Type         int                 `json:"type"`
	DepartmentId uint64              `json:"-" gorm:"column:department_id;default:0"`
	Campus       PostCampusType.Enum `json:"campus"`
	// 0 已提问 1 已回复 2 已解决 3 已分发
	Solved PostSolveType.Enum `json:"solved" gorm:"default:0"`

	// 帖子内容
	Title    string `json:"title"`
	Content  string `json:"content"`
	Nickname string `json:"nickname"`

	// 各种数量
	FavCount  uint64 `json:"fav_count" gorm:"default:0"`
	LikeCount uint64 `json:"like_count" gorm:"default:0"`
	DisCount  uint64 `json:"-" gorm:"default:0"`

	// 评分
	Rating uint64 `json:"rating" gorm:"default:0"`
	// 置顶值
	Value uint64 `json:"value" gorm:"default:0"`

	// 分词
	Tokens string `json:"-"`

	UpdatedAt string `json:"-" gorm:"default:null;"`

	// etag
	Etag string `json:"e_tag" gorm:"column:extra_tag;"`
	// 能否评论
	Commentable bool `json:"commentable" gorm:"default:true"`
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

// 帖子返回数据
type PostResponse struct {
	Post
	Tag          *Tag            `json:"tag"`
	Floors       []FloorResponse `json:"floors"`
	CommentCount int             `json:"comment_count"`
	ImageUrls    []string        `json:"image_urls"`
	Department   *Department     `json:"department"`
	UInfo        UserInfo        `json:"user_info"`

	IsDeleted bool `json:"is_deleted"`
	// 用于处理链式数据
	Error error `json:"-"`
}

// 客户端帖子返回数据
type PostResponseUser struct {
	Post
	Tag          *Tag                `json:"tag"`
	Floors       []FloorResponseUser `json:"floors"`
	CommentCount int                 `json:"comment_count"`
	ImageUrls    []string            `json:"image_urls"`
	Department   *Department         `json:"department"`

	IsLike     bool `json:"is_like"`
	IsDis      bool `json:"is_dis"`
	IsFav      bool `json:"is_fav"`
	IsOwner    bool `json:"is_owner"`
	VisitCount int  `json:"visit_count"`

	UInfo UserInfo `json:"user_info"`

	IsDeleted bool `json:"is_deleted"`
	// 用于处理链式数据
	Error error `json:"-"`
}

func (p *Post) geneResponse(unscoped bool) PostResponse {
	var pr PostResponse

	imgs, err := GetImageInPost(p.Id)
	if err != nil {
		pr.Error = err
		return pr
	}
	pr = PostResponse{
		Post:         *p,
		CommentCount: GetCommentCount(p.Id, true, unscoped),
		ImageUrls:    imgs,
		// user info
		UInfo: GetUserInfo(util.AsStrU(p.Uid)),
	}

	if p.DepartmentId > 0 {
		d, err := GetDepartment(p.DepartmentId)
		if err != nil {
			pr.Error = err
			return pr
		}
		pr.Department = &d
	}
	tag, _ := GetTagInPost(util.AsStrU(p.Id))
	if tag != nil {
		pr.Tag = tag
	}
	pr.Error = err
	pr.IsDeleted = pr.DeletedAt.Valid

	return pr
}

func (p PostResponse) searchByUid(uid string) PostResponseUser {
	pr := PostResponseUser{
		Post:         p.Post,
		Tag:          p.Tag,
		CommentCount: p.CommentCount,
		ImageUrls:    p.ImageUrls,
		Department:   p.Department,
		UInfo:        p.UInfo,
		IsLike:       IsLikePostByUid(uid, util.AsStrU(p.Id)),
		IsDis:        IsDisPostByUid(uid, util.AsStrU(p.Id)),
		IsFav:        IsFavPostByUid(uid, util.AsStrU(p.Id)),
		IsOwner:      IsOwnPostByUid(uid, util.AsStrU(p.Id)),
		VisitCount:   GetPostVisitCount(util.AsStrU(p.Id)),
	}

	// frs, err := getShortFloorResponsesInPostWithUid(util.AsStrU(p.Id), uid)
	// if err != nil {
	// 	pr.Error = err
	// 	return pr
	// }
	// pr.Floors = frs
	return pr
}

// 将post数组转化为返回结果，后台使用
func transPostsToResponses(posts *[]Post) ([]PostResponse, error) {
	var prs = []PostResponse{}
	var err error
	for _, p := range *posts {

		pr := p.geneResponse(true)

		if pr.Type == POST_SCHOOL_TYPE {
			var user User
			db.Where("id = ?", pr.Uid).Find(&user)
			pr.Nickname = user.realnameFull()
		}
		if pr.Error != nil {
			err = giterrors.Wrap(err, pr.Error.Error())
		} else {
			prs = append(prs, pr)
		}
	}
	return prs, err
}

// 将post数组转化为用户返回结果， 前端使用
func transPostsToResponsesWithUid(posts *[]Post, uid string) ([]PostResponseUser, error) {
	var prs = []PostResponseUser{}
	var err error
	for _, p := range *posts {
		pr := p.geneResponse(false).searchByUid(uid)
		if pr.Error != nil {
			err = giterrors.Wrap(err, pr.Error.Error())
		} else {
			prs = append(prs, pr)
		}
	}
	return prs, err
}

func GetPostVisitCount(postId string) int {
	var cnt int64
	db.Model(&LogVisitHistory{}).Where("post_id = ?", postId).Count(&cnt)
	return int(cnt)
}

func GetPost(postId string) (Post, error) {
	var post Post
	err := db.Where("id = ?", postId).First(&post).Error
	return post, err
}

// 获取未被分发的帖子
func GetUndistributedPosts(c *gin.Context) ([]PostResponse, error) {
	var posts = []Post{}
	err := db.Scopes(util.Paginate(c)).Where("type = 1 AND solved = 0").Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	ret, err := transPostsToResponses(&posts)
	return ret, err
}

// 后台使用
func GetPostResponse(postId string) (PostResponse, error) {
	var p Post
	var pr PostResponse
	err := db.Unscoped().Where("id = ?", postId).First(&p).Error
	if err != nil {
		return pr, err
	}
	pr = p.geneResponse(true)
	if pr.Type == POST_SCHOOL_TYPE {
		var user User
		db.Where("id = ?", pr.Uid).Find(&user)
		pr.Nickname = user.realnameFull()
	}
	return pr, pr.Error
}

// 前端使用
func GetPostResponseUser(postId string, uid string) (PostResponseUser, error) {
	var post Post
	var pr PostResponseUser
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		return pr, err
	}
	ret := post.geneResponse(false).searchByUid(uid)
	return ret, ret.Error
}

// front表示是否为前端请求
func getPosts(c *gin.Context, maps map[string]interface{}) ([]Post, int, error) {
	var (
		posts []Post
		cnt   int64
		err   error
	)
	content := maps["content"].(string)
	postType := maps["type"].(int)
	searchMode := maps["search_mode"].(PostSearchModeType.Enum)
	departmentId := maps["department_id"].(string)
	solved := maps["solved"].(string)
	tagId := maps["tag_id"].(string)
	valueMode := maps["value_mode"].(PostValueModeType.Enum)
	front := maps["front"].(bool)
	isDeleted := maps["is_deleted"].(string)
	commentable := maps["commentable"].(string)
	etag := maps["etag"].(string)

	var d = db.Model(&Post{})
	// 如果是前端
	if !front {
		d = d.Unscoped()
	}
	if isDeleted == "1" {
		d = d.Where("deleted_at IS NOT NULL")
	} else {
		d = d.Where("deleted_at IS NULL")
	}
	// 置顶帖搜索
	if valueMode == PostValueModeType.DEFAULT {
		d = d.Order("value DESC")
	} else if valueMode == PostValueModeType.ONLY {
		d = d.Where("value <> 0")
	} else if valueMode == PostValueModeType.NONE {
		// VALUE_NONE 不做操作
	}

	// 当搜索不为空时加上全文检索
	if content != "" {
		d = db.Select("p.*", "ts_rank(p.tokens, q) as score").
			Table("(?) as p, plainto_tsquery(?) as q", d, segment.Cut(content, " ")).
			Where("q @@ p.tokens").Order("score DESC")
	}
	// 排序方式
	if searchMode == PostSearchModeType.TIME {
		d = d.Order("created_at DESC")
	} else if searchMode == PostSearchModeType.UPDATE {
		d = d.Order("updated_at DESC")
	}

	// 分区 不为全部时加上区分
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
	}
	// 获取是否评论的
	if commentable == "1" {
		d = d.Where("commentable = ?", true)
	} else if commentable == "0" {
		d = d.Where("commentable = ?", false)
	}

	// etag区分
	if etag != "" && PostEtagType.Contains(etag) {
		d = d.Where("extra_tag = ?", etag)
	}

	// 开始搜索
	if err = d.Count(&cnt).Error; err != nil {
		return posts, int(cnt), err
	}
	// 分页
	d = d.Scopes(util.Paginate(c))
	// 这里还得加一次，上面的是子查询的
	if !front {
		d = d.Unscoped()
	}
	err = d.Find(&posts).Error
	return posts, int(cnt), err
}

// 获取帖子返回数据，后台使用
func GetPostResponses(c *gin.Context, maps map[string]interface{}) ([]PostResponse, int, error) {
	maps["front"] = false
	posts, cnt, err := getPosts(c, maps)
	if err != nil {
		return nil, 0, err
	}
	ret, err := transPostsToResponses(&posts)
	return ret, cnt, err
}

// 获取帖子返回数据带uid，前端使用
func GetPostResponsesWithUid(c *gin.Context, uid string, maps map[string]interface{}) ([]PostResponseUser, error) {
	maps["front"] = true
	posts, _, err := getPosts(c, maps)
	if err != nil {
		return nil, err
	}
	return transPostsToResponsesWithUid(&posts, uid)
}

func GetUserPostResponseWithUid(c *gin.Context, uid string) ([]PostResponseUser, error) {
	var posts []Post
	if err := db.Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("id DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponsesWithUid(&posts, uid)
}

func GetUserPostResponses(c *gin.Context, uid string, deleted bool) ([]PostResponse, error) {
	var posts []Post
	d := db.Unscoped().Where("uid = ?", uid).Scopes(util.Paginate(c)).Order("id DESC")
	if deleted {
		if err := d.Where("deleted_at IS NOT NULL").Find(&posts).Error; err != nil {
			return nil, err
		}
	} else {
		if err := d.Find(&posts).Error; err != nil {
			return nil, err
		}
	}
	return transPostsToResponses(&posts)
}

func GetFavPostResponseWithUid(c *gin.Context, uid string) ([]PostResponseUser, error) {
	var posts []Post
	if err := db.Joins(`JOIN qnhd.log_post_fav
	ON qnhd.post.id = qnhd.log_post_fav.post_id
	AND qnhd.log_post_fav.uid = ?`, uid).Scopes(util.Paginate(c)).Order("id DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponsesWithUid(&posts, uid)
}

func GetHistoryPostResponseWithUid(c *gin.Context, uid string) ([]PostResponseUser, error) {
	var posts []Post
	var ids []string
	if err := db.Model(&LogVisitHistory{}).Where("uid = ?", uid).Distinct("post_id", "created_at").Order("created_at DESC").Scopes(util.Paginate(c)).Find(&ids).Error; err != nil {
		return nil, err
	}

	if err := db.Where("id IN (?)", ids).Scopes(util.Paginate(c)).Find(&posts).Error; err != nil {
		return nil, err
	}
	return transPostsToResponsesWithUid(&posts, uid)
}

func AddPost(maps map[string]interface{}) (uint64, error) {
	var err error
	uid := maps["uid"].(uint64)
	var user User
	db.Where("id = ?", uid).Find(&user)
	var post = &Post{
		Type:     maps["type"].(int),
		Uid:      uid,
		Nickname: user.Nickname,
		Campus:   maps["campus"].(PostCampusType.Enum),
		Title:    filter.CommonFilter.Filter(maps["title"].(string)),
		Content:  filter.CommonFilter.Filter(maps["content"].(string)),
	}
	if post.Type == POST_SCHOOL_TYPE {
		// 先对department_id进行查找，不存在要报错
		departId := maps["department_id"].(uint64)
		if err = db.Where("id = ?", departId).First(&Department{}).Error; err != nil {
			return 0, err
		}
		post.Nickname = user.realname()
		post.DepartmentId = departId
		imgs, img_ok := maps["image_urls"].([]string)
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(post).Error; err != nil {
				return err
			}

			if img_ok {
				if err := AddImageInPost(tx, post.Id, imgs); err != nil {
					return err
				}
			}
			return nil
		})
	} else if IsValidPostType(post.Type) {
    logging.Debug("测试调试1")
		imgs, img_ok := maps["image_urls"].([]string)
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(post).Error; err != nil {
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
				if err := AddPostWithTag(tx, post.Id, util.AsUint(tagId)); err != nil {
          logging.Debug("error: %v", err)
					return err
				}
				// 对帖子的tag增加记录
				addTagLog(util.AsUint(tagId), TagPointType.ADD_POST)
			}
			return nil
		})
	} else {
		return 0, fmt.Errorf("invalid post type")
	}
	if err != nil {
		return 0, err
	}
	if err := flushPostTokens(post.Id, post.Title, post.Content); err != nil {
		return 0, err
	}
	// 增加经验
	EditUserLevel(util.AsStrU(uid), UserLevelOperationType.ADD_POST)

	return post.Id, nil
}

func EditPost(postId string, maps map[string]interface{}) error {
	return db.Model(&Post{}).Where("id = ?", postId).Updates(maps).Error
}

func EditPostValue(uid, postId string, value int) error {
	var post Post
	if err := db.Where("id = ?", postId).Find(&post).Error; err != nil {
		return err
	}
	// 置顶值修改
	if post.Value == 0 && value > 0 {
		addNoticeWithTemplate(NoticeType.POST_VALUED, []uint64{post.Uid}, []string{post.Title})
	}
	// 这里对标签进行操作 如果置为0，则置为none，否则设置为top
	if value > 0 {
		EditPostEtag(uid, postId, PostEtagType.TOP)
		addManagerLog(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_TOP)
	} else {
		EditPostEtag(uid, postId, PostEtagType.NONE)
		addManagerLog(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_UNTOP)
	}
	return EditPost(postId, map[string]interface{}{"value": value})
}

func EditPostEtag(uid, postId string, t PostEtagType.Enum) error {
	if t == PostEtagType.NONE {
		addManagerLog(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_UNETAG)
	} else {
		addManagerLogWithDetail(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_ETAG,
			fmt.Sprintf("to: %s", t.GetSymbol()))
	}
	// 如果是帖子加精，加经验
	if t == PostEtagType.RECOMMEND {
		var pUserId string
		db.Model(&Post{}).Select("uid").Where("id = ?", postId).Find(&pUserId)
		EditUserLevel(pUserId, UserLevelOperationType.POST_RECOMMENDED)
	}
	return db.Model(&Post{}).Where("id = ?", postId).Update("extra_tag", t.GetSymbol()).Error
}

func EditPostDepartment(uid, postId string, departmentId string) error {
	// 判断是否存在部门
	var (
		newType Department
		rawType Department
	)
	if err := db.First(&newType, departmentId).Error; err != nil {
		return err
	}
	post, err := GetPost(postId)
	if err != nil {
		return err
	}
	if err := db.First(&rawType, post.DepartmentId).Error; err != nil {
		return err
	}
	// 如果类型相同
	if rawType.Id == newType.Id {
		return fmt.Errorf("不能修改为同类型")
	}
	// 通知帖子用户
	addNoticeWithTemplate(NoticeType.POST_DEPARTMENT_TRANSFER, []uint64{post.Uid}, []string{post.Title, newType.Name})
	addManagerLogWithDetail(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_DEPARTMENT_TRANSFER,
		fmt.Sprintf("from: %s, to: %s", rawType.Name, newType.Name))

	return EditPost(postId, map[string]interface{}{"department_id": departmentId})
}

func EditPostCommentable(uid, postId string, commentable bool) error {
	addManagerLog(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_EDIT_COMMENTABLE)
	return EditPost(postId, map[string]interface{}{"commentable": commentable})
}

// 分发帖子
func DistributePost(uid string, postId string, departmentId string) error {
	// 判断是否存在部门
	var (
		newType Department
	)
	fmt.Println(departmentId)
	if err := db.First(&newType, departmentId).Error; err != nil {
		return err
	}
	addManagerLogWithDetail(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_DEPARTMENT_DISTRIBUTE,
		fmt.Sprintf("to: %s", newType.Name))
	return EditPost(postId, map[string]interface{}{"department_id": departmentId, "solved": PostSolveType.DISTRIBUTED})
}

func EditPostType(uid, postId string, typeId string) error {
	// 判断是否存在类型
	var (
		newType PostType
		rawType PostType
	)
	if err := db.First(&newType, typeId).Error; err != nil {
		return err
	}
	post, err := GetPost(postId)
	if err != nil {
		return err
	}
	if err := db.First(&rawType, post.Type).Error; err != nil {
		return err
	}
	// 如果类型相同
	if rawType.Id == newType.Id {
		return fmt.Errorf("不能修改为同类型")
	}
	// 如果要修改为校务类型，禁止操作
	if util.AsInt(typeId) == int(POST_SCHOOL_TYPE) {
		return fmt.Errorf("不能修改为校务类型")
	}
	// 如果是校务类型，需要去掉部门
	if post.Type == POST_SCHOOL_TYPE {
		if err := EditPost(postId, map[string]interface{}{"type": typeId, "department_id": 0}); err != nil {
			return err
		}
		// 需要更帖子昵称和楼层的昵称
		if err := updatePostAndFloorNickname(post); err != nil {
			return err
		}
	}
	// 更新楼层type
	// 通知帖子用户
	addNoticeWithTemplate(NoticeType.POST_TYPE_TRANSFER, []uint64{post.Uid}, []string{rawType.Name, post.Title, newType.Name})
	addManagerLogWithDetail(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_TPYE_TRANSFER,
		fmt.Sprintf("from: %s, to: %s", rawType.Name, newType.Name))

	return EditPost(postId, map[string]interface{}{"type": typeId})
}

func updatePostAndFloorNickname(post Post) error {
	var u User
	var floors []Floor
	db.Where("id = ?", post.Uid).Find(&u)
	db.Where("post_id = ?", post.Id).Find(&floors)
	if u.Uid > 0 {
		db.Model(&Post{}).Update("nickname", u.Nickname)
	}
	for _, f := range floors {
		var u User
		db.Where("id = ?", f.Uid).Find(&u)
		if u.Uid > 0 {
			db.Model(&Floor{}).Update("nickname", u.Nickname)
		}
	}
	return nil
}

func DeletePostsUser(id, uid string) (uint64, error) {
	var post = Post{}
	if err := db.Where("id = ? AND uid = ?", id, uid).First(&post).Error; err != nil {
		return 0, err
	}
	err := deletePost(&post)
	return post.Id, err
}

func DeletePostAdmin(uid, postId string, reason string) (uint64, error) {
	var post, _ = GetPost(postId)
	// 找到举报过帖子的所有用户
	var uids []uint64
	db.Model(&Report{}).Select("uid").Where("type = ? AND post_id = ?", ReportType.POST, post.Id).Find(&uids)

	err := deletePost(&post)
	if err != nil {
		return 0, err
	}
	DEFAULT_REASON := "违反社区规范"
	if reason == "" {
		reason = DEFAULT_REASON
	}

	addNoticeWithTemplate(NoticeType.POST_REPORT_SOLVE, uids, []string{post.Title})
	// 通知被删除的用户
	addNoticeWithTemplate(NoticeType.POST_DELETED_WITH_REASON, []uint64{post.Uid}, []string{post.Title, reason})
	addManagerLog(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_DELETE)
	// 将被举报的人扣经验
	EditUserLevel(util.AsStrU(post.Uid), UserLevelOperationType.POST_DELETED)
	// 找如果有举报这个帖子的，找举报人加经验
	if len(uids) != 0 {
		for _, log := range uids {
			EditUserLevel(util.AsStrU(log), UserLevelOperationType.REPORT_PASSED)
		}
	}
	return post.Id, nil
}

// 删除帖子记录
func deletePost(post *Post) error {
	/*
		需要删除的内容
		reports
		post_reply
		floors
	*/
	return db.Transaction(func(tx *gorm.DB) error {
		if err := DeleteTagInPost(tx, post.Id); err != nil {
			return err
		}
		if err := deleteReports(tx, "post_id = ?", post.Id); err != nil {
			return err
		}
		// 删除log
		if err := DeletePostReplysInPost(tx, post.Id); err != nil {
			return err
		}
		if err := DeleteFloorsInPost(tx, post.Id); err != nil {
			return err
		}
		if err := tx.Delete(&post, post.Id).Error; err != nil {
			return err
		}
		return nil
	})

}

// 恢复帖子记录
func RecoverPost(postId string) error {
	/*
		需要恢复的内容
		reports
		post_reply
		floors
	*/
	return db.Transaction(func(tx *gorm.DB) error {
		var post Post
		if err := tx.Unscoped().Where("id = ?", postId).Find(&post).Error; err != nil {
			return err
		}
		if err := recoverReports(tx, "post_id = ?", post.Id); err != nil {
			return err
		}
		// 删除log
		if err := RecoverPostReplysInPost(tx, post.Id); err != nil {
			return err
		}
		if err := RecoverFloorsInPost(tx, post.Id); err != nil {
			return err
		}
		if err := tx.Unscoped().Model(&Post{}).Where("id = ?", post.Id).Update("deleted_at", gorm.Expr("NULL")).Error; err != nil {
			return err
		}
		return nil
	})
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
	if err := db.Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新收藏数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		return 0, err
	}
	if err := db.Model(&post).Update("fav_count", post.FavCount+1).Error; err != nil {
		return 0, err
	}

	if uid != util.AsStrU(post.Uid) {
		updatePostTime(post.Id)
		addTagLogInPost(post.Id, TagPointType.FAV_POST)
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
		return 0, err
	}
	if err := db.Model(&post).Update("fav_count", post.FavCount-1).Error; err != nil {
		return 0, err
	}
	if uid != util.AsStrU(post.Uid) {
		addTagLogInPost(post.Id, TagPointType.UNFAV_POST)
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
	if err := db.Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新点赞数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		return 0, err
	}
	if err := db.Model(&post).Update("like_count", post.LikeCount+1).Error; err != nil {
		return 0, err
	}

	if uid != util.AsStrU(post.Uid) {
		updatePostTime(post.Id)
		addTagLogInPost(post.Id, TagPointType.LIKE_POST)
	}
	addUnreadLike(post.Uid, LikeType.POST, post.Id)
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
		return 0, err
	}
	if err := db.Model(&post).Update("like_count", post.LikeCount-1).Error; err != nil {
		return 0, err
	}
	if uid != util.AsStrU(post.Uid) {
		addTagLogInPost(post.Id, TagPointType.UNLIKE_POST)
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
	if err := db.Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新点踩数
	var post Post
	if err := db.Where("id = ?", postId).First(&post).Error; err != nil {
		return 0, err
	}
	if err := db.Model(&post).Update("dis_count", post.DisCount+1).Error; err != nil {
		return 0, err
	}
	if uid != util.AsStrU(post.Uid) {
		updatePostTime(post.Id)
		addTagLogInPost(post.Id, TagPointType.DIS_POST)
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
		return 0, err
	}
	if err := db.Model(&post).Update("dis_count", post.DisCount-1).Error; err != nil {
		return 0, err
	}
	if uid != util.AsStrU(post.Uid) {
		addTagLogInPost(post.Id, TagPointType.UNDIS_POST)
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

func updatePostTime(postId uint64) error {
	return db.Model(&Post{}).Where("id = ?", postId).Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func ReturnPost(uid, postId string) error {
	var post Post
	db.Where("id = ?", postId).Find(&post)
	addManagerLogWithDetail(util.AsUint(uid), util.AsUint(postId), ManagerLogType.POST_RETURN,
		fmt.Sprintf("from %d", post.DepartmentId))
	return db.Model(&Post{}).Where("id = ?", postId).Update("solved", PostSolveType.UNDISTRIBUTED).Error
}
