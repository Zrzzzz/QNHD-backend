package backend

import (
	"qnhd/api/v1/frontend"
	"qnhd/enums/IdentityType"
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
	Game
	Sensitive
	PostType
	Banner
	Statistic
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
	Game,
	Sensitive,
	PostType,
	Banner,
	Statistic,
}

func Setup(g *gin.RouterGroup) {
	// 获取token
	g.GET("/auth", GetAuth)
	g.GET("/auth/:token", frontend.RefreshToken)
	g.Use(jwt.JWT())
	g.Use(permission.IdentityDemand(IdentityType.ADMIN))
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
		// 新建公告模板
		noticeGroup.POST("/notice/template", AddNoticeTemplate)
		// 修改公告
		noticeGroup.POST("/notice/modify", EditNoticeTemplate)
		// 删除指定公告
		noticeGroup.GET("/notice/delete", DeleteNotice)
	case User:
		// 新建单个用户
		g.POST("/user", permission.RightDemand(models.UserRight{Super: true}), AddUser)
		// 新建多个用户
		g.POST("/users", permission.RightDemand(models.UserRight{Super: true}), AddUsers)
		// 获取某用户详细信息
		g.GET("/user/detail", permission.RightDemand(models.UserRight{Super: true}), GetUserDetail)
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
		// 删除管理员
		g.GET("/user/manager/delete", permission.RightDemand(models.UserRight{Super: true}), DeleteManager)
		// 强制更新token
		g.GET("/user/update", permission.RightDemand(models.UserRight{Super: true}), ForceTokenUpdate)
	case Post:
		// 获取帖子列表
		g.GET("/posts", frontend.GetPosts(false))
		// 获取用户帖子
		g.GET("/posts/user", GetUserPosts)
		// 获取帖子
		g.GET("/post", frontend.GetPost(false))
		// 获取帖子回复
		g.GET("/post/replys", frontend.GetPostReplys)
		// 帖子回复校方回应
		g.POST("/post/reply", permission.RightDemand(models.UserRight{Super: true, SchAdmin: true}), AddPostReply)
		// 帖子转移部门
		g.POST("/post/transfer/department", permission.RightDemand(models.UserRight{Super: true, SchAdmin: true}), TransferPostDepartment)
		// 帖子换类型
		g.POST("/post/transfer/type", TransferPostType)
		// 修改帖子加精值
		g.POST("/post/value", permission.RightDemand(models.UserRight{Super: true, StuAdmin: true}), EditPostValue)
		// 删除指定帖子
		g.GET("/post/delete", permission.RightDemand(models.UserRight{Super: true, StuAdmin: true}), DeletePost)
		// 恢复指定帖子
		g.POST("/post/recover", permission.RightDemand(models.UserRight{Super: true}), RecoverPost)
		// 添加帖子标签
		g.POST("/post_tag", AddPostTag)
		// 删除帖子的标签
		g.GET("/post_tag/delete", DeletePostTag)
		// 删除帖子的图片
		g.GET("/post_image/delete", DeletePostImages)
	case Report:
		// 获取举报列表
		g.GET("/reports", GetReports)
		// 删除举报
		g.GET("/report/delete", SolveReport)
	case Floor:
		// 查询单个楼层
		g.GET("/floor", GetFloor)
		// 查询楼层内回复
		g.GET("/floor/replys", GetFloorReplys)
		// 获取用户楼层
		g.GET("/floors/user", GetUserFloors)
		// 查询多个楼层
		g.GET("/floors", GetFloors)
		// 删除指定楼层
		g.GET("/floor/delete", permission.RightDemand(models.UserRight{Super: true, StuAdmin: true}), DeleteFloor)
		// 恢复指定楼层
		g.POST("/floor/recover", permission.RightDemand(models.UserRight{Super: true}), RecoverFloor)
	case Tag:
		// 查询标签
		g.GET("/tags", GetTags)
		// 获取热议标签
		g.GET("/tags/hot", GetHotTag)
		tg := g.Group("", permission.RightDemand(models.UserRight{Super: true}))
		// 删除指定标签
		tg.GET("/tag/delete", DeleteTag)
		// 获取标签详情
		tg.GET("/tag/detail", GetTagDetail)
		// 清空标签热度
		tg.GET("/tag/clear", ClearTagPoint)
		// 增加标签热度
		tg.POST("/tag/point", AddTagPoint)
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
	case Game:
		// 获取游戏列表
		g.GET("/game", GetNewestGame)
		// 更新游戏列表
		g.POST("/game", permission.RightDemand(models.UserRight{Super: true}), AddNewGame)
	case Sensitive:
		sGroup := g.Group("", permission.RightDemand(models.UserRight{Super: true}))
		// 获取关键词文件
		sGroup.GET("/sensitive", GetSensitiveWordFile)
		// 上传关键词文件
		sGroup.POST("/sensitive", UploadSensitiveWordFile)
		// 追加词语
		sGroup.POST("/sensitive/words", AddWordsToSensitiveFile)
	case PostType:
		// 获取帖子类型
		g.GET("/posttypes", frontend.GetPostTypes)
		// 增加帖子类型
		g.POST("/posttype", permission.RightDemand(models.UserRight{Super: true}), AddPostType)
	case Banner:
		// 获取轮播图列表
		g.GET("/banners", GetBanners)
		// 更新轮播图列表
		g.POST("/banner", AddBanner)
		// 更新轮播图顺序
		g.POST("/banner/order", UpdateBannerOrder)
		// 删除轮播图
		g.GET("/banner/delete", DeleteBanner)
	case Statistic:
		// 获取帖子数量
		g.GET("/statistic/posts/count", GetPostCount)
		// 获取楼层数量
		g.GET("/statistic/floors/count", GetFloorCount)
		// 获取帖子浏览数量
		g.GET("/statistic/posts/visit/count", GetVisitPostCount)
	}
}
