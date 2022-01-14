package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"qnhd/pkg/file"
	"qnhd/pkg/logging"
	"qnhd/pkg/setting"
	"qnhd/pkg/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName + fmt.Sprintf("%d", time.Now().Unix()))
	return fileName + ext
}
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}
func GetRuntimePath() string {
	return setting.AppSetting.RuntimeRootPath
}
func GetImageFullPath() string {
	return GetRuntimePath() + GetImagePath()
}
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}
func CheckImageSize(f *multipart.FileHeader) bool {
	size := f.Size
	return int(size) <= setting.AppSetting.ImageMaxSize
}
func CheckPath(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}
	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
func CheckImage(image *multipart.FileHeader) error {
	fullPath := GetImageFullPath()
	imageName := GetImageName(image.Filename)
	var err error
	if !CheckImageExt(imageName) || !CheckImageSize(image) {
		err = fmt.Errorf("Check image failed")
		return err
	} else {
		if err = CheckPath(fullPath); err != nil {
			return err
		}
		return nil
	}
}
func GetImageSrc(image *multipart.FileHeader) string {
	imageName := GetImageName(image.Filename)
	fullPath := GetImageFullPath()
	src := fullPath + imageName
	return src
}
func SaveImagesFromFromData(imgs []*multipart.FileHeader, c *gin.Context) ([]string, error) {
	var imageUrls = []string{}
	var err error
	// 检查每张图
	for _, img := range imgs {
		err = CheckImage(img)
		if err != nil {
			logging.Error("Add post error: %v", err)
			return imageUrls, err
		}
	}
	// 对每个图片进行处理
	for _, img := range imgs {
		src := GetImageSrc(img)
		if err := c.SaveUploadedFile(img, src); err != nil {
			logging.Error("Add post error: %v", err)
			return imageUrls, err
		}
		imageName := GetImageName(img.Filename)
		imageUrls = append(imageUrls, GetImagePath()+imageName)
	}
	return imageUrls, nil
}
func DeleteImageUrls(urls []string) error {
	for _, url := range urls {
		if url == "" {
			continue
		}
		err := os.Remove(GetRuntimePath() + url)
		if err != nil {
			logging.Error(err.Error())
			return err
		}
	}
	return nil
}
