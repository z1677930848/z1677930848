// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package userconfigs

const (
	EmailVerificationDefaultLife   = 86400 * 2
	EmailResetPasswordDefaultLife  = 3600
	MobileVerificationDefaultLife  = 1800
	MobileResetPasswordDefaultLife = 1800
)

type UserRegisterConfig struct {
	IsOn                bool `yaml:"isOn" json:"isOn"`
	ComplexPassword     bool `yaml:"complexPassword" json:"complexPassword"`
	RequireVerification bool `yaml:"requireVerification" json:"requireVerification"`
	RequireIdentity     bool `yaml:"requireIdentity" json:"requireIdentity"`
	CheckClientRegion   bool `yaml:"checkClientRegion" json:"checkClientRegion"`

	EmailVerification struct {
		IsOn       bool   `yaml:"isOn" json:"isOn"`
		ShowNotice bool   `yaml:"showNotice" json:"showNotice"`
		Subject    string `yaml:"subject" json:"subject"`
		Body       string `yaml:"body" json:"body"`
		CanLogin   bool   `yaml:"canLogin" json:"canLogin"`
		Life       int32  `yaml:"life" json:"life"`
	} `yaml:"emailVerification" json:"emailVerification"`

	EmailResetPassword struct {
		IsOn    bool   `yaml:"isOn" json:"isOn"`
		Subject string `yaml:"subject" json:"subject"`
		Body    string `yaml:"body" json:"body"`
		Life    int32  `yaml:"life" json:"life"`
	} `yaml:"emailResetPassword" json:"emailResetPassword"`

	MobileVerification struct {
		IsOn       bool   `yaml:"isOn" json:"isOn"`
		ShowNotice bool   `yaml:"showNotice" json:"showNotice"`
		CanLogin   bool   `yaml:"canLogin" json:"canLogin"`
		Body       string `yaml:"body" json:"body"`
		Life       int32  `yaml:"life" json:"life"`
		Force      bool   `yaml:"force" json:"force"`
	} `yaml:"mobileVerification" json:"mobileVerification"`

	CDNIsOn   bool     `json:"cdnIsOn"`
	ClusterId int64    `yaml:"clusterId" json:"clusterId"`
	Features  []string `yaml:"features" json:"features"`

	NSIsOn bool `json:"nsIsOn"`
	ADIsOn bool `json:"adIsOn"`
}

func DefaultUserRegisterConfig() *UserRegisterConfig {
	cfg := &UserRegisterConfig{
		IsOn:                false,
		ComplexPassword:     true,
		CDNIsOn:             true,
		NSIsOn:              false,
		RequireVerification: false,
		Features: []string{
			UserFeatureCodeServerAccessLog,
			UserFeatureCodeServerViewAccessLog,
			UserFeatureCodeServerWAF,
			UserFeatureCodePlan,
		},
	}
	cfg.EmailVerification.CanLogin = true
	cfg.EmailVerification.ShowNotice = true
	cfg.EmailVerification.Subject = "[${product.name}] Email verification"
	cfg.EmailVerification.Body = `<p>Welcome to ${product.name}. Please visit the following link to verify your email.</p>
<p><a href="${url.verify}" target="_blank">${url.verify}</a></p>
<p>If the link is not clickable, copy and paste it into your browser.</p>
<p>${product.name} Team</p>
<p><a href="${url.home}" target="_blank">${url.home}</a></p>`

	cfg.EmailResetPassword.IsOn = true
	cfg.EmailResetPassword.Subject = "[${product.name}] Reset password"
	cfg.EmailResetPassword.Body = `<p>You requested to reset password. Input the following code:</p>
<p><strong>${code}</strong></p>
<p>${product.name} Team</p>
<p><a href="${url.home}" target="_blank">${url.home}</a></p>`

	cfg.MobileVerification.Body = "Your verification code is ${code}"

	return cfg
}
