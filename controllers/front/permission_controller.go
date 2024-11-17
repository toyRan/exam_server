package front

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// PermissionCreate 创建角色
func PermissionCreate(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 检查权限名称是否已经存在
	var existingPermission models.Permission
	if err := config.DB.Where("name = ?", permission.Name).First(&existingPermission).Error; err == nil {
		// 如果找到相同名称的权限，返回错误响应
		log.Printf("Permission %s already exists", permission.Name)
		utils.ErrorResponse(c, "Permission with the same name already exists", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&permission).Error; err != nil {
		log.Printf("Error creating permission: %v", err)
		utils.ErrorResponse(c, "Failed to create permission", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Permission created successfully", permission)
}

// PermissionDelete 删除权限
func PermissionDelete(c *gin.Context) {
	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取分类 ID
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var permission models.Permission

	if err := config.DB.First(&permission, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Permission not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&permission).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete permission", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Permission deleted successfully", nil)
}

// PermissionUpdate 更新角色
func PermissionUpdate(c *gin.Context) {
	// 定义结构体来接收请求体中的参数
	var request struct {
		ID          uint   `json:"id" binding:"required"` // 角色 ID
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	// 绑定并验证请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 查询角色是否存在
	var permission models.Permission
	if err := config.DB.First(&permission, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Permission not found", http.StatusNotFound)
		return
	}

	// 更新角色数据
	permission.Name = request.Name
	permission.Description = request.Description

	// 保存更新
	if err := config.DB.Save(&permission).Error; err != nil {
		utils.ErrorResponse(c, "Failed to update permission", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, "Permission updated successfully", permission)
}

// GetPermission 查询单个权限
func GetPermission(c *gin.Context) {
	var permission models.Permission
	if err := config.DB.First(&permission, c.Param("id")).Error; err != nil {
		utils.ErrorResponse(c, "Permission not found", http.StatusNotFound)
		return
	}

	utils.SuccessResponse(c, "Permission fetched successfully", permission)
}

// GetAllPermissions 查询所有权限
func GetAllPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		log.Printf("查询所有权限失败%v", err)
		fmt.Printf("查询所有权限失败%v", err)
		utils.ErrorResponse(c, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Roles fetched successfully", permissions)
}
