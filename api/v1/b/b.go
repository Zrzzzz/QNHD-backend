package b

import (
	"qnhd/middleware/jwt"
	"qnhd/pkg/util"

	"github.com/gin-gonic/gin"
)

type BackendType int

const (
	Admin BackendType = iota
	Banned
	Blocked
	Notice
	User
	Post
	Report
)

var BackendTypes = [...]BackendType{
	Admin,
	Banned,
	Blocked,
	Notice,
	User,
	Post,
	Report,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth)
	g.GET("/auth/:token", RefreshToken)
	g.Use(jwt.JWT(util.ADMIN))
	for _, t := range BackendTypes {
		initType(g, t)
	}
}

func initType(g *gin.RouterGroup, t BackendType) {
	switch t {
	case Admin:
		//获取管理员列表
		g.GET("/admin", GetAdmins)
		//新建管理员
		g.POST("/admin", AddAdmins)
		//修改管理员
		g.PUT("/admin", EditAdmins)
		//删除指定管理员
		g.DELETE("/admin", DeleteAdmins)
	case Banned:
		//获取封禁用户列表
		g.GET("/banned", GetBanned)
		//新建封禁用户
		g.POST("/banned", AddBanned)
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
		g.POST("/notice", AddNotices)
		//修改公告
		g.PUT("/notice", EditNotices)
		//删除指定公告
		g.DELETE("/notice", DeleteNotices)
	case User:
		//获取用户列表
		g.GET("/user", GetUsers)
		//新建用户
		g.POST("/user", AddUsers)
		//修改用户
		g.PUT("/user", EditUsers)
		//删除指定用户
		g.DELETE("/user", DeleteUsers)
	case Post:
		//获取帖子列表
		g.GET("/post", GetPosts)
		//删除指定帖子
		g.DELETE("/post", DeletePosts)
	case Report:
		//获取帖子列表
		g.GET("/report", GetReports)
	}
}
