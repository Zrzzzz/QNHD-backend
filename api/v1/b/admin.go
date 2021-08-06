package b

import (
	"log"
	"net/http"
	"qnhd/api/r"
	"qnhd/models"
	"qnhd/pkg/e"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetAdmins(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	maps["name"] = name
	list := models.GetAdmins(maps)
	data["list"] = list
	data["total"] = len(list)

	c.JSON(http.StatusOK, r.H(e.SUCCESS, data))
}

func AddAdmins(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	if valid.HasErrors() {
		for _, r := range valid.Errors {
			log.Printf("Add admin error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.AddAdmins(name, password)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

func EditAdmins(c *gin.Context) {
	name := c.Query("name")
	password := c.Query("password")

	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(password, "password")

	if valid.HasErrors() {
		for _, r := range valid.Errors {
			log.Printf("Add admin error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.EditAdmins(name, password)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}

func DeleteAdmins(c *gin.Context) {
	name := c.Query("name")

	valid := validation.Validation{}
	valid.Required(name, "name")

	if valid.HasErrors() {
		for _, r := range valid.Errors {
			log.Printf("Add admin error: %v", r)
		}
		c.JSON(http.StatusOK, r.H(e.INVALID_PARAMS, nil))
		return
	}

	models.DeleteAdmins(name)
	c.JSON(http.StatusOK, r.H(e.SUCCESS, nil))
}
