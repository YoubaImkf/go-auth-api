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
		host:     viper.GetString("smtp.host"),
		port:     viper.GetInt("smtp.port"),
		username: viper.GetString("smtp.username"),
		password: viper.GetString("smtp.password"),
		from:     viper.GetString("smtp.from"),
	}
}

func (s *emailService) SendPasswordResetEmail(to, token string) error {
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	} else {
		auth = nil
	}

	resetURL := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", token)
	msg := []byte(fmt.Sprintf(
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
		s.from, to, resetURL))

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
