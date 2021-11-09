package frontend

import (
	"qnhd/middleware/jwt"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

type FrontType int

const (
	User FrontType = iota
	Tag
	Post
	Floor
	History
)

var FrontTypes = [...]FrontType{
	User,
	Tag,
	Post,
	Floor,
	History,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth1)
	g.GET("/auth/:token", RefreshToken)
	g.Use(jwt.JWT(util.USER))
	for _, t := range FrontTypes {
		initType(g, t)
	}
}

func initType(g *gin.RouterGroup, t FrontType) {
	switch t {
	case User:
		//新建用户
		g.POST("/user", AddUsers)
		//修改用户
		g.PUT("/user", EditUsers)
	case Tag:
		//查询标签
		g.GET("/tags", GetTags)
		//新建标签
		g.POST("/tag", AddTag)
		//删除指定标签
		g.DELETE("/tag", DeleteTag)
		//获取热议标签
		g.GET("/tags/hot", GetHotTag)
	case Post:
		//查询多个帖子
		g.GET("/posts", GetPosts)
		//查询个人发帖
		g.GET("/posts/user", GetUserPosts)
		//查询收藏帖子
		g.GET("/posts/fav", GetFavPosts)
		//查询历史帖子
		g.GET("/posts/history", GetHistoryPosts)
		//查询单个帖子
		g.GET("/post", GetPost)
		// 收藏或者取消
		g.POST("/post/favOrUnfav", FavOrUnfavPost)
		// 点赞或者取消
		g.POST("/post/likeOrUnlike", LikeOrUnlikePost)
		// 点踩或者取消
		g.POST("/post/disOrUndis", DisOrUndisPost)
		//新建帖子
		g.POST("/post", AddPost)
		//删除指定帖子
		g.DELETE("/post", DeletePost)
	case Floor:
		//查询多个楼层
		g.GET("/floors", GetFloors)
		//新建楼层
		g.POST("/floor", AddFloor)
		//回复楼层
		g.POST("/floor/reply", ReplyFloor)
		// 点赞或者取消
		g.POST("/floor/likeOrUnlike", LikeOrUnlikeFloor)
		// 点踩或者取消
		g.POST("/floor/disOrUndis", DisOrUndisFloor)
		//删除指定楼层
		g.DELETE("/floor", DeleteFloor)
	}
}
