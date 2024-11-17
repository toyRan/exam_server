package middlewares

import (
	"exam_server/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 验证 JWT 并从中提取用户信息（必需的认证）
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 信息
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenString == "" {
			utils.ErrorResponse(c, "Authorization token is missing", http.StatusUnauthorized)
			c.Abort() // 中止请求
			return
		}

		// 验证和解析 JWT
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(c, "Invalid or expired token", http.StatusUnauthorized)
			c.Abort() // 中止请求
			return
		}

		// 将 userID 存入 context
		c.Set("userID", claims.UserID)
		log.Printf("jwt里解析出来的用户id为%d", claims.UserID)

		c.Next() // 继续处理请求
	}
}

// OptionalJWTAuth 可选的 JWT 认证中间件
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 信息
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenString != "" {
			// 验证和解析 JWT
			claims, err := utils.ValidateJWT(tokenString)
			if err == nil {
				// 将 userID 存入 context
				c.Set("userID", claims.UserID)
				log.Printf("可选jwt里解析出来的用户id为%d", claims.UserID)
			}
		}

		c.Next() // 无论是否有 token，都继续处理请求
	}
}
