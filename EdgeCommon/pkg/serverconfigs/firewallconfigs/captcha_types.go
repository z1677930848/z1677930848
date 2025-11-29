// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package firewallconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

type CaptchaType = string

const (
	CaptchaTypeDefault  CaptchaType = "default"
	CaptchaTypeOneClick CaptchaType = "oneClick"
	CaptchaTypeSlide    CaptchaType = "slide"
	CaptchaTypeGeeTest  CaptchaType = "geetest"
)

// FindAllCaptchaTypes returns supported captcha implementations
func FindAllCaptchaTypes() []*shared.Definition {
	return []*shared.Definition{
		{Code: CaptchaTypeDefault, Name: "Text Captcha", Description: "Input text to verify"},
		{Code: CaptchaTypeOneClick, Name: "One Click", Description: "Click button to verify"},
		{Code: CaptchaTypeSlide, Name: "Slide", Description: "Slide block to verify"},
		{Code: CaptchaTypeGeeTest, Name: "Geetest", Description: "Use Geetest service"},
	}
}

func DefaultCaptchaType() *shared.Definition {
	types := FindAllCaptchaTypes()
	if len(types) > 0 {
		return types[0]
	}
	return &shared.Definition{Code: CaptchaTypeDefault, Name: "Text Captcha"}
}

func FindCaptchaType(code CaptchaType) *shared.Definition {
	if len(code) == 0 {
		code = CaptchaTypeDefault
	}
	for _, t := range FindAllCaptchaTypes() {
		if t.Code == code {
			return t
		}
	}
	return DefaultCaptchaType()
}
