package f

import (
	"qnhd/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type FrontType int

const (
	User FrontType = iota
	Tag
	Post
	Floor
)

var FrontTypes = [...]FrontType{
	User,
	Tag,
	Post,
	Floor,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth)
	g.GET("/auth/:token", RefreshToken)
	g.Use(jwt.JWT())
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
		g.GET("/tag", GetTag)
		//新建标签
		g.POST("/tag", AddTag)
		//删除指定标签
		g.DELETE("/tag", DeleteTag)
	case Post:
		//查询多个帖子
		g.GET("/posts", GetPosts)
		//查询单个帖子
		g.GET("/post", GetPost)
		//新建帖子
		g.POST("/post", AddPosts)
		//删除指定帖子
		g.DELETE("/post", DeletePosts)
	case Floor:
		//查询多个楼层
		g.GET("/floors", GetFloors)
		//新建楼层
		g.POST("/floor", AddFloors)
		//回复楼层
		g.POST("/floor/reply", ReplyFloor)
		//删除指定楼层
		g.DELETE("/floor", DeleteFloor)
	}
}
