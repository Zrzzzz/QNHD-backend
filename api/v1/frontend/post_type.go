package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /b/posttypes
func GetPostTypes(c *gin.Context) {
	list, err := models.GetPostTypes()
	if err != nil {
		logging.Error("Get posttypes error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	// 把校务放在最后
	var retList []models.PostType
	var schType models.PostType
	for _, t := range list {
		if t.Id != models.POST_SCHOOL_TYPE {
			retList = append(retList, t)
		} else {
			schType = t
		}
	}
	retList = append(retList, schType)

	data := make(map[string]interface{})
	data["list"] = retList
	data["total"] = len(retList)

	r.OK(c, e.SUCCESS, data)
}
