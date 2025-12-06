package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailSender struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

var emailSender *EmailSender

func InitEmailSender() {
	// Get Mailtrap credentials from environment variables
	// These can be set via .env file or system environment
	host := os.Getenv("MAILTRAP_HOST")
	if host == "" {
		host = "sandbox.smtp.mailtrap.io"
	}

	port := os.Getenv("MAILTRAP_PORT")
	if port == "" {
		port = "2525"
	}

	username := os.Getenv("MAILTRAP_USERNAME")
	password := os.Getenv("MAILTRAP_PASSWORD")
	from := os.Getenv("MAILTRAP_FROM")
	if from == "" {
		from = "noreply@taskqueue.local"
	}

	emailSender = &EmailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}

	if username != "" && password != "" {
		log.Printf("✅ Email sender initialized (SMTP credentials configured)")
	} else {
		log.Println("⚠️  Warning: Mailtrap credentials not set. Email will be logged only.")
		log.Println("   Set MAILTRAP_USERNAME and MAILTRAP_PASSWORD to enable real email sending.")
	}
}

func SendEmail(to, subject, body string) error {
	if emailSender == nil {
		InitEmailSender()
	}

	// Log the email being sent
	log.Printf("📧 Email Task: TO=%s, SUBJECT=%s", to, subject)
	log.Printf("   BODY=%s", body)

	// If credentials are not set, just log (demo mode)
	if emailSender.Username == "" || emailSender.Password == "" {
		log.Println("   [DEMO MODE] Email logged but not actually sent (credentials not configured)")
		return nil
	}

	// Send via SMTP (Mailtrap)
	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		emailSender.From,
		to,
		subject,
		body,
	)

	auth := smtp.PlainAuth("", emailSender.Username, emailSender.Password, emailSender.Host)
	addr := fmt.Sprintf("%s:%s", emailSender.Host, emailSender.Port)

	err := smtp.SendMail(addr, auth, emailSender.From, []string{to}, []byte(message))
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("✅ Email sent successfully to %s via %s", to, emailSender.Host)
	return nil
}
