// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package firewallconfigs

type HTTPFirewallJavascriptCookieAction struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"`

	Life              int32  `yaml:"life" json:"life"`                         // 鏈夋晥鏈?
	MaxFails          int    `yaml:"maxFails" json:"maxFails"`                 // 鏈€澶уけ璐ユ鏁?
	FailBlockTimeout  int    `yaml:"failBlockTimeout" json:"failBlockTimeout"` // 澶辫触鎷︽埅鏃堕棿
	Scope             string `yaml:"scope" json:"scope"`
	FailBlockScopeAll bool   `yaml:"failBlockScopeAll" json:"failBlockScopeAll"`
}

func NewHTTPFirewallJavascriptCookieAction() *HTTPFirewallJavascriptCookieAction {
	return &HTTPFirewallJavascriptCookieAction{
		Life:              600,
		MaxFails:          100,
		FailBlockTimeout:  3600,
		Scope:             FirewallScopeServer,
		FailBlockScopeAll: true,
	}
}
