package admin

import (
	"errors"
	"exam_server/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"exam_server/models"
	"exam_server/utils"
)

// UserRequest 用户请求结构体
type UserRequest struct {
	ID        uint   `json:"id,omitempty"` // 添加 ID 字段，在更新时使用
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email" binding:"required,email"`
	Status    int64  `json:"status" binding:"oneof=0 1"` // 状态字段: 0=禁用, 1=启用
	RoleID    int64  `json:"role_id" binding:"required"` // 改为 RoleID
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 验证角色是否存在
	var role models.Role
	if err := config.DB.First(&role, req.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(c, "选择的角色不存在", http.StatusBadRequest)
			return
		}
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logrus.WithError(err).Error("Failed to hash password")
		utils.ErrorResponse(c, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  hashedPassword,
		Email:     req.Email,
		Status:    req.Status,
		RoleID:    req.RoleID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		logrus.WithError(err).Error("Failed to create user")
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	// 清除密码后返回用户信息
	user.Password = ""
	utils.SuccessResponse(c, "User created successfully", user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("Failed to parse user update request")
		utils.ErrorResponse(c, "参数错误"+err.Error(), http.StatusBadRequest)
		return
	}
	// 验证用户ID是否提供
	if req.ID == 0 {
		utils.ErrorResponse(c, "User ID is required", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      req.Email,
		"role_id":    req.RoleID, // 更新用户角色
		"status":     req.Status,
	}

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logrus.WithError(err).Error("Failed to hash password")
			utils.ErrorResponse(c, "Failed to process password", http.StatusInternalServerError)
			return
		}
		updates["password"] = hashedPassword
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		logrus.WithError(err).Error("Failed to update user")
		utils.ErrorResponse(c, "Failed to update user"+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "User updated successfully", nil)
}

//// GetUserList 获取用户列表
//func GetUserList(c *gin.Context) {
//	var users []models.User
//
//	err := config.DB.Select("id, username, email, role, created_at, updated_at").Find(&users).Error
//	if err != nil {
//		logrus.WithError(err).Error("Failed to get user list")
//		utils.ErrorResponse(c, "Failed to get user list", http.StatusInternalServerError)
//		return
//	}
//
//	utils.SuccessResponse(c, "User list retrieved successfully", users)
//}

// UserResponse 用户响应结构体
type UserResponse struct {
	ID        int64             `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	RoleID    int64             `json:"role_id"`
	RoleName  string            `json:"role_name"`
	Status    int64             `json:"status"`
	CreatedAt *models.LocalTime `json:"created_at"`
	UpdatedAt *models.LocalTime `json:"updated_at"`
}

// GetUsersPaginated 获取分页用户列表
func GetUsersPaginated(c *gin.Context) {
	// 获取查询参数
	queryStr := c.DefaultQuery("query", "")
	pageStr := c.DefaultQuery("currentPage", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 分页查询
	var users []models.User
	var total int64

	// 构建查询，添加 Preload("Role") 预加载角色数据
	query := config.DB.Model(&models.User{}).Preload("Role")

	// 如果 query 参数存在，则进行模糊搜索
	if queryStr != "" {
		searchPattern := "%" + queryStr + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	query.Order("id DESC")

	// 获取总数
	query.Count(&total)

	// 进行分页
	offset := (page - 1) * pageSize
	query.Limit(pageSize).Offset(offset).Find(&users)

	// 转换响应数据
	var responseUsers []UserResponse
	for _, user := range users {
		responseUser := UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			RoleID:    user.RoleID,
			RoleName:  user.Role.Name, // 现在可以访问 Role.Name
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		responseUsers = append(responseUsers, responseUser)
	}

	// 构建响应数据
	response := gin.H{
		"users":       responseUsers,
		"total":       total,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalPage":   (total + int64(pageSize) - 1) / int64(pageSize),
	}

	utils.SuccessResponse(c, "User list retrieved successfully", response)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, "无效的用户ID", http.StatusBadRequest)
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := config.DB.First(&user, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(c, "用户不存在", http.StatusNotFound)
			return
		}
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	// 执行软删除
	if err := config.DB.Delete(&user).Error; err != nil {
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	utils.SuccessResponse(c, "用户删除成功", nil)
}

// GetRoleOptions 获取角色选项列表
func GetRoleOptions(c *gin.Context) {
	var roles []models.Role

	// 只获取启用状态的角色
	if err := config.DB.Where("status = ?", 1).
		Order("sort asc").
		Find(&roles).Error; err != nil {
		message, statusCode := utils.HandleMySQLError(err)
		utils.ErrorResponse(c, message, statusCode)
		return
	}

	// 转换为前端需要的格式
	var options []map[string]interface{}
	for _, role := range roles {
		options = append(options, map[string]interface{}{
			"value": role.ID,
			"label": role.Name,
		})
	}

	utils.SuccessResponse(c, "获取角色列表成功", options)
}
