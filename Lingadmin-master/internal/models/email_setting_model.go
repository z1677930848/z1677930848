package models

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo/dbs"
)

type EmailSettingModel struct {
	Id           uint32 `field:"id"`
	SmtpHost     string `field:"smtp_host"`
	SmtpPort     int    `field:"smtp_port"`
	SmtpUsername string `field:"smtp_username"`
	SmtpPassword string `field:"smtp_password"`
	FromEmail    string `field:"from_email"`
	FromName     string `field:"from_name"`
	UseTLS       bool   `field:"use_tls"`
	IsEnabled    bool   `field:"is_enabled"`
}

func (m *EmailSettingModel) Get() (*EmailSettingModel, error) {
	one, err := dbs.Default().
		Table("email_settings").
		Where("is_enabled=1").
		Find()
	if err != nil {
		return nil, err
	}
	if one == nil {
		return nil, nil
	}
	return one.(*EmailSettingModel), nil
}

func (m *EmailSettingModel) Save() error {
	_, err := dbs.Default().
		Table("email_settings").
		Where("id=?", m.Id).
		Update(map[string]interface{}{
			"smtp_host":     m.SmtpHost,
			"smtp_port":     m.SmtpPort,
			"smtp_username": m.SmtpUsername,
			"smtp_password": m.SmtpPassword,
			"from_email":    m.FromEmail,
			"from_name":     m.FromName,
			"use_tls":       m.UseTLS,
		})
	return err
}

func (m *EmailSettingModel) Create() error {
	_, err := dbs.Default().
		Table("email_settings").
		Insert(map[string]interface{}{
			"smtp_host":     m.SmtpHost,
			"smtp_port":     m.SmtpPort,
			"smtp_username": m.SmtpUsername,
			"smtp_password": m.SmtpPassword,
			"from_email":    m.FromEmail,
			"from_name":     m.FromName,
			"use_tls":       m.UseTLS,
			"is_enabled":    1,
		})
	return err
}

func (m *EmailSettingModel) GetSMTPConfig() *utils.SMTPConfig {
	return &utils.SMTPConfig{
		Host:     m.SmtpHost,
		Port:     m.SmtpPort,
		Username: m.SmtpUsername,
		Password: m.SmtpPassword,
		From:     m.FromEmail,
		UseTLS:   m.UseTLS,
	}
}
