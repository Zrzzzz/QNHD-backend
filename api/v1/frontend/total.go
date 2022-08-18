package frontend

import (
	"qnhd/enums/IdentityType"
	"qnhd/middleware/jwt"
	"qnhd/middleware/permission"

	"github.com/gin-gonic/gin"
)

type FrontType int

const (
	Tag FrontType = iota
	Post
	Floor
	History
	Department
	Report
	Message
	Game
	PostType
	Banner
	User
	Share
)

var FrontTypes = [...]FrontType{
	Tag,
	Post,
	Floor,
	History,
	Department,
	Report,
	Message,
	Game,
	PostType,
	Banner,
	User,
	Share,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth/passwd", GetAuthPasswd)
	g.GET("/auth/token", GetAuthToken)
	g.GET("/auth/:token", RefreshToken)
	g.Use(jwt.JWT())
	g.Use(permission.IdentityDemand(IdentityType.USER))
	// 封号的话不能访问
	g.Use(permission.ValidBanned())
	for _, t := range FrontTypes {
		initType(g, t)
	}
}

func initType(g *gin.RouterGroup, t FrontType) {
	switch t {
	case Tag:
		// 查询标签
		g.GET("/tags", GetTags)
		// 新建标签
		g.POST("/tag", AddTag)
		// 删除指定标签
		g.GET("/tag/delete", DeleteTag)
		// 获取热议标签
		g.GET("/tags/hot", GetHotTag)
		// 获取推荐标签
		g.GET("/tag/recommend", GetRecommendTag)
	case Post:
		// 查询多个帖子
		g.GET("/posts", GetPosts())
		// 查询个人发帖
		g.GET("/posts/user", GetUserPosts)
		// 查询收藏帖子
		g.GET("/posts/fav", GetFavPosts)
		// 查询历史帖子
		g.GET("/posts/history", GetHistoryPosts)
		// 查询单个帖子
		g.GET("/post", GetPost())
		// 新建帖子
		g.POST("/post", permission.ValidBlocked(), AddPost)
		// 解决问题
		g.POST("/post/solve", EditPostSolved)
		// 获取帖子回复
		g.GET("/post/replys", GetPostReplys)
		// 帖子回复校方回应
		g.POST("/post/reply", AddPostReply)
		// 收藏或者取消
		g.POST("/post/fav", FavOrUnfavPost)
		// 点赞或者取消
		g.POST("/post/like", LikeOrUnlikePost)
		// 点踩或者取消
		g.POST("/post/dis", DisOrUndisPost)
		// 访问记录
		g.POST("/post/visit", VisitPost)
		// 删除指定帖子
		g.GET("/post/delete", DeletePost)
	case Floor:
		// 查询多个楼层
		g.GET("/floors", GetFloors)
		// 查询单个楼层
		g.GET("/floor", GetFloor)
		// 查询楼层内回复
		g.GET("/floor/replys", GetFloorReplys)
		// 新建楼层
		g.POST("/floor", permission.ValidBlocked(), AddFloor)
		// 回复楼层
		g.POST("/floor/reply", permission.ValidBlocked(), ReplyFloor)
		//  点赞或者取消
		g.POST("/floor/like", LikeOrUnlikeFloor)
		//  点踩或者取消
		g.POST("/floor/dis", DisOrUndisFloor)
		// 删除指定楼层
		g.GET("/floor/delete", DeleteFloor)
	case Department:
		// 查询部门
		g.GET("/departments", GetDepartments)
	case Report:
		// 添加举报
		g.POST("/report", AddReport)
	case Message:
		// 获取未读楼层
		g.GET("/message/floors", GetMessageFloors)
		// 获取未读回复
		g.GET("/message/replys", GetMessagePostReplys)
		// 获取未读通知
		g.GET("/message/notices", GetMessageNotices)
		// 获取管理员通知
		g.GET("/message/notices/department", GetMessageDepartmentNotices)
		// 获取未读点赞
		g.GET("/message/likes", GetMessageLikes)
		// 获取未读数量
		g.GET("/message/count", GetMessageCount)
		// 已读通知
		g.POST("/message/notice/read", ReadNotice)
		// 删除通知记录
		g.GET("/message/notices/delete", DeleteMessageNotices)
		// 已读楼层
		g.POST("/message/floor/read", ReadFloor)
		// 已读评论内楼层
		g.POST("/message/floor/read_in_post", ReadFloorInPost)
		// 已读回复
		g.POST("/message/reply/read", ReadReply)
		// 已读点赞
		g.POST("/message/like/read", ReadLike)
		// 全部已读
		g.POST("/message/all", ReadAllMessage)
	case Game:
		// 获取游戏列表
		g.GET("/game", GetNewestGame)
	case PostType:
		// 获取帖子类型
		g.GET("/posttypes", GetPostTypes)
	case Banner:
		// 获取游戏列表
		g.GET("/banners", GetBanners)
	case User:
		// 获取自己信息
		g.GET("/user", GetUserInfo)
		// 修改昵称
		g.POST("/user/name", EditUserName)
		// 修改头像
		g.POST("/user/avatar", EditUserAvatar)
	case Share:
		// 分享记录
		g.POST("/share", ShareLog)
	}
}
