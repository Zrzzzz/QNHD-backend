package frontend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [post]
// @way [formdata]
// @param type, post_id, floor_id, reason
// @return
// @route /f/report
func AddReport(c *gin.Context) {
	uid := r.GetUid(c)
	rType := c.PostForm("type")
	postId := c.PostForm("post_id")
	floorId := c.PostForm("floor_id")
	reason := c.PostForm("reason")
	valid := validation.Validation{}
	valid.Required(rType, "type")
	valid.Numeric(rType, "type")
	valid.Required(postId, "post_id")
	valid.Numeric(postId, "post_id")
	if floorId != "" {
		valid.Numeric(floorId, "floor_id")
	}
	ok, verr := r.E(&valid, "Add report")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	rTypeint := util.AsInt(rType)
	valid.Range(rTypeint, 1, 2, "type")
	ok, verr = r.E(&valid, "Add report")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}
	var floorIdint uint64 = 0
	if floorId != "" {
		floorIdint = util.AsUint(floorId)
	}
	if rTypeint == 2 && floorIdint == 0 {
		r.Success(c, e.INVALID_PARAMS, nil)
		return
	}
	maps := map[string]interface{}{
		"uid":      util.AsUint(uid),
		"type":     rTypeint,
		"post_id":  util.AsUint(postId),
		"floor_id": floorIdint,
		"reason":   reason,
	}
	err := models.AddReport(maps)
	if err != nil {
		logging.Error(" error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
