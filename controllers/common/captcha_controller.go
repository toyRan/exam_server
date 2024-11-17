package admin

import (
	"bytes"
	"encoding/base64"
	"exam_server/utils"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// GetCaptcha 生成验证码图片
func GetCaptcha(c *gin.Context) {
	captchaID := c.Param("captchaID")

	if captchaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Captcha ID is required"})
		return
	}

	// 设置响应头为图片格式
	c.Writer.Header().Set("Content-Type", "image/png")

	// 生成并写入验证码图片
	if err := utils.GenerateCaptchaImage(c.Writer, captchaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate captcha image"})
		return
	}
}

// RefreshCaptcha - 生成并返回新的 captchaID 和验证码图片的 URL
func RefreshCaptchaOld(c *gin.Context) {
	captchaID := utils.GenerateCaptchaID()
	c.JSON(http.StatusOK, gin.H{
		"captchaID":  captchaID,
		"captchaURL": "/api/v1/captcha/" + captchaID,
	})
}

// RefreshCaptcha - 生成并返回新的 captchaID 和 Base64 编码的图片数据
func RefreshCaptcha(c *gin.Context) {

	captchaID := captcha.NewLen(4) // 生成一个新的验证码 ID

	// 生成验证码图片，并将其编码为 Base64 格式
	var imgBuffer bytes.Buffer
	if err := captcha.WriteImage(&imgBuffer, captchaID, 240, 100); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate captcha image"})
		return
	}

	// 将图片数据转换为 Base64 格式
	encodedImage := base64.StdEncoding.EncodeToString(imgBuffer.Bytes())

	// 返回 JSON 响应，包含 captchaID 和 Base64 图片数据
	c.JSON(http.StatusOK, gin.H{
		"captchaID":    captchaID,
		"captchaImage": "data:image/png;base64," + encodedImage, // 这就是前端可以直接展示的图片
	})
}
