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
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
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
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		logging.Warn("%v", err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
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

func CheckImage(file *multipart.File, image *multipart.FileHeader) (string, error) {
	imageName := GetImageName(image.Filename)
	fullPath := GetImageFullPath()

	src := fullPath + imageName
	var err error
	if !CheckImageExt(imageName) || !CheckImageSize(*file) {
		err = fmt.Errorf("Check image failed")
		return "", err
	} else {
		if err = CheckPath(fullPath); err != nil {
			return "", err
		}
		return src, nil
	}

}
