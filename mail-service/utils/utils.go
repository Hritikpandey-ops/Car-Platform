package utils

import (
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(to, subject, body string) error {
	port := 587

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_USERNAME"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), port, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))

	return d.DialAndSend(m)
}
