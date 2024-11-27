package email

import (
	"errors"
	"log"
	"net/smtp"
	"os"
)

// SendEmail sends an email to the specified recipient
func SendEmail(to, subject, body string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return errors.New("email environment variables not set")
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}
