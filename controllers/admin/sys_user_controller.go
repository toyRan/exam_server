package admin

import (
	"exam_server/config"
	"exam_server/models"
	"exam_server/utils"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var tokenBlacklist = make(map[string]bool)

//// SysUserCreate 显示注册页面并生成验证码
//func SysUserCreate(c *gin.Context) {
//	captchaID := GenerateCaptchaID() // 使用封装好的生成函数
//	log.Printf(captchaID)
//	c.HTML(http.StatusOK, "register.html", gin.H{"CaptchaID": captchaID, "Title": "Register"})
//}

// SysUserStore 处理后台用户添加逻辑
func SysUserStore(c *gin.Context) {
	// 表单数据绑定
	var input struct {
		Username   string  `json:"username" binding:"required"`
		Email      string  `json:"email" binding:"required,email"`
		Password   string  `json:"password" binding:"required,min=6"`
		SysRoleIDs []int64 `json:"sys_role_ids"`                 // 角色ID数组
		Status     int64   `json:"status" validate:"oneof=0 1" ` // 使用状态，无需 'required' 标签
	}

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		// 打印并返回详细错误信息
		log.Println("Binding Error:", err)
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	// 检查后台用户是否已存在
	var user models.SysUser
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err == nil {
		utils.ErrorResponse(c, "该邮件已经存在", http.StatusBadRequest)
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.ErrorResponse(c, "密码加密失败", http.StatusInternalServerError)
		return
	}

	// 开启事务
	tx := config.DB.Begin()

	// 创建后台用户
	user = models.SysUser{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Status:   input.Status, // 使用状态
	}

	// 保存用户到数据库
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("创建用户失败: %v", err)
		utils.ErrorResponse(c, "创建用户失败", http.StatusInternalServerError)
		return
	}

	// 分配角色
	if len(input.SysRoleIDs) > 0 {
		var roles []models.SysRole
		if err := tx.Where("id IN ?", input.SysRoleIDs).Find(&roles).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, "部分角色未找到", http.StatusBadRequest)
			return
		}

		// 将角色分配给用户
		if err := tx.Model(&user).Association("SysRoles").Replace(&roles); err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, "用户角色分配失败", http.StatusInternalServerError)
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败: %v", err)
		utils.ErrorResponse(c, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	data := gin.H{
		"email":    user.Email,
		"username": user.Username,
	}

	// 返回成功响应
	utils.SuccessResponse(c, "后台用户创建成功", data)
}

// VerifySysUserEmail 验证邮箱
//func VerifySysUserEmail(c *gin.Context) {
//	var user models.SysUser
//	userID := c.Param("id")
//
//	if err := config.DB.First(&user, userID).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "SysUser not found"})
//		return
//	}
//
//	// 更新 email_verified_at 为当前时间
//	currentTime := time.Now()
//	user.EmailVerifiedAt = &currentTime
//
//	if err := config.DB.Save(&user).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
//}

// ActivateSysUser 通过激活令牌激活后台用户
func ActivateSysUser(c *gin.Context) {
	var input struct {
		Token string `json:"token" binding:"required"`
	}

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 假设激活令牌是一个简单的字符串，实际应用中应更复杂
	var user models.SysUser
	if err := config.DB.Where("activation_token = ?", input.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired activation token"})
		return
	}

	// 激活后台用户
	user.Status = 1
	//now := time.Now() // 创建一个中间变量来保存当前时间
	//user.EmailVerifiedAt = &now // 使用中间变量的地址赋值
	//user.ActivationToken = nil  // 清空激活令牌，避免重复使用

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SysUser activated successfully"})
}

// SysUserLogin - 后台用户登录
func SysUserLogin(c *gin.Context) {
	// 接收并绑定请求体的数据
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		//Captcha   string `json:"captcha" binding:"required"`
		//CaptchaID string `json:"captchaID" binding:"required"`
	}

	// 绑定并验证请求体
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}
	//
	//// 验证验证码
	//if !utils.VerifyCaptcha(input.CaptchaID, input.Captcha) {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid captcha"})
	//	return
	//}

	// 查找后台用户
	var user models.SysUser
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, "Invalid email", http.StatusUnauthorized)
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.ErrorResponse(c, "Invalid  password", http.StatusUnauthorized)
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.ErrorResponse(c, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 登录成功，返回 token
	data := gin.H{
		"token": token,
		"user": gin.H{
			"email":    user.Email,
			"username": user.Username,
		},
	}
	utils.SuccessResponse(c, "登录成功", data)
}

// RequestPasswordReset - 处理后台用户密码重置请求
func RequestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找后台用户是否存在
	var user models.SysUser
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "后台用户不存在"})
		return
	}

	// 生成密码重置 token
	resetToken := uuid.New().String()

	// 插入 token 到 password_reset_tokens 表
	resetTokenRecord := models.PasswordResetToken{
		Email:     input.Email,
		Token:     resetToken,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&resetTokenRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存重置密码请求"})
		return
	}

	// 生成重置密码链接
	resetLink := "http://your-domain.com/reset-password?token=" + resetToken + "&email=" + input.Email

	// 发送重置密码邮件
	err := utils.SendResetPasswordEmail(input.Email, resetLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送重置密码邮件失败"})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{"message": "重置密码的链接已经发送到您的邮箱"})
}

// ResetPassword - 处理后台用户密码重置
func ResetPassword(c *gin.Context) {
	var input struct {
		Password string `json:"password" binding:"required,min=6"`
		Token    string `json:"token" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找 password_reset_tokens 表中的记录
	var resetTokenRecord models.PasswordResetToken
	if err := config.DB.Where("email = ? AND token = ?", input.Email, input.Token).First(&resetTokenRecord).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或过期的重置密码链接"})
		return
	}

	// 验证 token 是否已过期（假设 token 有效期为 1 小时）
	if time.Since(resetTokenRecord.CreatedAt) > time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "重置密码链接已过期"})
		return
	}

	// 查找后台用户并更新密码
	var user models.SysUser
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "后台用户不存在"})
		return
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法加密密码"})
		return
	}

	// 更新后台用户密码
	user.Password = hashedPassword
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新密码"})
		return
	}

	// 删除 password_reset_tokens 表中的记录
	config.DB.Delete(&resetTokenRecord)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "密码已成功重置"})
}

// Logout - 退出登录
func Logout(c *gin.Context) {
	// 清除客户端的 JWT Token （实际由前端来处理）
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// SysRoleResponse 用于返回角色的完整信息
type SysRoleResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // 假设角色表包含名称字段
}

// SysUserResponse 定义返回的后台用户字段
type SysUserResponse struct {
	ID        int64             `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Status    int64             `json:"status"`
	CreatedAt *models.LocalTime `json:"created_at"`
	UpdatedAt *models.LocalTime `json:"updated_at"`
	SysRoles  []SysRoleResponse `json:"sys_roles"` // 用户的完整角色信息
}

func GetAllSysUsers(c *gin.Context) {
	var users []models.SysUser
	var userResponses []SysUserResponse

	// 获取查询参数
	query := c.DefaultQuery("query", "")
	currentPage, _ := strconv.Atoi(c.DefaultQuery("currentPage", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 计算偏移量
	offset := (currentPage - 1) * pageSize

	// 构建查询
	dbQuery := config.DB.Preload("SysRoles").Order("created_at DESC")

	// 如果有查询参数，按名称或邮箱模糊查询
	if query != "" {
		dbQuery = dbQuery.Where("username LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 执行查询并分页
	if err := dbQuery.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		log.Printf("获取管理员列表失败 %v\n", err)
		utils.ErrorResponse(c, "获取管理员列表失败", http.StatusInternalServerError)
		return
	}

	// 映射到 SysUserResponse
	for _, user := range users {
		var roles []SysRoleResponse
		for _, role := range user.SysRoles {
			roles = append(roles, SysRoleResponse{
				ID:   role.ID,
				Name: role.Name, // 假设角色包含名称字段
			})
		}

		userResponses = append(userResponses, SysUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			SysRoles:  roles, // 添加完整角色信息
		})
	}

	// 获取后台用户总数以计算总页数
	var total int64
	config.DB.Model(&models.SysUser{}).Count(&total)
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize)) // 计算总页数

	// 返回响应
	utils.SuccessResponse(c, "获取管理员列表成功", gin.H{
		"sys_users":   userResponses,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"totalPages":  totalPages,
	})
}

// SysUserUpdate 更新后台用户信息
func SysUserUpdate(c *gin.Context) {
	// 定义请求数据结构
	var request struct {
		ID         int64   `json:"id" binding:"required"`
		Email      string  `json:"email" binding:"required,email"`
		Username   string  `json:"username" binding:"required"`
		Password   string  `json:"password"`     // 可选：只有当提供时更新密码
		SysRoleIDs []int64 `json:"sys_role_ids"` // 角色ID数组
		Status     int64   `json:"status"`       // 使用状态  status
	}

	// 绑定并验证请求体中的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("绑定请求体失败: %v", err) // 打印详细的错误信息
		utils.ErrorResponse(c, fmt.Sprintf("Invalid request data: %v", err), http.StatusBadRequest)
		return
	}

	//在 Go 的 gin 框架中，binding:"required" 会检查字段是否为零值，对于 int64 类型的字段，零值即为 0。
	//如果 Status 字段值为 0，gin 会认为该字段未被赋值，从而导致验证失败。
	// 手动检查 Status 的合法性（如果需要）
	if request.Status != 0 && request.Status != 1 {
		utils.ErrorResponse(c, "Invalid status value", http.StatusBadRequest)
		return
	}

	// 开启事务
	tx := config.DB.Begin()

	// 查询用户是否存在
	var sysUser models.SysUser
	if err := tx.First(&sysUser, request.ID).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "SysUser not found", http.StatusNotFound)
		return
	}

	// 更新用户的基本信息
	sysUser.Email = request.Email
	sysUser.Username = request.Username
	sysUser.Status = request.Status

	// 如果提供了密码，则对其进行加密后更新
	if request.Password != "" {
		hashedPassword, err := utils.HashPassword(request.Password) // 假设 utils.HashPassword 是一个密码加密函数
		if err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		sysUser.Password = hashedPassword
	}

	// 更新用户的角色关联（多对多关系）
	if len(request.SysRoleIDs) > 0 {
		var sysRoles []models.SysRole
		if err := tx.Where("id IN ?", request.SysRoleIDs).Find(&sysRoles).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, "Some roles not found", http.StatusBadRequest)
			return
		}
		// 更新多对多关系
		if err := tx.Model(&sysUser).Association("SysRoles").Replace(&sysRoles); err != nil {
			tx.Rollback()
			log.Printf("更新用户角色失败 %v", err)
			utils.ErrorResponse(c, "Failed to update user roles", http.StatusInternalServerError)
			return
		}
	}

	// 保存用户信息
	if err := tx.Save(&sysUser).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("提交事务失败 %v", err)
		utils.ErrorResponse(c, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	utils.SuccessResponse(c, "SysUser updated successfully", sysUser)
}

// SysUserDelete 删除用户
func SysUserDelete(c *gin.Context) {
	var request struct {
		ID int64 `json:"uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, "Invalid request", http.StatusBadRequest)
		return
	}

	// 获取当前登录用户 ID，假设已从中间件中设置了用户 ID
	currentUserID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 防止用户删除自己
	if request.ID == currentUserID.(int64) {
		log.Println("You cannot delete your own account")
		utils.ErrorResponse(c, "You cannot delete your own account", http.StatusForbidden)
		return
	}

	// 检查用户是否存在
	var user models.SysUser
	if err := config.DB.First(&user, request.ID).Error; err != nil {
		utils.ErrorResponse(c, "User not found", http.StatusNotFound)
		return
	}

	// 执行删除操作
	if err := config.DB.Delete(&user).Error; err != nil {
		utils.ErrorResponse(c, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	//// 操作日志记录（可选）
	//// 假设有一个日志记录表，记录用户操作
	//if err := logAction(currentUserID.(uint), "delete_user", request.ID); err != nil {
	//	utils.ErrorResponse(c, "Failed to record action log", http.StatusInternalServerError)
	//	return
	//}

	utils.SuccessResponse(c, "SysUser deleted successfully", nil)
}
