package b

import (
	"log"
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetBanned(c *gin.Context) {
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

	list := models.GetBanned(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

func AddBanned(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			log.Printf("Get Banned Error %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if !models.IfBannedByUid(intuid) {
		models.AddBannedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_BANNED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))

}

func DeleteBanned(c *gin.Context) {
	uid := c.Query("uid")
	valid := validation.Validation{}
	valid.Required(uid, "uid")
	valid.Numeric(uid, "uid")
	if valid.HasErrors() {
		for _, r := range valid.Errors {
			log.Printf("Get Banned Error %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}
	intuid, _ := strconv.ParseUint(uid, 10, 64)

	code := 0
	if models.IfBannedByUid(intuid) {
		models.DeleteBannedByUid(intuid)
		code = e.SUCCESS
	} else {
		code = e.ERROR_NOT_BANNED_USER
	}
	c.JSON(http.StatusOK, r.H(code, nil))
}
