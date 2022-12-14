package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// require content have page and page_size param
// return overnum, neednum
func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		enable := c.Query("page_disable")
		if enable == "1" {
			return db
		}
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		base, _ := strconv.Atoi(c.Query("page_base"))

		offset := (page - 1) * pageSize
		return db.Offset(base + offset).Limit(pageSize)
	}
}
