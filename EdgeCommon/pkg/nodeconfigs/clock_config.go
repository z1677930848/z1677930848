// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

func DefaultClockConfig() *ClockConfig {
	return &ClockConfig{
		AutoSync:    true,
		Server:      "",
		CheckChrony: true,
	}
}

// ClockConfig 鏃堕挓鐩稿叧閰嶇疆
type ClockConfig struct {
	AutoSync    bool   `yaml:"autoSync" json:"autoSync"`       // 鑷姩灏濊瘯鍚屾鏃堕挓
	Server      string `yaml:"server" json:"server"`           // 鏃堕挓鍚屾鏈嶅姟鍣?
	CheckChrony bool   `yaml:"checkChrony" json:"checkChrony"` // 妫€鏌?chronyd 鏄惁鍦ㄨ繍琛?
}

func (this *ClockConfig) Init() error {
	return nil
}
