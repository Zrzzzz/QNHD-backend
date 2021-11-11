package backend

import (
	"qnhd/api/v1/frontend"
	"qnhd/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type BackendType int

const (
	Banned BackendType = iota
	Blocked
	Notice
	User
	Post
	Report
	Floor
)

var BackendTypes = [...]BackendType{
	Banned,
	Blocked,
	Notice,
	User,
	Post,
	Report,
	Floor,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth)
	g.GET("/auth/:token", frontend.RefreshToken)
	//新建用户，不需要token
	g.POST("/user", AddUser)

	g.Use(jwt.JWT(jwt.ADMIN))
	for _, t := range BackendTypes {
		initType(g, t)
	}
}

func initType(g *gin.RouterGroup, t BackendType) {
	switch t {
	case Banned:
		//获取封禁用户列表
		g.GET("/banned", GetBanned)
		//新建封禁用户
		g.POST("/banned", AddBanned)
		//删除封禁用户
		g.DELETE("/banned", DeleteBanned)
	case Blocked:
		//获取禁言用户列表
		g.GET("/blocked", GetBlocked)
		//新建禁言用户
		g.POST("/blocked", AddBlocked)
		//删除指定禁言用户
		g.DELETE("/blocked", DeleteBlocked)
	case Notice:
		//获取公告列表
		g.GET("/notice", GetNotices)
		//新建公告
		g.POST("/notice", AddNotice)
		//修改公告
		g.PUT("/notice", EditNotice)
		//删除指定公告
		g.DELETE("/notice", DeleteNotice)
	case User:
		//获取普通用户列表
		g.GET("/users/common", GetCommonUsers)
		//获取所有用户列表
		g.GET("/users/all", GetAllUsers)
		//修改用户密码
		g.PUT("/user", EditUser)
		//修改用户权限
		g.PUT("/user/right", EditUserRight)
	case Post:
		//获取帖子列表
		g.GET("/posts", frontend.GetPosts)
		//获取帖子
		g.GET("/post", frontend.GetPost)
		//删除指定帖子
		g.DELETE("/post", DeletePosts)
	case Report:
		//获取帖子列表
		g.GET("/report", GetReports)
	case Floor:
		//查询多个楼层
		g.GET("/floors", GetFloors)
		//删除指定楼层
		g.DELETE("/floor", DeleteFloor)
	}
}
