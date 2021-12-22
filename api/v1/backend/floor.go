package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param post_id, page=0, page_size
// @return floors
func GetFloors(c *gin.Context) {
	postId := c.Query("post_id")

	valid := validation.Validation{}
	valid.Required(postId, "postId")
	valid.Numeric(postId, "postId")
	ok, verr := r.E(&valid, "Get floors")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	list, err := models.GetFloorResponsesInPost(c, postId, "0")
	if err != nil {
		logging.Error("Get floors error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.Success(c, e.SUCCESS, data)
}

// @method [delete]
// @way [query]
// @param floor_id
// @return nil
func DeleteFloor(c *gin.Context) {
	floorId := c.Query("floor_id")

	valid := validation.Validation{}
	valid.Required(floorId, "floorId")
	valid.Numeric(floorId, "floorId")
	ok, verr := r.E(&valid, "Get floors")
	if !ok {
		r.Success(c, e.INVALID_PARAMS, map[string]interface{}{"error": verr.Error()})
		return
	}

	_, err := models.DeleteFloorByAdmin(floorId)
	if err != nil {
		logging.Error("Delete floor error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	r.Success(c, e.SUCCESS, nil)
}
