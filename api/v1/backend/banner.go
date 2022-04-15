package backend

import (
	"qnhd/models"
	"qnhd/pkg/e"
	"qnhd/pkg/r"
	"qnhd/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// @method [get]
// @way [query]
// @param
// @return
// @route /b/banners
func GetBanners(c *gin.Context) {
	list, err := models.GetBanners()
	if err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	data := make(map[string]interface{})
	data["list"] = list
	data["total"] = len(list)
	r.OK(c, e.SUCCESS, data)
}

// @method [post]
// @way [formdata]
// @param content
// @return
// @route /b/banner
func AddBanner(c *gin.Context) {
	name := c.PostForm("name")
	title := c.PostForm("title")
	image := c.PostForm("image")
	url := c.PostForm("url")
	valid := validation.Validation{}
	valid.Required(name, "name")
	valid.Required(title, "title")
	valid.Required(image, "image")
	valid.Required(url, "url")
	ok, verr := r.ErrorValid(&valid, "Add banner error")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}
	maps := map[string]interface{}{
		"name":  name,
		"title": title,
		"image": image,
		"url":   url,
	}

	if err := models.AddBanner(maps); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [post]
// @way [formdata]
// @param banner_id, order
// @return
// @route /b/banner/order
func UpdateBannerOrder(c *gin.Context) {
	bannerId := c.PostForm("banner_id")
	order := c.PostForm("order")
	valid := validation.Validation{}
	valid.Required(bannerId, "banner_id")
	valid.Numeric(bannerId, "banner_id")
	valid.Required(order, "order")
	valid.Numeric(order, "order")
	ok, verr := r.ErrorValid(&valid, "Update banner order error")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	if err := models.UpdateBannerOrder(util.AsUint(bannerId), util.AsInt(order)); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}

// @method [delete]
// @way [query]
// @param banner_id
// @return
// @route /b/banner/delete
func DeleteBanner(c *gin.Context) {
	bannerId := c.Query("banner_id")
	valid := validation.Validation{}
	valid.Required(bannerId, "banner_id")
	valid.Numeric(bannerId, "banner_id")
	ok, verr := r.ErrorValid(&valid, "Delete banner error")
	if !ok {
		r.Error(c, e.INVALID_PARAMS, verr.Error())
		return
	}

	if err := models.DeleteBanner(util.AsUint(bannerId)); err != nil {
		r.Error(c, e.ERROR_DATABASE, err.Error())
		return
	}
	r.OK(c, e.SUCCESS, nil)
}
