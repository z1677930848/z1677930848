// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

const DefaultOpenFileCacheMax = 1024

// OpenFileCacheConfig open file cache閰嶇疆
type OpenFileCacheConfig struct {
	IsOn bool `yaml:"isOn" json:"isOn"`
	Max  int  `yaml:"max" json:"max"`
}

func (this *OpenFileCacheConfig) Init() error {
	return nil
}
