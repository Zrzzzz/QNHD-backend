package models

import (
	"fmt"
	"qnhd/pkg/enums/NoticeType"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/util"
	"qnhd/request/twtservice"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const POST_SCHOOL_TYPE = 1

type Floor struct {
	Model
	Uid         uint64 `json:"uid"`
	PostId      uint64 `json:"post_id"`
	Content     string `json:"content"`
	Nickname    string `json:"nickname" `
	ImageURL    string `json:"image_url" gorm:"default:''"`
	ReplyTo     uint64 `json:"reply_to" gorm:"default:0"`
	ReplyToName string `json:"reply_to_name" gorm:"default:''"`
	SubTo       uint64 `json:"sub_to" gorm:"default:0"`
	LikeCount   uint64 `json:"like_count" gorm:"default:0"`
	DisCount    uint64 `json:"-" gorm:"default:0"`
}

type LogFloorLike struct {
	Uid     uint64 `json:"uid"`
	FloorId uint64 `json:"floor_id"`
}

type LogFloorDis struct {
	Uid     uint64 `json:"uid"`
	FloorId uint64 `json:"floor_id"`
}

// 楼层返回数据
type FloorResponse struct {
	Floor
	SubFloors   []FloorResponse `json:"sub_floors"`
	SubFloorCnt int             `json:"sub_floor_cnt"`
	// 处理链式错误
	Error error `json:"-"`
}

// 客户端返回楼层
type FloorResponseUser struct {
	Floor
	SubFloors   []FloorResponseUser `json:"sub_floors"`
	SubFloorCnt int                 `json:"sub_floor_cnt"`
	IsLike      bool                `json:"is_like"`
	IsDis       bool                `json:"is_dis"`
	IsOwner     bool                `json:"is_owner"`
	// 处理链式错误
	Error error `json:"-"`
}

func (f *Floor) geneResponse(searchSubFloors bool) FloorResponse {
	var fr = FloorResponse{
		Floor:     *f,
		SubFloors: []FloorResponse{},
	}
	if searchSubFloors {
		// 处理回复本条楼层的楼层
		rps, err := getHighlikeSubfloors(util.AsStrU(f.Id))
		if err != nil {
			fr.Error = err
			return fr
		}
		fr.SubFloors = rps
		// 获取子楼层总数
		fr.SubFloorCnt = getFloorSubFloorCount(util.AsStrU(f.Id))
	}
	return fr
}

func (f FloorResponse) searchWithUid(uid string, searchSubFloors bool) FloorResponseUser {
	var fr = FloorResponseUser{
		Floor:       f.Floor,
		IsLike:      IsLikeFloorByUid(uid, util.AsStrU(f.Id)),
		IsDis:       IsDisFloorByUid(uid, util.AsStrU(f.Id)),
		IsOwner:     IsOwnFloorByUid(uid, util.AsStrU(f.Id)),
		SubFloors:   []FloorResponseUser{},
		SubFloorCnt: f.SubFloorCnt,
	}
	if searchSubFloors {
		// 处理回复本条楼层的楼层
		rps, err := getHighlikeSubfloorsWithUid(util.AsStrU(f.Id), uid)
		if err != nil {
			fr.Error = err
			return fr
		}
		fr.SubFloors = rps
	}
	return fr
}

// 将楼层数组转为返回结果数组
func transFloorsToResponses(floor *[]Floor, searchSubFloors bool) ([]FloorResponse, error) {
	var frs = []FloorResponse{}
	var err error
	for _, f := range *floor {
		fr := f.geneResponse(searchSubFloors)
		if fr.Error != nil {
			err = errors.Wrap(err, fr.Error.Error())
		} else {
			frs = append(frs, fr)
		}
	}
	return frs, err
}

// 将楼层数组转为返回结果数组(有uid)
func transFloorsToResponsesWithUid(floor *[]Floor, uid string, searchSubFloors bool) ([]FloorResponseUser, error) {
	var frs = []FloorResponseUser{}
	var err error
	for _, f := range *floor {
		fr := f.geneResponse(searchSubFloors).searchWithUid(uid, searchSubFloors)
		if fr.Error != nil {
			err = errors.Wrap(err, fr.Error.Error())
		} else {
			frs = append(frs, fr)
		}
	}
	return frs, err
}

const OWNER_NAME = "青年湖"

var FLOOR_NAME = []string{
	"9教", "19教", "23教", "26教", "33教", "44教", "45教", "46教", "47教", "规圆楼", "矩方楼", "天麟广场", "北洋广场", "北洋纪念亭", "求实会堂", "天津大学星", "三问桥", "太雷广场", "北洋纪念林", "鹏翔公寓", "克拉公寓", "留园", "格园", "诚园", "正园", "修园", "治园", "平园", "知园", "梅园餐厅", "兰园餐厅", "桃园餐厅", "棠园餐厅", "竹园餐厅", "留园餐厅", "青园餐厅", "菊园餐厅", "天大美食广场", "校训石", "尚贤石", "敬业湖", "友谊湖", "爱晚湖", "1895行政楼", "郑东图书馆", "科学图书馆", "海棠书屋", "土立方", "大学生活动中心", "体育馆", "游泳馆", "土木馆", "水利馆", "校史馆", "冯骥才艺术研究院", "王学仲艺术研究所", "天南楼", "北洋门诊部", "北洋超市", "罗森便利店", "京东便利店", "小诚食", "菜鸟驿站", "天大四季村", "天津大学幼儿园", "斗兽场", "求实影院", "海小棠", "洗衣房", "天大纪念品店", "理发店", "修车铺", "隔壁南开", "地科院",
}

func getFloorNickname(cnt int) string {
	idx := cnt % len(FLOOR_NAME)
	if cnt/len(FLOOR_NAME) >= 1 {
		return fmt.Sprintf("%s%d号", FLOOR_NAME[idx], cnt/len(FLOOR_NAME)+1)
	}
	return FLOOR_NAME[idx]
}

// 根据id返回
func GetFloor(floorId string) (Floor, error) {
	var floor Floor
	err := db.Where("id = ?", floorId).First(&floor).Error
	return floor, err
}

// 返回单个楼层带Response
func GetFloorResponse(floorId string) (FloorResponse, error) {
	var ret FloorResponse
	floor, err := GetFloor(floorId)
	if err != nil {
		return ret, err
	}
	fr := floor.geneResponse(true)
	return fr, fr.Error
}

// 返回单个楼层带Response,有uid
func GetFloorResponseWithUid(floorId, uid string) (FloorResponseUser, error) {
	var ret FloorResponseUser
	floor, err := GetFloor(floorId)
	if err != nil {
		return ret, err
	}
	fr := floor.geneResponse(true).searchWithUid(uid, true)
	return fr, fr.Error
}

// // 缩略返回帖子内楼层，即返回5条
// func getShortFloorResponsesInPost(postId string) ([]FloorResponse, error) {
// 	var floors []Floor
// 	if err := db.Where("post_id = ? AND reply_to = 0", postId).Order("created_at").Limit(5).Find(&floors).Error; err != nil {
// 		return nil, err
// 	}
// 	return transFloorsToResponses(&floors, true)
// }

// // 缩略返回帖子内楼层，即返回5条 含用户id
// func getShortFloorResponsesInPostWithUid(postId, uid string) ([]FloorResponseUser, error) {
// 	var floors []Floor
// 	if err := db.Where("post_id = ? AND reply_to = 0", postId).Order("created_at").Limit(5).Find(&floors).Error; err != nil {
// 		return nil, err
// 	}
// 	return transFloorsToResponsesWithUid(&floors, uid, true)
// }

// 分页返回帖子里的楼层
func GetFloorResponses(c *gin.Context, postId string) ([]FloorResponse, error) {
	var floors []Floor
	if err := db.Where("post_id = ? AND reply_to = 0", postId).Order("created_at DESC").Scopes(util.Paginate(c)).Find(&floors).Error; err != nil {
		return nil, err
	}
	return transFloorsToResponses(&floors, true)
}

// 分页返回帖子里的楼层，带uid
func GetFloorResponsesWithUid(c *gin.Context, postId, uid string) ([]FloorResponseUser, error) {
	var floors []Floor
	if err := db.Where("post_id = ? AND reply_to = 0", postId).Order("created_at DESC").Scopes(util.Paginate(c)).Find(&floors).Error; err != nil {
		return nil, err
	}
	return transFloorsToResponsesWithUid(&floors, uid, true)
}

// 返回楼层内最高赞的5条楼层
func getHighlikeSubfloors(floorId string) ([]FloorResponse, error) {
	var floors []Floor
	// 按照点赞降序，创建时间降序
	err := db.Where("sub_to = ?", floorId).Order("like_count DESC, created_at DESC").Limit(5).Find(&floors).Error
	if err != nil {
		return []FloorResponse{}, err
	}
	return transFloorsToResponses(&floors, false)
}

// 返回楼层内最高赞的5条楼层，携带用户id
func getHighlikeSubfloorsWithUid(floorId, uid string) ([]FloorResponseUser, error) {
	var floors []Floor
	// 按照点赞降序，创建时间降序
	err := db.Where("sub_to = ?", floorId).Order("like_count DESC, created_at DESC").Limit(5).Find(&floors).Error
	if err != nil {
		return []FloorResponseUser{}, err
	}
	return transFloorsToResponsesWithUid(&floors, uid, false)
}

func getFloorSubFloorCount(floorId string) int {
	var ret int64
	db.Model(&Floor{}).Where("sub_to = ?", floorId).Count(&ret)
	return int(ret)
}

func GetCommentCount(postId uint64, withSubfloors bool) int {
	var ret int64
	a := db.Model(&Floor{}).Where("post_id = ?", postId)
	if !withSubfloors {
		a = a.Where("sub_to = 0")
	}
	a.Count(&ret)
	return int(ret)
}

// 分页返回楼层内的回复
func GetFloorReplyResponses(c *gin.Context, floorId string) ([]FloorResponse, error) {
	var floors []Floor
	err := db.Where("sub_to = ?", floorId).Order("created_at").Scopes(util.Paginate(c)).Find(&floors).Error
	if err != nil {
		return []FloorResponse{}, err
	}
	return transFloorsToResponses(&floors, false)
}

// 分页返回楼层内的回复带uid
func GetFloorReplyResponsesWithUid(c *gin.Context, floorId, uid string) ([]FloorResponseUser, error) {
	var floors []Floor
	err := db.Where("sub_to = ?", floorId).Order("created_at").Scopes(util.Paginate(c)).Find(&floors).Error
	if err != nil {
		return []FloorResponseUser{}, err
	}
	return transFloorsToResponsesWithUid(&floors, uid, false)
}

// 添加楼层评论
func AddFloor(maps map[string]interface{}) (uint64, error) {
	var post Post
	var nickname string
	uid := maps["uid"].(uint64)
	postId := maps["postId"].(uint64)
	// 先找到post主人
	if err := db.First(&post, postId).Error; err != nil {
		return 0, err
	}
	// 如果是帖子主人
	if post.Uid == uid {
		nickname = OWNER_NAME
	} else {
		var floor Floor
		if err := db.Unscoped().Where("uid = ? AND post_id = ?", uid, postId).Find(&floor).Error; err != nil {
			return 0, err
		}
		// 是否已经发过言
		if floor.Id > 0 {
			nickname = floor.Nickname
		} else {
			var cnt int64
			// 除去owner, 这里加上被删除的楼
			if err := db.Model(&Floor{}).Unscoped().Where("post_id = ? AND uid <> ?", postId, post.Uid).Distinct("uid").Count(&cnt).Error; err != nil {
				return 0, err
			}
			nickname = getFloorNickname(int(cnt))
		}
	}
	var newFloor = Floor{
		Uid:      uid,
		PostId:   postId,
		Content:  filter.Filter(maps["content"].(string)),
		Nickname: nickname,
		ImageURL: maps["image_url"].(string),
	}
	if err := db.Create(&newFloor).Error; err != nil {
		return 0, err
	}
	// 如果不是回复自己的帖子，通知帖子主人
	if post.Uid != uid {
		addUnreadFloor(post.Uid, newFloor.Id)
		var user User
		if err := db.Where("id = ?", post.Uid).Find(&user).Error; err == nil {
			if err := twtservice.NotifyPost(post.Title, user.Number); err != nil {
				logging.Error(err.Error())
			}
		}
	}
	// 对帖子的tag增加记录, 当不是校务才会有
	if post.Type != POST_SCHOOL_TYPE {
		addTagLogInPost(post.Id, TAG_ADD_FLOOR)
	}
	updatePostTime(post.Id)
	return newFloor.Id, nil
}

// 添加楼层回复
func ReplyFloor(maps map[string]interface{}) (uint64, error) {
	var post Post
	var nickname string
	uid := maps["uid"].(uint64)
	// 判断存在floor
	floorId := maps["replyToFloor"].(uint64)
	var toFloor Floor
	if err := db.First(&toFloor, floorId).Error; err != nil {
		return 0, err
	}
	postId := toFloor.PostId
	// 先找到post主人
	if err := db.First(&post, postId).Error; err != nil {
		return 0, err
	}

	if post.Uid == uid {
		nickname = OWNER_NAME
	} else {
		// 还有可能已经发过言
		var floor Floor
		if err := db.Unscoped().Where("uid = ? AND post_id = ?", uid, postId).Find(&floor).Error; err != nil {
			return 0, err
		}
		if floor.Id > 0 {
			nickname = floor.Nickname
		} else {
			var cnt int64
			// 除去owner
			if err := db.Model(&Floor{}).Where("post_id = ? AND uid <> ?", postId, post.Uid).Distinct("uid").Count(&cnt).Error; err != nil {
				return 0, err
			}
			nickname = getFloorNickname(int(cnt))
		}
	}

	var newFloor = Floor{
		Uid:         uid,
		PostId:      toFloor.PostId,
		Content:     maps["content"].(string),
		Nickname:    nickname,
		ImageURL:    maps["image_url"].(string),
		ReplyTo:     toFloor.Id,
		ReplyToName: toFloor.Nickname,
	}
	// 判断子楼层
	// 如果没有subto，说明回复的不是子楼层
	if toFloor.SubTo == 0 {
		newFloor.SubTo = toFloor.Id
	} else {
		newFloor.SubTo = toFloor.SubTo
	}

	if err := db.Create(&newFloor).Error; err != nil {
		return 0, err
	}

	// 如果不是回复自己的帖子，通知帖子主人
	if post.Uid != uid {
		addUnreadFloor(post.Uid, newFloor.Id)
	}
	// 如果回复的楼层不是子楼层，通知回复的楼层的主人，这里开始避免重复
	if toFloor.Uid != uid && toFloor.Uid != post.Uid {
		addUnreadFloor(toFloor.Uid, newFloor.Id)
		var user User
		if err := db.Where("id = ?", toFloor.Uid).Find(&user).Error; err == nil {

			if err := twtservice.NotifyFloor(toFloor.Content, user.Number); err != nil {
				logging.Error(err.Error())
			}
		}
	}
	// 如果回复的帖子是子楼层，通知层主
	if toFloor.SubTo != 0 {
		subToFloor, _ := GetFloor(util.AsStrU(newFloor.SubTo))
		if subToFloor.Uid != uid && subToFloor.Uid != toFloor.Uid && subToFloor.Uid != post.Uid {
			addUnreadFloor(subToFloor.Uid, newFloor.Id)
		}
	}

	// 对帖子的tag增加记录, 当不是校务才会有
	if post.Type != POST_SCHOOL_TYPE {
		addTagLogInPost(post.Id, TAG_ADD_FLOOR)
	}

	updatePostTime(post.Id)

	return newFloor.Id, nil
}

func DeleteFloorByAdmin(uid, floorId string) (uint64, error) {
	var floor Floor
	var post Post
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		return 0, err
	}
	if err := db.Where("id = ?", floor.PostId).Find(&post).Error; err != nil {
		return 0, err
	}

	if err := deleteFloor(&floor); err != nil {
		return 0, err
	}

	updatePostTime(floor.PostId)
	// 通知举报过楼层的所有用户
	var uids []uint64
	db.Model(&Report{}).Select("uid").Where("type = ? AND floor_id = ?", ReportTypeFloor, floor.Id).Find(&uids)
	addNoticeWithTemplate(NoticeType.FLOOR_REPORT_SOLVE, uids, []string{post.Title, floor.Content})
	// 通知被删除的用户
	addNoticeWithTemplate(NoticeType.FLOOR_DELETED, []uint64{floor.Uid}, []string{post.Title, floor.Content})
	return floor.Id, nil
}

func DeleteFloorByUser(uid, floorId string) (uint64, error) {
	var floor = Floor{}
	if err := db.Where("uid = ? AND id = ?", uid, floorId).First(&floor).Error; err != nil {
		return 0, err
	}
	if err := deleteFloor(&floor); err != nil {
		return 0, err
	}

	updatePostTime(floor.PostId)

	return floor.Id, nil
}

func deleteFloor(floor *Floor) error {
	/*
		删除楼层逻辑
		subto的帖子, reply_to的帖子
		log_floor_dis, log_floor_like, log_unread_floor
		reports
	*/
	return db.Transaction(func(tx *gorm.DB) error {

		// 先找到所有楼层
		var (
			floors        = map[uint64]bool{}
			ids           []uint64
			subToFloors   []Floor
			replyToFloors []Floor
		)
		if err := db.Where("sub_to = ?", floor.Id).Find(&subToFloors).Error; err != nil {
			return err
		}
		if err := db.Where("reply_to = ?", floor.Id).Find(&replyToFloors).Error; err != nil {
			return err
		}
		// 这里需要避免重复, 合并到floors里
		for _, f := range subToFloors {
			_, ok := floors[f.Id]
			if !ok {
				floors[f.Id] = true
			}
		}
		for _, f := range replyToFloors {
			_, ok := floors[f.Id]
			if !ok {
				floors[f.Id] = true
			}
		}
		for k := range floors {
			ids = append(ids, k)
		}
		// 加上自己
		ids = append(ids, floor.Id)
		// 删除log
		if err := tx.Where("floor_id IN (?)", ids).Delete(&LogFloorLike{}).Error; err != nil {
			return err
		}
		if err := tx.Where("floor_id IN (?)", ids).Delete(&LogFloorDis{}).Error; err != nil {
			return err
		}
		// 删除reports
		if err := tx.Where("id IN (?) AND type = ?", ids, LIKE_FLOOR).Delete(&LogUnreadLike{}).Error; err != nil {
			return err
		}
		if err := deleteReports(tx, "floor_id IN (?)", ids); err != nil {
			return err
		}

		return db.Delete(&Floor{}, ids).Error
	})

}

func DeleteFloorsInPost(tx *gorm.DB, postId uint64) error {
	if tx == nil {
		tx = db
	}
	var floorIds []uint64
	if err := tx.Model(&Floor{}).Select("id").Where("post_id = ?", postId).Find(&floorIds).Error; err != nil {
		return err
	}
	// 楼层回复记录
	if err := tx.Where("floor_id IN (?)", floorIds).Delete(&LogUnreadFloor{}).Error; err != nil {
		return err
	}
	// 楼层点赞记录
	if err := tx.Where("id IN (?) AND type = ?", floorIds, LIKE_FLOOR).Delete(&LogUnreadLike{}).Error; err != nil {
		return err
	}
	if err := tx.Where("post_id = ?", postId).Delete(&Floor{}).Error; err != nil {
		return err
	}
	return nil
}

// 点赞楼层
func LikeFloor(floorId string, uid string) (uint64, error) {
	var log LogFloorLike
	// 首先判断点没点过赞
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid > 0 {
		return 0, fmt.Errorf("已被点赞")
	}

	log.Uid = util.AsUint(uid)
	log.FloorId = util.AsUint(floorId)
	if err := db.Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("like_count", floor.LikeCount+1).Error; err != nil {
		return 0, err
	}

	updatePostTime(floor.PostId)
	addUnreadLike(floor.Uid, LIKE_FLOOR, floor.Id)
	UndisFloor(floorId, uid)
	addTagLogInPost(floor.PostId, TAG_LIKE_FLOOR)
	return floor.LikeCount, nil
}

// 取消点赞楼层
func UnlikeFloor(floorId string, uid string) (uint64, error) {
	var log LogFloorLike
	// 首先判断点没点过赞
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		return 0, err
	}

	if log.Uid == 0 {
		return 0, fmt.Errorf("未被点赞")
	}

	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Delete(&log).Error; err != nil {
		return 0, err
	}

	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("like_count", floor.LikeCount-1).Error; err != nil {
		return 0, err
	}
	addTagLogInPost(floor.PostId, TAG_UNLIKE_FLOOR)
	return floor.LikeCount, nil
}

// 点踩楼层
func DisFloor(floorId string, uid string) (uint64, error) {
	var log LogFloorDis
	// 首先判断点没点过踩
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid > 0 {
		return 0, fmt.Errorf("已被点踩")
	}

	log.Uid = util.AsUint(uid)
	log.FloorId = util.AsUint(floorId)
	if err := db.Create(&log).Error; err != nil {
		return 0, err
	}
	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("dis_count", floor.DisCount+1).Error; err != nil {
		return 0, err
	}
	UnlikeFloor(floorId, uid)
	addTagLogInPost(floor.PostId, TAG_DIS_FLOOR)
	return floor.DisCount, nil
}

func UndisFloor(floorId string, uid string) (uint64, error) {
	var log LogFloorDis
	// 首先判断点没点过踩
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		return 0, err
	}
	if log.Uid == 0 {
		return 0, fmt.Errorf("未被点踩")
	}

	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Delete(&log).Error; err != nil {
		return 0, err
	}

	// 更新楼的likes
	var floor Floor
	if err := db.Where("id = ?", floorId).First(&floor).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
	}
	if err := db.Model(&floor).Update("dis_count", floor.DisCount-1).Error; err != nil {
		return 0, err
	}
	addTagLogInPost(floor.PostId, TAG_UNDIS_FLOOR)
	return floor.DisCount, nil
}

func IsLikeFloorByUid(uid, floorId string) bool {
	var log LogFloorLike
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Uid > 0
}

func IsDisFloorByUid(uid, floorId string) bool {
	var log LogFloorDis
	if err := db.Where("uid = ? AND floor_id = ?", uid, floorId).Find(&log).Error; err != nil {
		logging.Error(err.Error())
		return false
	}
	return log.Uid > 0
}

func IsOwnFloorByUid(uid, floorId string) bool {
	var floor, err = GetFloor(floorId)
	if err != nil {
		return false
	}
	return util.AsStrU(floor.Uid) == uid
}
