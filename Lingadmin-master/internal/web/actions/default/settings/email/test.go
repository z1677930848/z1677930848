package email

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) RunPost(params struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	FromEmail    string
	UseTLS       bool
	TestEmail    string

	Must *actions.Must
}) {
	params.Must.
		Field("testEmail", params.TestEmail).
		Require("").
		Match(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "")

	config := &utils.SMTPConfig{
		Host:     params.SmtpHost,
		Port:     params.SmtpPort,
		Username: params.SmtpUsername,
		Password: params.SmtpPassword,
		From:     params.FromEmail,
		UseTLS:   params.UseTLS,
	}

	msg := &utils.EmailMessage{
		To:      []string{params.TestEmail},
		Subject: "LingCDN ",
		Body:    "<h2></h2><p>SMTP</p><p>LingCDN </p>",
		IsHTML:  true,
	}

	err := utils.SendEmail(config, msg)
	if err != nil {
		this.Fail(": " + err.Error())
	}

	this.Success()
}
