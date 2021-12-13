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
	Tag
	Department
)

var BackendTypes = [...]BackendType{
	Banned,
	Blocked,
	Notice,
	User,
	Post,
	Report,
	Floor,
	Tag,
	Department,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	// 昵称认证
	g.GET("/auth/number", GetAuthNumber)
	// 电话认证
	g.GET("/auth/phone", GetAuthPhone)
	g.GET("/auth/:token", frontend.RefreshToken)
	// 新建用户，不需要token
	g.POST("/user", AddUser)

	g.Use(jwt.JWT(jwt.ADMIN))
	for _, t := range BackendTypes {
		initType(g, t)
	}
}

func initType(g *gin.RouterGroup, t BackendType) {
	switch t {
	case Banned:
		// 获取封禁用户列表
		g.GET("/banned", GetBanned)
		// 新建封禁用户
		g.POST("/banned", AddBanned)
		// 删除封禁用户
		g.GET("/banned/delete", DeleteBanned)
	case Blocked:
		// 获取禁言用户列表
		g.GET("/blocked", GetBlocked)
		// 新建禁言用户
		g.POST("/blocked", AddBlocked)
		// 删除指定禁言用户
		g.GET("/blocked/delete", DeleteBlocked)
	case Notice:
		// 获取公告列表
		g.GET("/notice", GetNotices)
		// 新建公告
		g.POST("/notice", AddNotice)
		// 修改公告
		g.POST("/notice/modify", EditNotice)
		// 删除指定公告
		g.GET("/notice/delete", DeleteNotice)
	case User:
		// 获取普通用户列表
		g.GET("/users/common", GetCommonUsers)
		// 获取所有用户列表
		g.GET("/users/all", GetAllUsers)
		// 修改用户密码
		g.POST("/user/modify", EditUser)
		// 修改用户权限
		g.POST("/user/right/modify", EditUserRight)
		// 修改用户部门
		g.POST("/user/department/modify", EditUserDepartment)
	case Post:
		// 获取帖子列表
		g.GET("/posts", frontend.GetPosts)
		// 获取帖子
		g.GET("/post", frontend.GetPost)
		// 删除指定帖子
		g.GET("/post/delete", DeletePosts)
	case Report:
		// 获取帖子列表
		g.GET("/report", GetReports)
	case Floor:
		// 查询多个楼层
		g.GET("/floors", GetFloors)
		// 删除指定楼层
		g.GET("/floor/delete", DeleteFloor)
	case Tag:
		// 查询标签
		g.GET("/tags", GetTags)
		// 删除指定标签
		g.GET("/tag/delete", DeleteTag)
		// 获取热议标签
		g.GET("/tags/hot", GetHotTag)
	case Department:
		// 查询部门
		g.GET("/departments", GetDepartments)
		// 添加部门
		g.POST("/department", AddDepartment)
		// 修改部门资料
		g.POST("/department/modify", EditDepartment)
		// 删除部门
		g.GET("/department/delete", DeleteDepartment)
	}
}
