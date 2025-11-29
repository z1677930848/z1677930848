package utils

import (
	"fmt"
	"strings"
)

type EmailService struct {
	config *SMTPConfig
}

func NewEmailService(config *SMTPConfig) *EmailService {
	return &EmailService{config: config}
}

func (s *EmailService) SendTemplateEmail(to []string, template string, vars map[string]string) error {
	subject, body := s.renderTemplate(template, vars)
	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}
	return SendEmail(s.config, msg)
}

func (s *EmailService) renderTemplate(template string, vars map[string]string) (string, string) {
	templates := map[string]struct {
		subject string
		body    string
	}{
		"license_expire": {
			"LingCDN ",
			"<h2></h2><p> {{customer_name}}</p><p> <strong>{{license_code}}</strong>  <strong>{{expire_date}}</strong> </p><p></p><p>LingCDN </p>",
		},
		"system_alert": {
			"LingCDN ",
			"<h2></h2><p>{{alert_time}}</p><p>{{alert_level}}</p><p>{{alert_message}}</p><p></p>",
		},
		"welcome": {
			" LingCDN",
			"<h2> LingCDN</h2><p> {{username}}</p><p></p><p>{{login_url}}</p><p>{{username}}</p><p></p><p>LingCDN </p>",
		},
	}

	tpl, ok := templates[template]
	if !ok {
		return "", ""
	}

	subject := tpl.subject
	body := tpl.body
	for k, v := range vars {
		body = strings.ReplaceAll(body, fmt.Sprintf("{{%s}}", k), v)
	}

	return subject, body
}
