package api

import (
	"qnhd/api/v1/b"

	"github.com/gin-gonic/gin"
)

func initUsersBackend(g *gin.RouterGroup) {
	//获取用户列表
	g.GET("/users", b.GetUsers)
	//新建用户
	g.POST("/users", b.AddUsers)
	//修改用户
	g.PUT("/users", b.EditUsers)
	//删除指定用户
	g.DELETE("/users", b.DeleteUsers)
}

func initBannedBackend(g *gin.RouterGroup) {
	//获取封禁用户列表
	g.GET("/banned", b.GetBanned)
	//新建封禁用户
	g.POST("/banned", b.AddBanned)
}

func initBlockedBackend(g *gin.RouterGroup) {
	//获取禁言用户列表
	g.GET("/blocked", b.GetBlocked)
	//新建禁言用户
	g.POST("/blocked", b.AddBlocked)
	//删除指定禁言用户
	g.DELETE("/blocked", b.DeleteBlocked)
}

func initAdminBackend(g *gin.RouterGroup) {
	//获取管理员列表
	g.GET("/admin", b.GetAdmins)
	//新建管理员
	g.POST("/admin", b.AddAdmins)
	//修改用户
	g.PUT("/admin", b.EditAdmins)
	//删除指定管理员
	g.DELETE("/admin", b.DeleteAdmins)
}

func initAuthBackend(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", b.GetAuth)
}
