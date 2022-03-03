package backend

import (
	"bufio"
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

func AddWordsToSensitiveFile(c *gin.Context) {
	words := c.PostFormArray("words")

	filePath := "conf/sensitive.txt"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logging.Error("open file error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	//及时关闭file句柄
	defer file.Close()
	//读原来文件的内容，并且显示在终端
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for _, w := range words {
		write.WriteString(w + "\n")
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	if err := filter.Reload(); err != nil {
		logging.Error("reload filter error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

func UploadSensitiveWordFile(c *gin.Context) {
	word, err := c.FormFile("word")
	if err != nil {
		logging.Error("upload sensitive word error: %v", err)
		r.Error(c, e.INVALID_PARAMS, err.Error())
		return
	}
	if err := c.SaveUploadedFile(word, "conf/sensitive.txt"); err != nil {
		logging.Error("save file error: %v", err)
		r.Error(c, e.ERROR_SAVE_FILE, err.Error())
		return
	}
	if err := filter.Reload(); err != nil {
		logging.Error("reload filter error: %v", err)
		r.Error(c, e.ERROR_SERVER, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
