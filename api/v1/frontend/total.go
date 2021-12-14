package frontend

import (
	"qnhd/middleware/jwt"

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
)

var FrontTypes = [...]FrontType{
	Tag,
	Post,
	Floor,
	History,
	Department,
	Report,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth)
	g.GET("/auth/:token", RefreshToken)
	g.Use(jwt.JWT(jwt.USER))
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
	case Post:
		// 查询多个帖子
		g.GET("/posts", GetPosts)
		// 查询个人发帖
		g.GET("/posts/user", GetUserPosts)
		// 查询收藏帖子
		g.GET("/posts/fav", GetFavPosts)
		// 查询历史帖子
		g.GET("/posts/history", GetHistoryPosts)
		// 查询单个帖子
		g.GET("/post", GetPost)
		//  收藏或者取消
		g.POST("/post/favOrUnfav/modify", FavOrUnfavPost)
		//  点赞或者取消
		g.POST("/post/likeOrUnlike/modify", LikeOrUnlikePost)
		//  点踩或者取消
		g.POST("/post/disOrUndis/modify", DisOrUndisPost)
		// 新建帖子
		g.POST("/post", AddPost)
		// 删除指定帖子
		g.GET("/post/delete", DeletePost)
	case Floor:
		// 查询多个楼层
		g.GET("/floors", GetFloors)
		// 查询楼层内回复
		g.GET("/f/floor/replys", GetFloorReplys)
		// 新建楼层
		g.POST("/floor", AddFloor)
		// 回复楼层
		g.POST("/floor/reply", ReplyFloor)
		//  点赞或者取消
		g.POST("/floor/likeOrUnlike/modify", LikeOrUnlikeFloor)
		//  点踩或者取消
		g.POST("/floor/disOrUndis/modify", DisOrUndisFloor)
		// 删除指定楼层
		g.GET("/floor/delete", DeleteFloor)
	case Department:
		// 查询部门
		g.GET("/departments", GetDepartments)
	case Report:
		g.POST("/report", AddReport)
	}
}
