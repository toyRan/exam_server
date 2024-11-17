package admin

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleRequest 角色请求结构体
type RoleRequest struct {
	ID          int64  `json:"id"`                                      // 创建时不需要，更新时必填
	Name        string `json:"name" binding:"required,max=50"`          // 角色名称，必填，最大50字符
	Code        string `json:"code" binding:"required,max=50,alphanum"` // 角色编码，必填，最大50字符，只允许字母和数字
	Status      int64  `json:"status" binding:"oneof=0 1"`              // 状态，必填，只允许0或1
	Sort        int64  `json:"sort" binding:"gte=0"`                    // 排序，必须大于等于0
	Description string `json:"description" binding:"omitempty,max=255"` // 描述，选填，最大255字符
}

// CreateRole 创建管理员角色
func CreateRole(c *gin.Context) {
	var request RoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 检查角色名称是否已经存在
	var existingRole models.Role
	if err := config.DB.Where("name = ?", request.Name).First(&existingRole).Error; err == nil {
		log.Printf("角色名称 %s already exists", request.Name)
		utils.ErrorResponse(c, "Role with the same name already exists", http.StatusBadRequest)
		return
	}

	// 检查角色编码是否已经存在
	if err := config.DB.Where("code = ?", request.Code).First(&existingRole).Error; err == nil {
		log.Printf("角色编码 %s already exists", request.Code)
		utils.ErrorResponse(c, "Role with the same code already exists", http.StatusBadRequest)
		return
	}

	// 创建角色
	role := models.Role{
		Name:        request.Name,
		Code:        request.Code,
		Status:      request.Status,
		Sort:        request.Sort,
		Description: request.Description,
	}

	if err := config.DB.Create(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to create role", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Role created successfully", role)
}

// DeleteRole 删除管理员角色
func DeleteRole(c *gin.Context) {
	var request struct {
		ID int64 `json:"id" binding:"required,gt=0"` // ID必填且大于0
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data: "+err.Error(), http.StatusBadRequest)
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

// UpdateRole 更新管理员角色
func UpdateRole(c *gin.Context) {
	var request RoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 更新时ID必填
	if request.ID <= 0 {
		utils.ErrorResponse(c, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var role models.Role
	if err := config.DB.First(&role, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "Role not found", http.StatusNotFound)
		return
	}

	// 检查更新的名称是否与其他角色重复
	if err := config.DB.Where("name = ? AND id != ?", request.Name, request.ID).First(&models.Role{}).Error; err == nil {
		utils.ErrorResponse(c, "Role name already exists", http.StatusBadRequest)
		return
	}

	// 检查更新的编码是否与其他角色重复
	if err := config.DB.Where("code = ? AND id != ?", request.Code, request.ID).First(&models.Role{}).Error; err == nil {
		utils.ErrorResponse(c, "Role code already exists", http.StatusBadRequest)
		return
	}

	// 更新管理员角色数据
	role.Name = request.Name
	role.Code = request.Code
	role.Status = request.Status
	role.Sort = request.Sort
	role.Description = request.Description

	// 保存更新
	if err := config.DB.Save(&role).Error; err != nil {
		utils.ErrorResponse(c, "Failed to update role", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, "Role updated successfully", role)
}

// GetRole  查询单个管理员角色
func GetRole(c *gin.Context) {
	var role models.Role
	if err := config.DB.First(&role, c.Param("id")).Error; err != nil {
		utils.ErrorResponse(c, "该管理员角色不存在", http.StatusNotFound)
		return
	}

	utils.SuccessResponse(c, "获取成功", role)
}

// GetRoleList 查询所有管理员角色 （无分页）
func GetRoleList(c *gin.Context) {

	var roles []models.Role

	//查询出所有的 roles 不分页
	if err := config.DB.Find(&roles).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	// 返回响应
	utils.SuccessResponse(c, "获取管理员角色列表成功", gin.H{
		"roles": roles,
	})
}

// GetRolesPaginated 查询所有管理员角色带分页
func GetRolesPaginated(c *gin.Context) {

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
	config.DB.Model(&models.Role{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取管理员角色列表成功", gin.H{
		"roles":       roles,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
	})
}
