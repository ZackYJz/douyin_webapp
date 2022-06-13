package oss

import (
	"fmt"
	"go_webapp/global"
	"mime/multipart"
	"os"

	"go.uber.org/zap"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func UploadOss(file *multipart.FileHeader, ossFilename string) string {
	Endpoint := "https://oss-cn-chengdu.aliyuncs.com" //
	AccessKeyID := "yours"                            //
	AccessKeySecret := "yours"                        //
	BucketName := "zackyj-typora"

	client, err := oss.New(Endpoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		global.Logger.Error("Create Oss Client Error", zap.Error(err))
	}

	// 指定bucket
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		global.Logger.Error("Get Oss bucket Error", zap.Error(err))
	}
	src, err := file.Open()
	if err != nil {
		global.Logger.Error("Open Uploaded File Error", zap.Error(err))
	}
	defer src.Close()

	// 将文件流上传至test目录下
	path := "douyin/" + ossFilename
	err = bucket.PutObject(path, src)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	return fmt.Sprintf("https://%s.%s/%s", BucketName, Endpoint[8:], path)
}
