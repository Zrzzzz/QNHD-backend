package api

import (
	v1 "qnhd/api/v1"
	"qnhd/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() (r *gin.Engine) {
	r = gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)
	apibackv1 := r.Group("/api/b/v1")
	{
		initUsers(apibackv1)
		initBanned(apibackv1)
	}
	apifrontv1 := r.Group("api/f/v1")
	{
		initHashTag(apifrontv1)
		initUsers(apifrontv1)
	}

	return r
}

func initHashTag(g *gin.RouterGroup) {
	//获取标签列表
	g.GET("/hashtags", v1.GetHashTag)
	//新建标签
	g.POST("/hashtags", v1.AddHashTag)
	//删除指定标签
	g.DELETE("/hashtags/:id", v1.DeleteHashTag)
}

func initUsers(g *gin.RouterGroup) {
	//获取用户列表
	g.GET("/users", v1.GetUsers)
	//新建用户
	g.POST("/users", v1.AddUsers)
	//修改用户
	g.PUT("/users", v1.EditUsers)
	//删除指定用户
	g.DELETE("/users", v1.DeleteUsers)
}

func initBanned(g *gin.RouterGroup) {
	//获取禁言用户列表
	g.GET("/banned", v1.GetBanned)
	//新建禁言用户
	g.POST("/banned", v1.AddBanned)
	//修改禁言用户
	g.PUT("/banned/:id", v1.EditBanned)
	//删除指定禁言用户
	g.DELETE("/banned/:id", v1.DeleteBanned)
}
