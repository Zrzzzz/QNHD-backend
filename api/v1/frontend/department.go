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
// @param name
// @return departmentList
// @route /f/departments
func GetDepartments(c *gin.Context) {
	name := c.Query("name")

	data := make(map[string]interface{})
	list, err := models.GetDepartments(name)
	if err != nil {
		logging.Error("Get department error: %v", err)
		r.Success(c, e.ERROR_DATABASE, map[string]interface{}{"error": err.Error()})
		return
	}
	data["list"] = list
	data["total"] = len(list)
	r.Success(c, e.SUCCESS, data)
}
