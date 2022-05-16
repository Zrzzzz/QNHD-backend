package common

import (
	"qnhd/crypto"
	"qnhd/enums/PostSearchModeType"
	"qnhd/enums/PostValueModeType"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetPosts(front bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		postType := c.Query("type")
		searchMode := c.Query("search_mode")
		content := c.Query("content")
		departmentId := c.Query("department_id")
		solved := c.Query("solved")
		tagId := c.Query("tag_id")
		valueMode := c.Query("value_mode")
		isDeleted := c.Query("is_deleted")
		valid := validation.Validation{}
		valid.Required(postType, "type")
		valid.Numeric(postType, "type")
		valid.Required(searchMode, "search_mode")
		valid.Numeric(searchMode, "search_mode")
		if valueMode == "" {
			// 默认值
			valueMode = "0"
		}
		valid.Numeric(valueMode, "value_mode")
		ok, verr := r.ErrorValid(&valid, "Get posts")
		if !ok {
			r.Error(c, e.INVALID_PARAMS, verr.Error())
			return
		}
		valid.Numeric(solved, "solved")
		valid.Numeric(departmentId, "department_id")
		valid.Numeric(tagId, "tag_id")
		postTypeint := util.AsInt(postType)
		searchModeint := util.AsInt(searchMode)
		valid.Range(searchModeint, 0, 1, "search_mode")
		if solved != "" {
			solvedint := util.AsInt(solved)
			valid.Range(solvedint, 0, 2, "solved")
		}
		ok, verr = r.ErrorValid(&valid, "Get posts")
		if !ok {
			r.Error(c, e.INVALID_PARAMS, verr.Error())
			return
		}

		data := make(map[string]interface{})
		maps := map[string]interface{}{
			"type":          postTypeint,
			"search_mode":   PostSearchModeType.Enum(searchModeint),
			"content":       content,
			"solved":        solved,
			"department_id": departmentId,
			"tag_id":        tagId,
			"value_mode":    PostValueModeType.Enum(util.AsInt(valueMode)),
			"is_deleted":    isDeleted,
		}
		if front {
			uid := r.GetUid(c)
			list, err := models.GetPostResponsesWithUid(c, uid, maps)
			if err != nil {
				logging.Error("Get posts error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			for i := range list {
				list[i].Uid = crypto.Encrypt(list[i].Uid, list[i].Id)
			}
			data["list"] = list
			data["total"] = len(list)
		} else {
			list, cnt, err := models.GetPostResponses(c, maps)
			if err != nil {
				logging.Error("Get posts error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			data["list"] = list
			data["total"] = cnt
		}

		r.OK(c, e.SUCCESS, data)
	}
}

func GetPost(front bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		valid := validation.Validation{}
		valid.Required(id, "id")
		valid.Numeric(id, "id")

		ok, verr := r.ErrorValid(&valid, "Get Posts")
		if !ok {
			r.Error(c, e.INVALID_PARAMS, verr.Error())
			return
		}
		data := make(map[string]interface{})
		if front {
			uid := r.GetUid(c)
			pr, err := models.GetPostResponseUserAndVisit(id, uid)
			if err != nil {
				logging.Error("Get post error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			pr.Uid = crypto.Encrypt(pr.Uid, pr.Id)
			data["post"] = pr
		} else {
			pr, err := models.GetPostResponse(id)
			if err != nil {
				logging.Error("Get post error: %v", err)
				r.Error(c, e.ERROR_DATABASE, err.Error())
				return
			}
			data["post"] = pr
		}
		r.OK(c, e.SUCCESS, data)
	}
}
