// captcha.go
package utils

import (
	"github.com/dchest/captcha"
	"net/http"
)

// GenerateCaptchaID 生成一个新的验证码ID
func GenerateCaptchaID() string {
	return captcha.New()
}

// GenerateCaptchaImage 根据captchaID生成验证码图片并写入响应
func GenerateCaptchaImage(w http.ResponseWriter, captchaID string) error {
	w.Header().Set("Content-Type", "image/png")
	return captcha.WriteImage(w, captchaID, 240, 80)
}

// VerifyCaptcha 验证用户输入的验证码是否正确
func VerifyCaptcha(captchaID, captchaValue string) bool {
	return captcha.VerifyString(captchaID, captchaValue)
}
