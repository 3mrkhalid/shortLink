package utils

import (
	"fmt"
	"net/smtp"
)

func SendResetEmail(to string, resetLink string) error {

	from := "yourEmail@gmail.com"
	password := "your_app_password"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(
		"Subject: Reset Password\r\n" +
			"\r\n" +
			fmt.Sprintf("Click here to reset your password:\n%s", resetLink),
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)
}