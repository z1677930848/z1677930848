// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package firewallconfigs

var DefaultHTTPFirewallPolicyLogConfig = &HTTPFirewallPolicyLogConfig{
	IsOn:          true,
	RequestBody:   true,
	RegionDenying: false,
}

type HTTPFirewallPolicyLogConfig struct {
	IsPrior       bool `yaml:"isPrior" json:"isPrior"`
	IsOn          bool `yaml:"isOn" json:"isOn"`
	RequestBody   bool `yaml:"requestBody" json:"requestBody"`     // 鏄惁璁板綍RequestBody
	RegionDenying bool `yaml:"regionDenying" json:"regionDenying"` // 鏄惁璁板綍鍖哄煙灏佺鏃ュ織
}

func (this *HTTPFirewallPolicyLogConfig) Init() error {
	return nil
}
