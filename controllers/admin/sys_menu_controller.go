package admin

import (
	"exam_server/services"
	"exam_server/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// GetSysUserMenus 获取当前后台用户的菜单
func GetSysUserMenus(c *gin.Context) {
	// 从上下文中获取 userID
	userID, exists := c.Get("userID")
	log.Printf("当前登录的userID: %d", userID)
	if !exists {
		utils.ErrorResponse(c, "后台用户不存在", http.StatusUnauthorized)
		return
	}

	// 将 interface{} 转换为 uint
	userIdToInt64, ok := userID.(int64)
	if !ok {
		log.Printf("userID 转换失败: %v", userID)
		utils.ErrorResponse(c, "Invalid user ID format", http.StatusInternalServerError)
		return
	}

	// 从数据库获取菜单数据
	menus, err := services.GetUserSysMenus(userIdToInt64)
	if err != nil {
		log.Printf("Failed to fetch menus: %v", err)
		utils.ErrorResponse(c, "Failed to fetch menus", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, "Menus fetched successfully", menus)
}
