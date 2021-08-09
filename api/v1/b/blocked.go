package b

import (
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/logging"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetBlocked(c *gin.Context) {
	uid := c.Query("uid")

	valid := validation.Validation{}
	valid.Numeric(uid, "uid")
	if valid.HasErrors() {
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if uid != "" {
		maps["uid"] = uid
	}

	list := models.GetBlocked(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

func AddBlocked(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("Add blocked error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if !models.IfBlockedByUid(intuid) {
		models.AddBlockedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_BLOCKED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))

}

func DeleteBlocked(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			logging.Error("Delete blocked error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if models.IfBlockedByUid(intuid) {
		models.DeleteBlockedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_BLOCKED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))
}
