package front

import (
	"log"
	"net/http"

	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserRegisterRequest 用户注册请求结构体
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, "绑定参数失败", http.StatusBadRequest)
		return
	}

	// 检查用户是否已存在（邮箱检查）
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, "Email已经注册", http.StatusBadRequest)
		return
	}

	// 检查用户名是否已存在
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, "用户名已被使用", http.StatusBadRequest)
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 获取用户 IP
	clientIP := c.ClientIP()

	// 创建新用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Status:    0,
		IpAddress: clientIP,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, "Failed to create user", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "恭喜你，注册成功", gin.H{
		"user": gin.H{
			// "id":         user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

// UserLoginRequest 用户登录请求结构体
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, "Failed to validate login data", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := config.DB.Preload("Role").Where("email = ?", req.Email).First(&user).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		utils.ErrorResponse(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check user status
	if user.Status != 1 {
		utils.ErrorResponse(c, "Account is disabled", http.StatusForbidden)
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Error comparing password: %v", err)
		utils.ErrorResponse(c, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 生成JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Login successful", gin.H{
		"token": token,
		"user": gin.H{
			//"id":         user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

// UserLogout 用户退出
func UserLogout(c *gin.Context) {
	// TODO: 实现用户退出逻辑
	// 1. 清除token（如果需要）
	c.JSON(http.StatusOK, gin.H{
		"message": "退出成功",
	})
}

// GetUserProfile 获取用户信息
func GetUserProfile(c *gin.Context) {
	// TODO: 实现获取用户信息逻辑
	// 1. 从JWT中获取用户ID
	// 2. 查询用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户信息成功",
		// 添加用户信息
	})
}

// UpdateProfileRequest 更新用户信息请求结构体
type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=1,max=50"`
	LastName    string `json:"last_name" binding:"required,min=1,max=50"`
	OldPassword string `json:"old_password,omitempty"`                           // 可选，但如果要更新密码则必填
	NewPassword string `json:"new_password,omitempty" binding:"omitempty,min=6"` // 可选，但如果填写则至少6位
}

// UpdateProfile 更新用户信息
func UpdateProfile(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("更新用户信息失败，错误信息为：%v", err)
		utils.ErrorResponse(c, "Invalid input data", http.StatusBadRequest)
		return
	}

	// 获取当前用户信息
	var currentUser models.User
	if err := config.DB.First(&currentUser, userID).Error; err != nil {
		utils.ErrorResponse(c, "User not found", http.StatusNotFound)
		return
	}

	// // 检查邮箱是否被其他用户使用
	// var existingUser models.User
	// if err := config.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
	// 	utils.ErrorResponse(c, "Email already in use by another user", http.StatusBadRequest)
	// 	return
	// }

	// 处理密码更新
	if req.NewPassword != "" {
		// 如果要更新密码，必须提供旧密码
		if req.OldPassword == "" {
			utils.ErrorResponse(c, "Old password is required to update password", http.StatusBadRequest)
			return
		}

		// 验证旧密码
		if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.OldPassword)); err != nil {
			utils.ErrorResponse(c, "Invalid old password", http.StatusBadRequest)
			return
		}

		// 生成新密码的哈希值
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			utils.ErrorResponse(c, "Failed to process password update", http.StatusInternalServerError)
			return
		}
		currentUser.Password = string(hashedPassword)
	}

	// 更新用户信息
	updates := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		// Email:     req.Email,
		Password: currentUser.Password, // 如果密码已更新，使用新密码；否则保持原密码
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		utils.ErrorResponse(c, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// 获取更新后的用户信息
	var updatedUser models.User
	if err := config.DB.First(&updatedUser, userID).Error; err != nil {
		utils.ErrorResponse(c, "Failed to fetch updated profile", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", gin.H{
		"user": gin.H{
			//"id":         updatedUser.ID,
			"first_name": updatedUser.FirstName,
			"last_name":  updatedUser.LastName,
			"email":      updatedUser.Email,
		},
	})
}
