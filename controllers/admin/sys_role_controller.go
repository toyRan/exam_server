package admin

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// SysRoleCreate 创建管理员角色
func SysRoleCreate(c *gin.Context) {
	var role models.SysRole
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 检查管理员角色名称是否已经存在
	var existingSysRole models.SysRole
	if err := config.DB.Where("name = ?", role.Name).First(&existingSysRole).Error; err == nil {
		// 如果找到相同名称的管理员角色，返回错误响应
		log.Printf("管理员角色 %s already exists", role.Name)
		utils.ErrorResponse(c, "SysRole with the same name already exists", http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to create role", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "SysRole created successfully", role)
}

// SysRoleDelete 删除管理员角色
func SysRoleDelete(c *gin.Context) {

	var request struct {
		ID uint `json:"id" binding:"required"` // 从请求体中获取分类 ID
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	var role models.SysRole

	if err := config.DB.First(&role, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "SysRole not found", http.StatusNotFound)
		return
	}

	if err := config.DB.Delete(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "SysRole deleted successfully", nil)
}

// SysRoleUpdate 更新管理员角色
func SysRoleUpdate(c *gin.Context) {
	// 定义结构体来接收请求体中的参数
	var request struct {
		ID          uint   `json:"id" binding:"required"` // 管理员角色 ID
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	// 绑定并验证请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 查询管理员角色是否存在
	var role models.SysRole
	if err := config.DB.First(&role, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "SysRole not found", http.StatusNotFound)
		return
	}

	// 更新管理员角色数据
	role.Name = request.Name
	role.Description = request.Description

	// 保存更新
	if err := config.DB.Save(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to update role", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, "SysRole updated successfully", role)
}

// GetSysRole  查询单个管理员角色
func GetSysRole(c *gin.Context) {
	var role models.SysRole
	if err := config.DB.First(&role, c.Param("id")).Error; err != nil {
		utils.ErrorResponse(c, "该管理员角色不存在", http.StatusNotFound)
		return
	}

	utils.SuccessResponse(c, "获取成功", role)
}

// GetAllSysRoles 查询所有管理员角色 （无分页）
func GetAllSysRoles(c *gin.Context) {

	var sysRoles []models.SysRole

	//查询出所有的 sys_roles 不分页
	if err := config.DB.Find(&sysRoles).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	// 返回响应
	utils.SuccessResponse(c, "获取管理员角色列表成功", gin.H{
		"sys_roles": sysRoles,
	})
}

// SetSysRolePermissions 设置管理员角色的权限
func SetSysRolePermissions(c *gin.Context) {
	var request struct {
		SysPermissionIDs []int64 `json:"sys_permission_ids" binding:"required"` // 权限ID列表
		SysRoleID        int64   `json:"sys_role_id" binding:"required"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data", http.StatusBadRequest)
		return
	}

	// 删除当前管理员角色的所有权限
	if err := config.DB.Debug().Where("sys_role_id = ?", uint(request.SysRoleID)).Delete(&models.SysRolePermission{}).Error; err != nil {
		log.Printf("删除当前管理员角色的所有权限 %v", err)
		utils.ErrorResponse(c, "Failed to clear current permissions", http.StatusInternalServerError)
		return
	}

	// 使用 map 去重权限ID
	uniquePermissionIDs := make(map[int64]struct{})
	for _, permissionID := range request.SysPermissionIDs {
		uniquePermissionIDs[permissionID] = struct{}{}
	}

	// 分配去重后的新权限
	for permissionID := range uniquePermissionIDs { // 注意：这里的 permissionID 是 map 的键，类型为 int64
		rolePermission := models.SysRolePermission{
			SysRoleID:       request.SysRoleID,
			SysPermissionID: permissionID,
		}
		if err := config.DB.Create(&rolePermission).Error; err != nil {
			log.Printf("给管理员角色分配权限失败 %v", err)
			utils.ErrorResponse(c, "Failed to assign permissions", http.StatusInternalServerError)
			return
		}
	}

	utils.SuccessResponse(c, "Permissions updated successfully", nil)
}

// GetAllSysRolesPaginated 查询所有管理员角色带分页
func GetAllSysRolesPaginated(c *gin.Context) {

	var roles []models.SysRole

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
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
		log.Printf("获取管理员角色列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取管理员角色列表失败", http.StatusInternalServerError)
		return
	}

	// 获取用户总数以计算总页数
	var total int64
	config.DB.Model(&models.User{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取管理员角色列表成功", gin.H{
		"sys_roles":   roles,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
	})
}
