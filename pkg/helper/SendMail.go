package helper

import (
	"os"

	"github.com/go-gomail/gomail"
)

func SendResetEmail(to, resetLink string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/html", "Click the link to reset your password: <a href='"+resetLink+"'>Reset Password</a>")

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)
	return d.DialAndSend(m)
}