package service

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

type EmailService interface {
	SendPasswordResetEmail(to, token string) error
}

type emailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailService() EmailService {
	return &emailService{
		host:     viper.GetString("SMTP_HOST"),
		port:     viper.GetInt("SMTP_PORT"),
		username: viper.GetString("SMTP_USERNAME"),
		password: viper.GetString("SMTP_PASSWORD"),
		from:     viper.GetString("SMTP_FROM"),
	}
}

func (s *emailService) SendPasswordResetEmail(to, token string) error {
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	} else {
		auth = nil
	}

	host := viper.GetString("APP_HOST")
	protocol := "http"

	if viper.GetString("APP_ENVIRONMENT") == "production" {
		protocol = "https"
	}

	resetURL := fmt.Sprintf("%s://%s/reset-password?token=%s", protocol, host, token)
	msg := fmt.Appendf(nil,
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Fake Password Reset\r\n"+
			"\r\n"+
			"WARNING: You just have to add the token to the field in reset-password with your new password on Swagger\r\n"+
			"\r\n"+
			"Hello,\r\n\r\n"+
			"We received a request to reset your password. Please click the link below to reset your password:\r\n\r\n"+
			"%s\r\n\r\n"+
			"If you did not request a password reset, please ignore this email.\r\n\r\n"+
			"Thank you,\r\n"+
			"Your Team",
		s.from, to, resetURL)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
