package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
)

// @method [get]
// @way [query]
// @param from, to
// @return
// @route /b/statistic/posts/count
func GetPostCount(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	valid := validation.Validation{}
	valid.Required(from, "from")
	valid.Required(to, "to")
	ok, verr := r.ErrorValid(&valid, "Get posts count")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	if carbon.Parse(from).Error != nil || carbon.Parse(to).Error != nil {
		r.Error(c, e.INVALID_PARAMS, "时间格式应为YYYY-MM-dd hh:mm:ss")
		return
	}
	cnt, err := models.GetPostCount(from, to)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method [get]
// @way [query]
// @param from, to
// @return
// @route /b/statistic/floors/count
func GetFloorCount(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	valid := validation.Validation{}
	valid.Required(from, "from")
	valid.Required(to, "to")
	ok, verr := r.ErrorValid(&valid, "Get floors count")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	if carbon.Parse(from).Error != nil || carbon.Parse(to).Error != nil {
		r.Error(c, e.INVALID_PARAMS, "时间格式应为YYYY-MM-dd hh:mm:ss")
		return
	}
	cnt, err := models.GetFloorCount(from, to)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}

// @method [get]
// @way [query]
// @param from, to
// @return
// @route /b/statistic/posts/visit/count
func GetVisitPostCount(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	valid := validation.Validation{}
	valid.Required(from, "from")
	valid.Required(to, "to")
	ok, verr := r.ErrorValid(&valid, "Get posts visit count")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	if carbon.Parse(from).Error != nil || carbon.Parse(to).Error != nil {
		r.Error(c, e.INVALID_PARAMS, "时间格式应为YYYY-MM-dd hh:mm:ss")
		return
	}
	cnt, err := models.GetVisitPostCount(from, to)
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"count": cnt})
}
