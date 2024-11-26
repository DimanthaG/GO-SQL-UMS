
package email

import (
    "fmt"
    "net/smtp"
)

func SendEmail(to, subject, body string) error {
    from := "your_email@example.com"
    password := "your_email_password"
    smtpHost := "smtp.example.com"
    smtpPort := "587"

    message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

    auth := smtp.PlainAuth("", from, password, smtpHost)
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
}
