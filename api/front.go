package api

import (
	"qnhd/api/v1/f"

	"github.com/gin-gonic/gin"
)

func initHashTagFront(g *gin.RouterGroup) {
	//获取标签列表
	g.GET("/hashtags", f.GetHashTag)
	//新建标签
	g.POST("/hashtags", f.AddHashTag)
	//删除指定标签
	g.DELETE("/hashtags", f.DeleteHashTag)
}

func initUsersFront(g *gin.RouterGroup) {
	//新建用户
	g.POST("/users", f.AddUsers)
	//修改用户
	g.PUT("/users", f.EditUsers)
}

func initAuthFront(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", f.GetAuth)
}
