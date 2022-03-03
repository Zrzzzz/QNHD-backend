package backend

import (
	"os"
	"qnhd/pkg/e"
	"qnhd/pkg/filter"
	"qnhd/pkg/logging"
	"qnhd/pkg/r"

	"github.com/gin-gonic/gin"
)

func GetSensitiveWordFile(c *gin.Context) {
	c.Header("content-disposition", "attachment; filename=word.txt")
	file := "conf/sensitive.txt"
	data, err := os.ReadFile(file)
	if err != nil {
		logging.Error("get file error: %v", err)
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	c.Data(200, "txt", data)
}

func UploadSensitiveWordFile(c *gin.Context) {
	word, err := c.FormFile("word")
	if err != nil {
		logging.Error("upload sensitive word error: %v", err)
		r.Error(c, e.INVALID_PARAMS, err.Error())
		return
	}
	if err := c.SaveUploadedFile(word, "conf/sensitive.txt"); err != nil {
		logging.Error(" error: %v", err)
		r.Error(c, e.ERROR_SAVE_FILE, err.Error())
		return
	}
	if err := filter.Reload(); err != nil {
		logging.Error(" error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
