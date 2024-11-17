package services

import (
	"exam_server/config"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// getBucket 获取OSS bucket的公共方法
func getBucket() (*oss.Bucket, error) {
	bucket, err := config.OssClient.Bucket(os.Getenv("OSS_BUCKET_NAME"))
	if err != nil {
		log.Printf("Failed to get bucket: %v", err)
		return nil, fmt.Errorf("failed to get OSS bucket: %v", err)
	}
	return bucket, nil
}

// UploadFilesToOSS 上传多个文件到阿里云OSS
func UploadFilesToOSS(files []*multipart.FileHeader, folderName string) ([]string, error) {
	bucket, err := getBucket()
	if err != nil {
		return nil, err
	}

	var fileURLs []string
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}

		defer func() {
			if err := src.Close(); err != nil {
				log.Printf("关闭文件错误: %v", err)
			}
		}()

		fileName := fmt.Sprintf("%s/%d_%s", folderName, time.Now().UnixNano(), file.Filename)
		if err = bucket.PutObject(fileName, src); err != nil {
			return nil, fmt.Errorf("failed to upload file to OSS: %v", err)
		}

		fileURL := fmt.Sprintf("https://%s.%s/%s",
			os.Getenv("OSS_BUCKET_NAME"),
			os.Getenv("OSS_ENDPOINT"),
			fileName)
		fileURLs = append(fileURLs, fileURL)
	}

	return fileURLs, nil
}

// DeleteFilesFromOSS 批量删除OSS文件
func DeleteFilesFromOSS(urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	bucket, err := getBucket()
	if err != nil {
		return err
	}

	objectKeys := make([]string, 0, len(urls))
	for _, fileURL := range urls {
		if strings.Contains(fileURL, "http") {
			parsedURL, err := url.Parse(fileURL)
			if err != nil {
				log.Printf("Invalid URL format: %s", fileURL)
				continue
			}
			objectKeys = append(objectKeys, strings.TrimPrefix(parsedURL.Path, "/"))
		} else {
			objectKeys = append(objectKeys, fileURL)
		}
	}

	if len(objectKeys) == 0 {
		return nil
	}

	// 分批处理，阿里云 OSS 每次最多删除 1000 个文件
	const batchSize = 1000
	for i := 0; i < len(objectKeys); i += batchSize {
		end := i + batchSize
		if end > len(objectKeys) {
			end = len(objectKeys)
		}

		batch := objectKeys[i:end]
		if _, err := bucket.DeleteObjects(batch, oss.DeleteObjectsQuiet(true)); err != nil {
			log.Printf("Failed to delete batch files from OSS: %v", err)
			return err
		}
	}

	return nil
}
