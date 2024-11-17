// auth.go
package utils

import (
	"exam_server/config"
	"exam_server/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CheckUserRole 检查用户角色
func CheckUserRole(userID interface{}, roleCode string) bool {
	if userID == nil {
		return false
	}

	var user models.User
	err := config.DB.Preload("Role").
		Select("users.id, users.role_id, roles.code, roles.status").
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.id = ? AND roles.status = ?", userID, 1).
		First(&user).Error

	if err != nil {
		logrus.WithError(err).
			WithFields(logrus.Fields{
				"userID": userID,
				"role":   roleCode,
			}).
			Error("Failed to query user role")
		return false
	}

	// 添加日志以便调试
	logrus.WithFields(logrus.Fields{
		"userID":    userID,
		"roleCode":  roleCode,
		"userRole":  user.Role.Code,
		"isMatched": user.Role.Code == roleCode,
	}).Info("Role check result")

	return user.Role.Code == roleCode
}
