package oss

import (
	"go_webapp/global"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(999999-100000)+100000)

	return fileName + ext
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func CheckContainExt(name string) bool {
	ext := GetFileExt(name)
	ext = strings.ToUpper(ext)
	for _, allowExt := range global.AppSetting.UploadImageAllowExts {
		if strings.ToUpper(allowExt) == ext {
			return true
		}
	}

	return false
}

func CheckMaxSize(f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
		return true
	}

	return false
}
