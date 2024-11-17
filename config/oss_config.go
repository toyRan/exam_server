package config

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
)

var OssClient *oss.Client

// InitOSS 初始化 OSS 客户端
func InitOSS() {
	endpoint := os.Getenv("OSS_ENDPOINT")
	accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")

	// 设置超时时间
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret,
		oss.Timeout(60*5, 120*5)) // 第一个参数是连接超时，第二个是读写超时
	if err != nil {
		log.Fatalf("Error initializing OSS: %v", err)
	}

	OssClient = client
}
