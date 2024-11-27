package email

import (
	"log"
	"net/smtp"
)

func SendEmail(to, subject, body string) {
	from := "your_email@example.com"  // Replace with your email
	password := "your_email_password" // Replace with your email password

	smtpHost := "smtp.example.com" // Replace with your SMTP server
	smtpPort := "587"

	message := []byte("Subject: " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return
	}

	log.Println("Email sent successfully to", to)
}
