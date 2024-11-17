package front

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 检查角色名称是否已经存在
	var existingRole models.Role
	if err := config.DB.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		// 如果找到相同名称的角色，返回错误响应
		log.Printf("角色 %s already exists", role.Name)
		utils.ErrorResponse(c, "Role with the same name already exists", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to create role", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Role created successfully", role)
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {

	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取分类 ID
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var role models.Role

	if err := config.DB.First(&role, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Role not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Role deleted successfully", nil)
}

// RoleUpdate 更新角色
func RoleUpdate(c *gin.Context) {
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
	var role models.Role
	if err := config.DB.First(&role, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Role not found", http.StatusNotFound)
		return
	}

	// 更新角色数据
	role.Name = request.Name
	role.Description = request.Description

	// 保存更新
	if err := config.DB.Save(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to update role", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, "Role updated successfully", role)
}

// GetRole  查询单个角色
func GetRole(c *gin.Context) {
	var role models.Role
	if err := config.DB.First(&role, c.Param("id")).Error; err != nil {
		utils.ErrorResponse(c, "Role not found", http.StatusNotFound)
		return
	}

	utils.SuccessResponse(c, "Role fetched successfully", role)
}

// GetAllRoles 查询所有角色
func GetAllRoles(c *gin.Context) {

	var roles []models.Role

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
		dbQuery = dbQuery.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 执行查询并分页
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
		log.Printf("获取角色列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取角色列表失败", http.StatusInternalServerError)
		return
	}

	if err := config.DB.Find(&roles).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	// 获取用户总数以计算总页数
	var total int64
	config.DB.Model(&models.User{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取角色列表成功", gin.H{
		"users":       roles,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
	})
}

// SetRolePermissions 设置角色的权限
func SetRolePermissions(c *gin.Context) {
	var request struct {
		Permissions []int64 `json:"permissions" binding:"required"` // 权限ID列表
	}

	// 获取并转换角色ID
	roleIDStr := c.Param("id")                         // 获取字符串类型的角色ID
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64) // 将字符串转换为 uint64 类型
	if err != nil {
		utils.ErrorResponse(c, "Invalid role ID", http.StatusBadRequest)
		return
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 删除当前角色的所有权限
	if err := config.DB.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		log.Printf("删除当前角色的所有权限 %v", err)
		utils.ErrorResponse(c, "Failed to clear current permissions", http.StatusInternalServerError)
		return
	}

	// 分配新权限
	for _, permissionID := range request.Permissions {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		if err := config.DB.Create(&rolePermission).Error; err != nil {
			utils.ErrorResponse(c, "Failed to assign permissions", http.StatusInternalServerError)
			return
		}
	}

	utils.SuccessResponse(c, "Permissions updated successfully", nil)
}
