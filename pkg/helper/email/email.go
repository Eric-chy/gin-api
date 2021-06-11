package email

import (
	"crypto/tls"
	"gin-api/config"
	"gopkg.in/gomail.v2"
)

type Email struct {
	*SMTPInfo
}

type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}

func SendMail(to []string, subject, body string) error {
	cfg := config.Conf.Email
	email := NewEmail(&SMTPInfo{
		Host:     cfg.Host,
		Port:     cfg.Port,
		IsSSL:    cfg.IsSSL,
		UserName: cfg.UserName,
		Password: cfg.Password,
		From:     cfg.From,
	})
	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(email.Host, email.Port, email.UserName, email.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: email.IsSSL}
	return dialer.DialAndSend(m)
}
