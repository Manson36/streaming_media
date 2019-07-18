package ossops

import (
	"github.com/streaming_media/streamserver/config"
	"log"
)

var EP string
var AK string  //access key 在阿里云个人信息可以查看
var SK string  //同AK

func init() {
	EP = config.GetOssAddr()
	AK = ""
	SK = ""
}

func UploadToOss(filename string, path string, bn string) bool {//bn: bucketName
	client, err := oss.New(EP, AK, SK)
	if err != nil {
		log.Printf("Init oss sevice error: %s", err)
		return false
	}

	bucket, err := client.Bucket(bn)
	if err != nil {
		log.Printf("Getting bucket error: %s", err)
		return false
	}

	//在这里我们为了增加上传速度，使用批量上传
	err = bucket.UploadFile(filename, path, 500*1024, oss.routines(3))
	if err != nil {
		log.Printf("Uploading object error: %s", err)
		return false
	}

	return true
}

func DeleteObject(filename string, bn string) bool {
	client, err := oss.New(EP, AK, SK)
	if err != nil {
		log.Printf("Init oss service error: %s", err)
		return false
	}

	bucket, err := client.Bucket(bn)
	if err != nil {
		log.Printf("Getting bucket error: %s", err)
		return false
	}

	err = bucket.DeleteObject(filename)
	if err != nil {
		log.Printf("Deleting object error: %s", err)
		return false
	}

	return true
}