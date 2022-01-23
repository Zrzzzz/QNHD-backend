package backend

import (
	"qnhd/api/v1/frontend"
	"qnhd/middleware/jwt"
	"qnhd/middleware/permission"
	"qnhd/models"

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
	g.GET("/auth", GetAuth)
	g.GET("/auth/:token", frontend.RefreshToken)
	// 新建用户，不需要token
	g.POST("/user", AddUser)

	g.Use(jwt.JWT())
	g.Use(permission.IdentityDemand(permission.ADMIN))
	for _, t := range BackendTypes {
		initType(g, t)
	}
}

func initType(g *gin.RouterGroup, t BackendType) {
	switch t {
	case Banned:
		bannedGroup := g.Group("", permission.RightDemand(models.UserRight{Super: true}))
		// 获取封禁用户列表
		bannedGroup.GET("/banned", GetBanned)
		// 新建封禁用户
		bannedGroup.POST("/banned", AddBanned)
		// 删除封禁用户
		bannedGroup.GET("/banned/delete", DeleteBanned)
	case Blocked:
		blockedGroup := g.Group("", permission.RightDemand(models.UserRight{Super: true, StuAdmin: true}))
		// 获取禁言用户列表
		blockedGroup.GET("/blocked", GetBlocked)
		// 新建禁言用户
		blockedGroup.POST("/blocked", AddBlocked)
		// 删除指定禁言用户
		blockedGroup.GET("/blocked/delete", DeleteBlocked)
	case Notice:
		noticeGroup := g.Group("", permission.RightDemand(models.UserRight{Super: true}))
		// 获取公告列表
		noticeGroup.GET("/notices", GetNotices)
		// 新建公告
		noticeGroup.POST("/notice", AddNotice)
		// 修改公告
		noticeGroup.POST("/notice/modify", EditNotice)
		// 删除指定公告
		noticeGroup.GET("/notice/delete", DeleteNotice)
	case User:
		// 获取请求者用户信息
		g.GET("/user/info", GetUserInfo)
		// 获取普通用户列表
		g.GET("/users/common", GetCommonUsers)
		// 获取单个普通用户
		g.GET("/user/common", GetCommonUser)
		// 获取管理员列表
		g.GET("/users/manager", permission.RightDemand(models.UserRight{Super: true}), GetManagers)
		// 修改管理员密码
		g.POST("/user/modify/super", permission.RightDemand(models.UserRight{Super: true}), EditUserPasswdBySuper)
		// 修改自己密码
		g.POST("/user/passwd/modify", EditUserPasswd)
		// 修改自己手机
		g.POST("/user/phone/modify", EditUserPhone)
		// 修改用户权限
		g.POST("/user/right/modify", permission.RightDemand(models.UserRight{Super: true}), EditUserRight)
		// 修改用户部门
		g.POST("/user/department/modify", permission.RightDemand(models.UserRight{Super: true}), EditUserDepartment)
	case Post:
		// 获取帖子列表
		g.GET("/posts", GetPosts)
		// 获取帖子
		g.GET("/post", GetPost)
		// 获取帖子回复
		g.GET("/post/replys", frontend.GetPostReplys)
		// 帖子回复校方回应
		g.POST("/post/reply", permission.RightDemand(models.UserRight{Super: true, SchAdmin: true}), AddPostReply)
		// 删除指定帖子
		g.GET("/post/delete", DeletePosts)
	case Report:
		// 获取举报列表
		g.GET("/reports", GetReports)
	case Floor:
		// 查询多个楼层
		g.GET("/floors", GetFloors)
		// 删除指定楼层
		g.GET("/floor/delete", DeleteFloor)
	case Tag:
		// 查询标签
		g.GET("/tags", GetTags)
		// 删除指定标签
		g.GET("/tag/delete", permission.RightDemand(models.UserRight{Super: true, SchAdmin: true}), DeleteTag)
		// 获取热议标签
		g.GET("/tags/hot", GetHotTag)
	case Department:
		// 查询部门
		g.GET("/departments", GetDepartments)
		departGroup := g.Group("", permission.RightDemand(models.UserRight{Super: true}))
		// 添加部门
		departGroup.POST("/department", AddDepartment)
		// 修改部门资料
		departGroup.POST("/department/modify", EditDepartment)
		// 删除部门
		departGroup.GET("/department/delete", DeleteDepartment)
	}
}
