package admin

import (
	"errors"
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetAllBrands 获取所有产品品牌
func GetAllBrands(c *gin.Context) {
	var brand []models.Brand
	if err := config.DB.Find(&brand).Error; err != nil {
		log.Println("Failed to fetch brands: ", err)
		utils.ErrorResponse(c, "Failed to fetch brands", http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, "Brand fetched successfully", brand)
}
func GetAllBrandsPaginated(c *gin.Context) {

	var brand []models.Brand

	// 获取查询参数
	query := c.DefaultQuery("query", "")
	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算偏移量
	offset := (currentPage - 1) * pageSize

	// 构建查询
	dbQuery := config.DB.Order("id DESC")

	// 如果有查询参数，按名称或邮箱模糊查询
	if query != "" {
		dbQuery = dbQuery.Where("name LIKE ? ", "%"+query+"%")
	}

	// 执行查询并分页
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&brand).Error; err != nil {
		log.Printf("获取品牌列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取品牌列表失败", http.StatusInternalServerError)
		return
	}

	// 获取总数以计算总页数
	var total int64
	config.DB.Model(&models.Brand{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取品牌列表成功", gin.H{
		"brands":      brand,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
	})
}

// GetBrand 获取单个产品品牌
func GetBrand(c *gin.Context) {
	// 获取 URL 中的 ID 参数并转换为 uint 类型
	var brand models.Brand
	brandID := c.Param("id")

	// 查询数据库中的品牌信息
	if err := config.DB.First(&brand, brandID).Error; err != nil {
		log.Printf("Failed to fetch brand: %v\n", err)
		utils.ErrorResponse(c, "Brand not found", http.StatusNotFound)
		return
	}
	utils.SuccessResponse(c, "Brand fetched successfully", brand)
}

// CreateBrand 创建新产品品牌
func CreateBrand(c *gin.Context) {
	var brand models.Brand

	if err := c.ShouldBindJSON(&brand); err != nil {
		log.Printf("Invalid request data: %v\n", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&brand).Error; err != nil {
		// 检查是否是唯一约束冲突错误
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			utils.ErrorResponse(c, "Brand name already exists. Please choose a different name.", http.StatusConflict)
			return
		}
		// 检查错误信息中是否包含 "Duplicate entry" 关键字
		//if strings.Contains(err.Error(), "Duplicate entry") {
		//	utils.ErrorResponse(c, "Brand name already exists. Please choose a different name.", http.StatusConflict)
		//	return
		//}
		log.Printf("Failed to create brand: %v", err)
		utils.ErrorResponse(c, "Failed to create brand", http.StatusInternalServerError)
		return
	}
	fmt.Printf("brand:%v\n", brand)

	utils.SuccessResponse(c, "Brand created successfully", brand)
}

// UpdateBrand 更新产品品牌
func UpdateBrand(c *gin.Context) {
	var request struct {
		ID          int64  `json:"id" binding:"required"`   // 从请求体中获取品牌 ID
		Name        string `json:"name" binding:"required"` // 其他要更新的字段
		Description string `json:"description"`
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

	var brand models.Brand
	// 根据 ID 查询品牌
	if err := config.DB.First(&brand, "id = ?", request.ID).Error; err != nil {
		log.Printf("Failed to fetch brand: %v", err)
		utils.ErrorResponse(c, "Brand not found", http.StatusNotFound)
		return
	}

	// 更新品牌信息
	brand.Name = request.Name
	brand.Description = request.Description

	// 保存更新
	if err := config.DB.Save(&brand).Error; err != nil {
		log.Printf("Failed to update brand: %v", err)
		utils.ErrorResponse(c, "Failed to update brand", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Brand updated successfully", brand)
}

// DeleteBrand 删除产品品牌
func DeleteBrand(c *gin.Context) {
	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取品牌 ID
	}
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var brand models.Brand

	// 根据 ID 查找品牌
	if err := config.DB.First(&brand, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Brand not found", http.StatusNotFound)
		return
	}

	// 删除品牌
	if err := config.DB.Delete(&brand).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete brand", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Brand soft-deleted successfully", nil)
}

// DeleteBrandsBatch 批量删除
func DeleteBrandsBatch(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"` // 从请求体中获取 ID 列表
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request data: %v", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 执行软删除
	if err := config.DB.Where("id IN ?", request.IDs).Delete(&models.Brand{}).Error; err != nil {
		log.Printf("Failed to delete brands: %v", err)
		utils.ErrorResponse(c, "Failed to delete brands", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Brands soft-deleted successfully", nil)
}
