// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

// NewTOAConfig 榛樿鐨凾OA閰嶇疆
func NewTOAConfig() *TOAConfig {
	return &TOAConfig{}
}

// TOAConfig TOA鐩稿叧閰嶇疆
type TOAConfig struct {
	IsOn bool `yaml:"isOn" json:"isOn"`
}

func (this *TOAConfig) Init() error {
	return nil
}

func (this *TOAConfig) RandLocalPort() uint16 {
	return 0
}
