package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPConfig SMTP 配置信息
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// EmailMessage 邮件内容
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// SendEmail 发送邮件
func SendEmail(config *SMTPConfig, msg *EmailMessage) error {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	contentType := "text/plain; charset=UTF-8"
	if msg.IsHTML {
		contentType = "text/html; charset=UTF-8"
	}

	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s\r\n\r\n%s",
		config.From,
		strings.Join(msg.To, ","),
		msg.Subject,
		contentType,
		msg.Body,
	)

	if config.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: config.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, config.Host)
		if err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}

		if err = client.Mail(config.From); err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}

		for _, to := range msg.To {
			if err = client.Rcpt(to); err != nil {
				return fmt.Errorf("SMTP: %w", err)
			}
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}

		_, err = w.Write([]byte(body))
		if err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}

		err = w.Close()
		if err != nil {
			return fmt.Errorf("SMTP: %w", err)
		}

		return client.Quit()
	}

	return smtp.SendMail(addr, auth, config.From, msg.To, []byte(body))
}

// TestSMTPConnection 测试 SMTP 是否配置正确
func TestSMTPConnection(config *SMTPConfig) error {
	testMsg := &EmailMessage{
		To:      []string{config.From},
		Subject: "LingCDN SMTP 测试",
		Body:    "这是一封测试邮件，如果收到说明 SMTP 设置正确。",
		IsHTML:  false,
	}

	return SendEmail(config, testMsg)
}
