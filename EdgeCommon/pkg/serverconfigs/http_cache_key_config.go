// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

// HTTPCacheKeyConfig 缂撳瓨Key閰嶇疆
type HTTPCacheKeyConfig struct {
	IsOn   bool   `yaml:"isOn" json:"isOn"`
	Scheme string `yaml:"scheme" json:"scheme"`
	Host   string `yaml:"host" json:"host"`
}

func (this *HTTPCacheKeyConfig) Init() error {
	return nil
}
