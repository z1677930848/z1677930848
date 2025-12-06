package models

import (
	"fmt"
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
	db, err := dbs.Default()
	if err != nil {
		return nil, err
	}
	one, err := dbs.NewQuery(nil).
		DB(db).
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
	db, err := dbs.Default()
	if err != nil {
		return err
	}
	_, err = dbs.NewQuery(nil).
		DB(db).
		Table("email_settings").
		Where(fmt.Sprintf("id=%d", m.Id)).
		Sets(map[string]any{
			"smtp_host":     m.SmtpHost,
			"smtp_port":     m.SmtpPort,
			"smtp_username": m.SmtpUsername,
			"smtp_password": m.SmtpPassword,
			"from_email":    m.FromEmail,
			"from_name":     m.FromName,
			"use_tls":       m.UseTLS,
		}).
		Update()
	return err
}

func (m *EmailSettingModel) Create() error {
	db, err := dbs.Default()
	if err != nil {
		return err
	}
	_, err = dbs.NewQuery(nil).
		DB(db).
		Table("email_settings").
		Sets(map[string]any{
			"smtp_host":     m.SmtpHost,
			"smtp_port":     m.SmtpPort,
			"smtp_username": m.SmtpUsername,
			"smtp_password": m.SmtpPassword,
			"from_email":    m.FromEmail,
			"from_name":     m.FromName,
			"use_tls":       m.UseTLS,
			"is_enabled":    1,
		}).
		Insert()
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
