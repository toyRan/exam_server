// jwt.go
package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtKey = []byte("your_secret_key") // 应该从配置文件或环境变量中读取
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Claims 定义JWT的声明
type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT Token
func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour) // Token 有效期为72小时
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT 验证并解析JWT Token
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}

	return claims, nil
}
