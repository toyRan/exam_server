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

// GetAllFrameMaterialsPaginated 查询所有管理员角色带分页
func GetAllFrameMaterialsPaginated(c *gin.Context) {

	var frameMaterial []models.FrameMaterial

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
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&frameMaterial).Error; err != nil {
		log.Printf("获取框材质列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取框材质列表失败", http.StatusInternalServerError)
		return
	}

	// 获取框材质总数以计算总页数
	var total int64
	config.DB.Model(&models.FrameMaterial{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取列表成功", gin.H{
		"frame_materials": frameMaterial,
		"total":           total,
		"currentPage":     currentPage,
		"pageSize":        pageSize,
		"totalPages":      totalPages,
	})
}

// GetAllFrameMaterials 获取所有产品框材质 不分页
func GetAllFrameMaterials(c *gin.Context) {
	var frameMaterials []models.FrameMaterial
	if err := config.DB.Find(&frameMaterials).Error; err != nil {
		log.Println("Failed to fetch frame_materials: ", err)
		utils.ErrorResponse(c, "Failed to fetch frame_materials", http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(c, "FrameMaterials fetched successfully", frameMaterials)
}

// GetFrameMaterial 获取单个产品框材质
func GetFrameMaterial(c *gin.Context) {
	// 获取 URL 中的 ID 参数并转换为 uint 类型
	var frameMaterial models.FrameMaterial
	frameMaterialID := c.Param("id")

	// 查询数据库中的框材质信息
	if err := config.DB.First(&frameMaterial, frameMaterialID).Error; err != nil {
		log.Printf("Failed to fetch frame_material: %v\n", err)
		utils.ErrorResponse(c, "FrameMaterial not found", http.StatusNotFound)
		return
	}
	utils.SuccessResponse(c, "FrameMaterial fetched successfully", frameMaterial)
}

// CreateFrameMaterial 创建新产品框材质
func CreateFrameMaterial(c *gin.Context) {
	var frameMaterial models.FrameMaterial

	if err := c.ShouldBindJSON(&frameMaterial); err != nil {
		log.Printf("Invalid request data: %v\n", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&frameMaterial).Error; err != nil {
		// 检查是否是唯一约束冲突错误
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			utils.ErrorResponse(c, "FrameMaterial name already exists. Please choose a different name.", http.StatusConflict)
			return
		}
		// 检查错误信息中是否包含 "Duplicate entry" 关键字
		//if strings.Contains(err.Error(), "Duplicate entry") {
		//	utils.ErrorResponse(c, "FrameMaterial name already exists. Please choose a different name.", http.StatusConflict)
		//	return
		//}
		log.Printf("Failed to create frame_material: %v", err)
		utils.ErrorResponse(c, "Failed to create frame_material", http.StatusInternalServerError)
		return
	}
	fmt.Printf("frame_material:%v\n", frameMaterial)

	utils.SuccessResponse(c, "FrameMaterial created successfully", frameMaterial)
}

// UpdateFrameMaterial 更新产品框材质
func UpdateFrameMaterial(c *gin.Context) {
	var request struct {
		ID          int64  `json:"id" binding:"required"`   // 从请求体中获取框材质 ID
		Name        string `json:"name" binding:"required"` // 其他要更新的字段
		Description string `json:"description"`
	}
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		// 将错误信息格式化为框材质友好的消息
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

	var frameMaterial models.FrameMaterial
	// 根据 ID 查询框材质
	if err := config.DB.First(&frameMaterial, "id = ?", request.ID).Error; err != nil {
		log.Printf("Failed to fetch frame_material: %v", err)
		utils.ErrorResponse(c, "FrameMaterial not found", http.StatusNotFound)
		return
	}

	// 更新框材质信息
	frameMaterial.Name = request.Name
	frameMaterial.Description = request.Description

	// 保存更新
	if err := config.DB.Save(&frameMaterial).Error; err != nil {
		log.Printf("Failed to update frame_material: %v", err)
		utils.ErrorResponse(c, "Failed to update frame_material", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "FrameMaterial updated successfully", frameMaterial)
}

// DeleteFrameMaterial 删除产品框材质
func DeleteFrameMaterial(c *gin.Context) {
	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取框材质 ID
	}
	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var frameMaterial models.FrameMaterial

	// 根据 ID 查找框材质
	if err := config.DB.First(&frameMaterial, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "FrameMaterial not found", http.StatusNotFound)
		return
	}

	// 删除框材质
	if err := config.DB.Delete(&frameMaterial).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete frame_material", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "FrameMaterial soft-deleted successfully", nil)
}
