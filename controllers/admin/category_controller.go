package admin

import (
	"errors"
	"exam_server/config"
	"exam_server/models"
	"exam_server/services"
	"exam_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetCategories 获取所有产品类别 无分类嵌套
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := config.DB.Find(&categories).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, "Categories fetched successfully", categories)
}

// GetCategory 获取单个产品类别
func GetCategory(c *gin.Context) {
	// 获取 URL 中的 ID 参数并转换为 uint 类型
	var category models.Category
	categoryID := c.Param("id")

	// 查询数据库中的分类信息
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		utils.ErrorResponse(c, "Category not found", http.StatusNotFound)
		return
	}
	utils.SuccessResponse(c, "Category fetched successfully", category)
}

// CreateCategory 创建新产品类别
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 如果 pid 不为 0，则检查对应父分类是否存在
	if category.Pid != 0 {
		var parentCategory models.Category
		if err := config.DB.First(&parentCategory, category.Pid).Error; err != nil {
			// 如果未找到对应的父分类，返回错误
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.ErrorResponse(c, "Parent category not found", http.StatusBadRequest)
				return
			}
			log.Printf("Failed to find parent category: %v", err)
			utils.ErrorResponse(c, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// 插入新分类
	if err := config.DB.Create(&category).Error; err != nil {
		//// 检查是否是唯一约束冲突错误
		//var mysqlErr *mysql.MySQLError
		//if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		//	utils.ErrorResponse(c, "Category name already exists. Please choose a different name.", http.StatusConflict)
		//	return
		//}
		// 检查错误信息中是否包含 "Duplicate entry" 关键字
		if strings.Contains(err.Error(), "Duplicate entry") {
			utils.ErrorResponse(c, "Category name already exists. Please choose a different name.", http.StatusConflict)
			return
		}
		log.Printf("Failed to create category: %v", err)
		utils.ErrorResponse(c, "Failed to create category", http.StatusInternalServerError)
		return
	}
	fmt.Printf("category:%v\n", category)

	utils.SuccessResponse(c, "Category created successfully", category)
}

// UpdateCategory 更新产品类别
func UpdateCategory(c *gin.Context) {
	var request struct {
		ID           int64  `json:"id" binding:"required"`   // 从请求体中获取分类 ID
		Name         string `json:"name" binding:"required"` // 其他要更新的字段
		Slug         string `json:"slug" binding:"required"` // 其他要更新的字段
		Description  string `json:"description"`             // 其他要更新的字段
		DisplayOrder int64  `json:"display_order"`
		Pid          int64  `json:"pid"` // 父分类ID
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

	var category models.Category
	// 根据 ID 查询分类
	if err := config.DB.First(&category, request.ID).Error; err != nil {
		log.Printf("Failed to fetch category: %v", err)
		utils.ErrorResponse(c, "Category not found", http.StatusNotFound)
		return
	}

	// 更新分类信息
	category.Name = request.Name
	category.Pid = request.Pid
	category.Slug = strings.ToLower(request.Slug) //别名转小写，因为它是用来前台显示到url的。为了url友好而产生的
	category.Description = request.Description
	category.DisplayOrder = request.DisplayOrder

	// 保存更新
	if err := config.DB.Save(&category).Error; err != nil {
		log.Printf("Failed to update category: %v", err)
		utils.ErrorResponse(c, "Failed to update category", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Category updated successfully", category)
}

// DeleteCategory 删除产品类别
func DeleteCategory(c *gin.Context) {
	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取分类 ID
	}
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var category models.Category

	// 根据 ID 查找分类
	if err := config.DB.First(&category, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Category not found", http.StatusNotFound)
		return
	}

	// 删除分类
	if err := config.DB.Delete(&category).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Category soft-deleted successfully", nil)
}

// GetCategoriesForCascader 返回分类的级联选择数据 级联下拉框需要的接口
func GetCategoriesForCascader(c *gin.Context) {
	data, err := services.GetCategoryCascaderData()
	if err != nil {
		utils.ErrorResponse(c, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, "Categories cascade fetched successfully", data)
}

// GetCategoriesTree 获取多级分类数据的接口 后台table显示用的
func GetCategoriesTree(c *gin.Context) {
	// 获取搜索关键字和分页参数
	query := c.DefaultQuery("query", "")
	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	fmt.Printf("query :%s", query)

	// 获取多级分类结构，包含查询条件
	categories, err := services.GetCategoriesTree(query)
	if err != nil {
		log.Println("Failed to retrieve categories:", err)
		utils.ErrorResponse(c, "error", http.StatusInternalServerError)
		return
	}

	log.Printf("categories:%v\n", categories)

	// 分页处理
	start := (currentPage - 1) * pageSize
	end := start + pageSize
	if end > len(categories) {
		end = len(categories)
	}
	paginatedCategories := categories[start:end]

	// 返回响应
	utils.SuccessResponse(c, "获取分类成功", paginatedCategories)
}
