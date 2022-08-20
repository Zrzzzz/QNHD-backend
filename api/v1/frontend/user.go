package frontend

import (
	"fmt"
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /f/user
func GetUserInfo(c *gin.Context) {
	uid := r.GetUid(c)
	user, err := models.GetUser(map[string]interface{}{"id": uid})
	if err != nil {
		logging.Error("get user error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, map[string]interface{}{"user": user})
}

// @method [post]
// @way [name]
// @param
// @return
// @route /f/user/name
func EditUserName(c *gin.Context) {
	uid := r.GetUid(c)
	name := c.PostForm("name")
	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.MaxSize(name, 20, "name")
	ok, verr := r.ErrorValid(&valid, "Edit user name")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	err := checkName(name)
	if err != nil {
		logging.Error("edit user name error: %v", err)
		r.Error(c, e.INVALID_PARAMS, err.Error())
		return
	}
	err = models.EditUserName(uid, name)
	if err != nil {
		logging.Error("edit user name error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [get]
// @way [query]
// @param old, new
// @return
// @route /f/user/update_num
func UpdateUserNumber(c *gin.Context) {
	uid := r.GetUid(c)
	old := c.PostForm("old")
	new := c.PostForm("new")
	valid := validation.Validation{}
	valid.Required(old, "old")
	valid.MaxSize(old, 20, "old")
	valid.Required(new, "new")
	valid.MaxSize(new, 20, "new")
	ok, verr := r.ErrorValid(&valid, "Update user number")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	err := models.UpdateUserNumber(uid, old, new)
	if err != nil {
		logging.Error("update user number error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

func checkName(name string) error {
	if strings.Contains(name, " ") {
		return fmt.Errorf("包含空格")
	}
	ok, e := filter.NicknameFilter.Validate(name)
	if !ok {
		return fmt.Errorf("含有敏感词: %s", e)
	}
	return nil
}
