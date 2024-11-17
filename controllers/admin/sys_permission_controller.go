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
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SysPermissionCreate 创建权限
func SysPermissionCreate(c *gin.Context) {
	var sysPermission models.SysPermission
	if err := c.ShouldBindJSON(&sysPermission); err != nil {
		// 检查是否是验证错误类型
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			// 遍历验证错误并构建错误信息
			errorMessages := make([]string, 0)
			for _, ve := range validationErrors {
				// 提取字段和错误类型，生成友好的错误信息
				errorMessage := fmt.Sprintf("字段 '%s' 验证失败: %s", ve.Field(), ve.ActualTag())
				errorMessages = append(errorMessages, errorMessage)
			}
			// 返回详细的验证错误信息
			utils.ErrorResponse(c, strings.Join(errorMessages, "; "), http.StatusBadRequest)
			return
		}

		// 如果不是验证错误，则打印并返回通用的错误信息
		log.Printf("请求数据绑定失败: %v", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 检查权限名称是否已经存在
	var existingSysPermission models.SysPermission
	if err := config.DB.Where("name = ? OR (route = ? AND method = ?)", sysPermission.Name, sysPermission.Route,
		sysPermission.Method).First(&existingSysPermission).Error; err == nil {
		// 如果找到相同名称的权限，返回错误响应
		log.Printf("权限名称或相同的路由和方法组合已存在 %v", existingSysPermission)
		utils.ErrorResponse(c, "权限名称或相同的路由和方法组合已存在", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&sysPermission).Error; err != nil {
		log.Printf("创建权限失败: %v", err)
		utils.ErrorResponse(c, "创建权限失败", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "权限创建成功", sysPermission)
}

// SysPermissionDelete 删除权限
func SysPermissionDelete(c *gin.Context) {
	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取分类 ID
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var sysPermission models.SysPermission

	if err := config.DB.First(&sysPermission, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "SysPermission not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&sysPermission).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete sys_permission", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "SysPermission deleted successfully", nil)
}

// SysPermissionUpdate 更新权限
func SysPermissionUpdate(c *gin.Context) {
	// 定义结构体来接收请求体中的参数
	var request struct {
		ID          int64  `json:"id" binding:"required"` // 权限 ID
		Name        string `json:"name" binding:"required"`
		Route       string `json:"route" binding:"required"`
		Method      string `json:"method" binding:"required"`
		Description string `json:"description"`
		ParentID    int64  `json:"parent_id"`
	}

	// 绑定并验证请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request data: %v", err)
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 查询权限是否存在
	var sysPermission models.SysPermission
	if err := config.DB.First(&sysPermission, request.ID).Error; err != nil {
		log.Printf("SysPermission %d not found", request.ID)
		utils.ErrorResponse(c, "权限未找到", http.StatusNotFound)
		return
	}

	// 更新权限数据
	sysPermission.Name = request.Name
	sysPermission.Description = request.Description
	sysPermission.Route = request.Route
	sysPermission.Method = request.Method
	sysPermission.ParentID = request.ParentID

	// 保存更新
	if err := config.DB.Save(&sysPermission).Error; err != nil {
		log.Printf("SysPermission %d update failed: %v", request.ID, err)
		utils.ErrorResponse(c, "权限更新失败", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, "权限更新成功", sysPermission)
}

// GetSysPermission 查询单个权限
func GetSysPermission(c *gin.Context) {
	var sysPermission models.SysPermission
	if err := config.DB.First(&sysPermission, c.Param("id")).Error; err != nil {
		log.Printf("SysPermission %s not found", c.Param("id"))
		utils.ErrorResponse(c, "未找到该系统权限", http.StatusNotFound)
		return
	}

	utils.SuccessResponse(c, "该管理员权限信息获取成功", sysPermission)
}

// GetAllSysPermissions 查询所有权限
func GetAllSysPermissions(c *gin.Context) {
	var sysPermissions []models.SysPermission
	if err := config.DB.Find(&sysPermissions).Error; err != nil {
		log.Printf("获取管理员权限列表失败 %v", err)
		utils.ErrorResponse(c, "获取管理员权限列表失败", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "管理员权限列表获取成功", sysPermissions)
}

// GetAllSysPermissionsPaginated 查询所有后台用户权限 带分页
func GetAllSysPermissionsPaginated(c *gin.Context) {

	var permissions []models.SysPermission

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
		dbQuery = dbQuery.Where("name LIKE ? OR route LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 执行查询并分页
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&permissions).Error; err != nil {
		log.Printf("获取权限列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取权限列表失败", http.StatusInternalServerError)
		return
	}

	// 获取用户总数以计算总页数
	var total int64
	config.DB.Model(&models.User{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取权限列表成功", gin.H{
		"sys_permissions": permissions,
		"total":           total,
		"currentPage":     currentPage,
		"pageSize":        pageSize,
		"totalPages":      totalPages,
	})
}

// NestedSysPermission represents a permission with potential child permissions
type NestedSysPermission struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Route       string                `json:"route"`
	Method      string                `json:"method"`
	ParentID    int64                 `json:"parent_id"`
	Children    []NestedSysPermission `json:"children"` // 子权限列表
}

// GetSysPermissionsTree 获取带父子嵌套层级的权限列表
func GetSysPermissionsTree(c *gin.Context) {
	// 获取搜索关键字
	query := c.DefaultQuery("query", "")

	// 获取完整权限树结构，包含查询条件
	permissions, err := services.GetSysPermissionTree(query)
	if err != nil {
		log.Println("Failed to retrieve permissions:", err)
		utils.ErrorResponse(c, "获取权限失败", http.StatusInternalServerError)
		return
	}

	// 返回响应
	utils.SuccessResponse(c, "获取权限成功", gin.H{
		"sys_permissions": permissions,
	})
}
