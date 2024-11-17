package utils

import (
	"fmt"
	"log"

	//"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/smtp"
	"os"
)

//// SendResetPasswordEmail - 发送重置密码的邮件
//func SendResetPasswordEmail(toEmail string, resetLink string) error {
//	from := mail.NewEmail("Your App Name", "noreply@mydomain.com")
//	subject := "Reset Your Password"
//	to := mail.NewEmail("Recipient", toEmail)
//	plainTextContent := fmt.Sprintf("Click the link to reset your password: %s", resetLink)
//	htmlContent := fmt.Sprintf("<strong>Click the link to reset your password: </strong><a href=\"%s\">Reset Password</a>", resetLink)
//	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
//	client := sendgrid.NewSendClient("your_sendgrid_api_key")
//	response, err := client.Send(message)
//	if err != nil {
//		return err
//	}
//	fmt.Println(response.StatusCode, response.Body, response.Headers)
//	return nil
//}

// sendActivationEmail 发送激活邮件
func sendActivationEmail(toEmail string, token string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// 确保所有配置项都已加载
	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP configuration is missing")
	}

	to := []string{toEmail}

	// 构建激活链接
	appDomain := os.Getenv("APP_DOMAIN")
	if appDomain == "" {
		return fmt.Errorf("网站域名未配置")
	}
	activationLink := fmt.Sprintf("%s/activate?token=%s", appDomain, token)

	message := []byte("Subject: Account Activation\n" +
		"To: " + toEmail + "\n" +
		"Please click the following link to activate your account: " + activationLink + "\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Printf("发送给%s的激活链接邮件失败. %v\n", toEmail, err)
		return err
	}
	return nil
}

// SendResetPasswordEmail - 发送重置密码的邮件
func SendResetPasswordEmail(toEmail string, resetLink string) error {
	from := os.Getenv("SMTP_USER")         // 发件人邮箱
	password := os.Getenv("SMTP_PASSWORD") // 发件人邮箱密码或 SMTP 授权码
	smtpHost := os.Getenv("SMTP_HOST")     // SMTP 服务器地址
	smtpPort := os.Getenv("SMTP_PORT")     // SMTP 端口号

	// 设置邮件内容
	subject := "Subject: Reset Your Password\n" // 邮件标题
	body := fmt.Sprintf("Click the link to reset your password: %s", resetLink)
	message := []byte(subject + "\n" + body)

	// 设置发件人和收件人
	auth := smtp.PlainAuth("", from, password, smtpHost)
	to := []string{toEmail}

	// 拼接 SMTP 服务器地址和端口
	smtpAddress := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// 发送邮件
	err := smtp.SendMail(smtpAddress, auth, from, to, message)
	if err != nil {
		log.Printf("发送给%s的邮件失败. %v\n", toEmail, err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
