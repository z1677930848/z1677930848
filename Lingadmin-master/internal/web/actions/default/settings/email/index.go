package email

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/models"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "settings", "email")
}

func (this *IndexAction) RunGet(params struct{}) {
	model := &models.EmailSettingModel{}
	setting, err := model.Get()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if setting != nil {
		this.Data["smtpHost"] = setting.SmtpHost
		this.Data["smtpPort"] = setting.SmtpPort
		this.Data["smtpUsername"] = setting.SmtpUsername
		this.Data["fromEmail"] = setting.FromEmail
		this.Data["fromName"] = setting.FromName
		this.Data["useTLS"] = setting.UseTLS
	} else {
		this.Data["smtpHost"] = ""
		this.Data["smtpPort"] = 587
		this.Data["smtpUsername"] = ""
		this.Data["fromEmail"] = ""
		this.Data["fromName"] = "LingCDN"
		this.Data["useTLS"] = true
	}
	this.Data["testEmail"] = ""

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	FromEmail    string
	FromName     string
	UseTLS       bool

	Must *actions.Must
}) {
	params.Must.
		Field("smtpHost", params.SmtpHost).
		Require("SMTP").
		Field("smtpPort", params.SmtpPort).
		Require("SMTP").
		Gt(0, "0").
		Field("smtpUsername", params.SmtpUsername).
		Require("SMTP").
		Field("smtpPassword", params.SmtpPassword).
		Require("SMTP").
		Field("fromEmail", params.FromEmail).
		Require("").
		Match(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "")

	model := &models.EmailSettingModel{}
	existing, err := model.Get()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if existing != nil {
		existing.SmtpHost = params.SmtpHost
		existing.SmtpPort = params.SmtpPort
		existing.SmtpUsername = params.SmtpUsername
		existing.SmtpPassword = params.SmtpPassword
		existing.FromEmail = params.FromEmail
		existing.FromName = params.FromName
		existing.UseTLS = params.UseTLS
		err = existing.Save()
	} else {
		model.SmtpHost = params.SmtpHost
		model.SmtpPort = params.SmtpPort
		model.SmtpUsername = params.SmtpUsername
		model.SmtpPassword = params.SmtpPassword
		model.FromEmail = params.FromEmail
		model.FromName = params.FromName
		model.UseTLS = params.UseTLS
		err = model.Create()
	}

	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
