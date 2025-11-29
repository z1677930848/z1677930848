// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package firewallconfigs

// SYNFloodConfig Syn flood闃叉姢璁剧疆
type SYNFloodConfig struct {
	IsPrior        bool  `yaml:"isPrior" json:"isPrior"`
	IsOn           bool  `yaml:"isOn" json:"isOn"`
	MinAttempts    int32 `yaml:"minAttempts" json:"minAttempts"`       // 鏈€灏忓皾璇曟鏁?鍒嗛挓
	TimeoutSeconds int32 `yaml:"timeoutSeconds" json:"timeoutSeconds"` // 鎷︽埅瓒呮椂鏃堕棿
	IgnoreLocal    bool  `yaml:"ignoreLocal" json:"ignoreLocal"`       // 蹇界暐鏈湴IP
}

func NewSYNFloodConfig() *SYNFloodConfig {
	return &SYNFloodConfig{
		IsOn:           false,
		MinAttempts:    10,
		TimeoutSeconds: 1800,
		IgnoreLocal:    true,
	}
}

func (this *SYNFloodConfig) Init() error {
	return nil
}
