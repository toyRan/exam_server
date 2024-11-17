package admin

import (
	"errors"
	"exam_server/config"
	"exam_server/models"
	"exam_server/services"
	"exam_server/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

// GetAllSeries 获取所有产品系列
func GetAllSeries(c *gin.Context) {
	var series []models.Series
	if err := config.DB.Find(&series).Error; err != nil {
		log.Println("Failed to fetch seriess: ", err)
		utils.ErrorResponse(c, "Failed to fetch seriess", http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, "Series fetched successfully", series)
}

func GetAllSeriesPaginated(c *gin.Context) {
	var seriesRecords []models.Series

	// 获取查询参数
	query := c.DefaultQuery("query", "")
	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算偏移量
	offset := (currentPage - 1) * pageSize

	// 构建查询
	dbQuery := config.DB.Preload("FrameMaterial").Order("id DESC")

	// 如果有查询参数，按名称或邮箱模糊查询
	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ? ", "%"+query+"%")
	}

	// 执行查询并分页
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&seriesRecords).Error; err != nil {
		log.Printf("获取系列号列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取系列号列表失败", http.StatusInternalServerError)
		return
	}

	// 创建扁平化的响应结构
	type FlattenedSeries struct {
		ID                int64  `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		PdfURL            string `json:"pdf_url"`
		IsNewDesign       bool   `json:"is_new_design"`
		FrameMaterialID   int64  `json:"frame_material_id"`
		FrameMaterialName string `json:"frame_material_name"`
		// ... 其他需要的字段
	}

	// 转换数据为扁平化结构
	flattenedSeries := make([]FlattenedSeries, len(seriesRecords))
	for i, series := range seriesRecords {
		flattenedSeries[i] = FlattenedSeries{
			ID:                series.ID,
			Name:              series.Name,
			Description:       series.Description,
			PdfURL:            series.PdfURL,
			IsNewDesign:       series.IsNewDesign,
			FrameMaterialID:   series.FrameMaterialID,
			FrameMaterialName: series.FrameMaterial.Name,
		}
	}

	// 获取总数并计算总页数
	var total int64
	config.DB.Model(&models.Series{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	// 返回扁平化的响应
	utils.SuccessResponse(c, "获取系列列表成功", gin.H{
		"series":      flattenedSeries,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
	})
}

// GetSeries 获取单个产品系列
func GetSeries(c *gin.Context) {
	var series models.Series
	seriesID := c.Param("id")

	// 查询数据库中的系列信息
	if err := config.DB.Preload("FrameMaterial").First(&series, seriesID).Error; err != nil {
		log.Printf("Failed to fetch series: %v\n", err)
		utils.ErrorResponse(c, "Series not found", http.StatusNotFound)
		return
	}

	// 创建扁平化的响应结构
	flattenedSeries := struct {
		ID                int64  `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		PdfURL            string `json:"pdf_url"`
		FrameMaterialID   int64  `json:"frame_material_id"`
		FrameMaterialName string `json:"frame_material_name"`
		// ... 其他需要的字段
	}{
		ID:                series.ID,
		Name:              series.Name,
		Description:       series.Description,
		PdfURL:            series.PdfURL,
		FrameMaterialID:   series.FrameMaterialID,
		FrameMaterialName: series.FrameMaterial.Name,
	}

	utils.SuccessResponse(c, "Series fetched successfully", flattenedSeries)
}

// CreateSeries 创建新产品系列
func CreateSeries(c *gin.Context) {
	// 定义请求结构体，添加验证规则
	var request struct {
		Name            string `json:"name" binding:"required"`
		Description     string `json:"description"`
		PdfURL          string `json:"pdf_url" binding:"required"`
		FrameMaterialID int64  `json:"frame_material_id" binding:"required"`
		IsNewDesign     bool   `json:"is_new_design"` // 新增字段：是否为最新设计
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			// 格式化验证错误信息
			var errMsgs []string
			for _, e := range validationErrors {
				switch e.Field() {
				case "FrameMaterialID":
					errMsgs = append(errMsgs, "Frame material ID is required")
				case "Name":
					errMsgs = append(errMsgs, "Name is required")
				case "PdfURL":
					errMsgs = append(errMsgs, "PDF URL is required")
				}
			}
			utils.ErrorResponse(c, strings.Join(errMsgs, "; "), http.StatusBadRequest)
		} else {
			// 处理其他类型的错误（如 JSON 解析错误）
			utils.ErrorResponse(c, "Invalid request format", http.StatusBadRequest)
		}
		return
	}

	// 验证frame_material是否存在
	var frameMaterial models.FrameMaterial
	if err := config.DB.First(&frameMaterial, request.FrameMaterialID).Error; err != nil {
		utils.ErrorResponse(c, "Invalid frame material ID", http.StatusBadRequest)
		return
	}

	// 创建系列记录
	series := models.Series{
		Name:            request.Name,
		Description:     request.Description,
		PdfURL:          request.PdfURL,
		FrameMaterialID: request.FrameMaterialID,
		IsNewDesign:     request.IsNewDesign, // 设置新设计标志
	}

	if err := config.DB.Create(&series).Error; err != nil {
		// 检查是否是唯一约束冲突错误
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			utils.ErrorResponse(c, "Series name already exists. Please choose a different name.", http.StatusConflict)
			return
		}
		log.Printf("Failed to create series: %v", err)
		utils.ErrorResponse(c, "Failed to create series", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Series created successfully", series)
}

// UpdateSeries 更新产品系列
func UpdateSeries(c *gin.Context) {
	var request struct {
		ID              int64  `json:"id" binding:"required"`
		Name            string `json:"name" binding:"required"`
		PdfURL          string `json:"pdf_url"  binding:"required"`
		Description     string `json:"description"`
		FrameMaterialID int64  `json:"frame_material_id" binding:"required"`
		IsNewDesign     bool   `json:"is_new_design"` // 新增字段：是否为最新设计
	}
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		// 将错误信息格式化为用户友好的消息
		var errs validator.ValidationErrors
		errors.As(err, &errs)
		var errMsg []string
		for _, e := range errs {
			// 获取每个字段的错误信息并格式化
			errMsg = append(errMsg, fmt.Sprintf("Field '%s' is %s", e.Field(), e.Tag()))
		}

		log.Printf("Invalid request data: %v", errMsg)
		utils.ErrorResponse(c, strings.Join(errMsg, ", "), http.StatusBadRequest)
		return
	}

	var series models.Series
	// 根据 ID 查询系列
	if err := config.DB.First(&series, "id = ?", request.ID).Error; err != nil {
		log.Printf("Failed to fetch series: %v", err)
		utils.ErrorResponse(c, "Series not found", http.StatusNotFound)
		return
	}

	// 验证frame_material是否存在
	var frameMaterial models.FrameMaterial
	if err := config.DB.First(&frameMaterial, request.FrameMaterialID).Error; err != nil {
		utils.ErrorResponse(c, "Invalid frame material ID", http.StatusBadRequest)
		return
	}

	// 如果PDF URL发生变化，删除旧文件
	if series.PdfURL != "" && series.PdfURL != request.PdfURL {
		if err := services.DeleteFilesFromOSS([]string{series.PdfURL}); err != nil {
			log.Printf("Failed to delete old PDF file from OSS: %v", err)
		}
	}

	// 更新系列信息
	series.Name = request.Name
	series.Description = request.Description
	series.PdfURL = request.PdfURL
	series.FrameMaterialID = request.FrameMaterialID
	series.IsNewDesign = request.IsNewDesign

	fmt.Printf("series.PdfURL :%v\n", request.PdfURL)

	// 保存更新
	if err := config.DB.Save(&series).Error; err != nil {
		log.Printf("Failed to update series: %v", err)
		utils.ErrorResponse(c, "Failed to update series", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Series updated successfully", series)
}

// DeleteSeries 删除产品系列
func DeleteSeries(c *gin.Context) {
	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取系列 ID
	}
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("DeleteSeries Invalid request data: %v", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var series models.Series

	// 根据 ID 查找系列
	if err := config.DB.First(&series, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Series not found", http.StatusNotFound)
		return
	}

	// 删除OSS中的PDF文件
	if series.PdfURL != "" {
		if err := services.DeleteFilesFromOSS([]string{series.PdfURL}); err != nil {
			log.Printf("Failed to delete PDF file from OSS: %v", err)
		}
	}

	// 删除系列
	if err := config.DB.Delete(&series).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete series", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Series soft-deleted successfully", nil)
}

// DeleteSeriesBatch 批量删除产品系列
func DeleteSeriesBatch(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"` // 从请求体中获取系列 ID 列表
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request data: %v", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 先获取要删除的记录
	var seriesToDelete []models.Series
	if err := config.DB.Where("id IN ?", request.IDs).Find(&seriesToDelete).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch series for deletion", http.StatusInternalServerError)
		return
	}

	// 收集所有需要删除的PDF URLs
	var pdfURLs []string
	for _, series := range seriesToDelete {
		if series.PdfURL != "" {
			pdfURLs = append(pdfURLs, series.PdfURL)
		}
	}

	// 批量删除OSS中的PDF文件
	if len(pdfURLs) > 0 {
		if err := services.DeleteFilesFromOSS(pdfURLs); err != nil {
			log.Printf("Failed to delete PDF files from OSS: %v", err)
		}
	}

	// 执行软删除
	if err := config.DB.Where("id IN ?", request.IDs).Delete(&models.Series{}).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete series", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Series deleted successfully", nil)
}
