package f

import (
	"qnhd/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type FrontType int

const (
	User FrontType = iota
	Tag
)

var FrontTypes = [...]FrontType{
	User,
	Tag,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth)
	g.Use(jwt.JWT())
	g.GET("/auth/:token", RefreshToken)
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
	}
}
