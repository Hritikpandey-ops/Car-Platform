package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendVerificationEmail(toEmail, token string) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	link := fmt.Sprintf("http://localhost:8081/verify?token=%s", token)
	body := fmt.Sprintf("Click the link to verify your email: %s", link)

	msg := []byte("Subject: Email Verification\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	}
}
