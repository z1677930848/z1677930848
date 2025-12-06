// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

// FollowProtocolConfig 鍗忚璺熼殢閰嶇疆
type FollowProtocolConfig struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"` // 鏄惁瑕嗙洊鐖剁骇閰嶇疆
	IsOn    bool `yaml:"isOn" json:"isOn"`       // 鏄惁鍚敤
	HTTP    struct {
		Port       int  `yaml:"port" json:"port"`             // 绔彛
		FollowPort bool `yaml:"followPort" json:"followPort"` // 璺熼殢绔彛
	} `yaml:"http" json:"http"` // HTTP閰嶇疆
	HTTPS struct {
		Port       int  `yaml:"port" json:"port"`             // 绔彛
		FollowPort bool `yaml:"followPort" json:"followPort"` // 璺熼殢绔彛
	} `yaml:"https" json:"https"` // HTTPS閰嶇疆
}

func NewFollowProtocolConfig() *FollowProtocolConfig {
	var p = &FollowProtocolConfig{}
	p.HTTP.FollowPort = true
	p.HTTPS.FollowPort = true
	return p
}

func (this *FollowProtocolConfig) Init() error {
	return nil
}
