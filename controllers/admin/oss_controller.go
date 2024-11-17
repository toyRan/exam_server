package admin

import (
	"exam_server/services"
	"exam_server/utils"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 添加允许的文件类型常量
var (
	allowedImageExts = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	allowedDocExts = map[string]bool{
		".pdf": true, ".xlsx": true, ".xls": true,
	}
)

// UploadFiles 支持多文件上传
func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Print("上传文件时解析表单错误: ", err)
		utils.ErrorResponse(c, "Failed to parse form", http.StatusBadRequest)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.ErrorResponse(c, "No files provided", http.StatusBadRequest)
		return
	}

	folderPath := c.DefaultPostForm("folder", "files")
	var responses []string
	var failedFiles []string

	// 检查所有文件类型
	for _, file := range files {
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		if !isAllowedFileType(fileExt) {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		// 如果未指定文件夹，根据文件类型自动设置存储目录
		uploadFolder := folderPath
		if folderPath == "files" {
			if _, isImage := allowedImageExts[fileExt]; isImage {
				uploadFolder = "images"
			} else if _, isDoc := allowedDocExts[fileExt]; isDoc {
				uploadFolder = "pdfs"
			}
		}

		fileURLs, err := services.UploadFilesToOSS([]*multipart.FileHeader{file}, uploadFolder)
		if err != nil {
			log.Printf("上传文件 %s 到OSS失败: %v", file.Filename, err)
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		responses = append(responses, fileURLs[0])
	}

	// 返回上传结果
	result := gin.H{
		"success_files": responses,
	}
	if len(failedFiles) > 0 {
		result["failed_files"] = failedFiles
	}

	if len(responses) > 0 {
		utils.SuccessResponse(c, "文件上传完成", result)
	} else {
		utils.ErrorResponse(c, "所有文件上传失败", http.StatusBadRequest)
	}
}

// 辅助函数：检查文件类型是否允许
func isAllowedFileType(ext string) bool {
	_, isImage := allowedImageExts[ext]
	_, isDoc := allowedDocExts[ext]
	return isImage || isDoc
}
