// Copyright 2024 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

// HLSEncryptingConfig HLS鍔犲瘑閰嶇疆
type HLSEncryptingConfig struct {
	IsOn bool `yaml:"isOn" json:"isOn"` // 鏄惁鍚敤

	OnlyURLPatterns   []*shared.URLPattern `yaml:"onlyURLPatterns" json:"onlyURLPatterns"`     // 浠呴檺鐨刄RL
	ExceptURLPatterns []*shared.URLPattern `yaml:"exceptURLPatterns" json:"exceptURLPatterns"` // 鎺掗櫎鐨刄RL
}

func (this *HLSEncryptingConfig) Init() error {
	for _, pattern := range this.OnlyURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	for _, pattern := range this.ExceptURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *HLSEncryptingConfig) MatchURL(url string) bool {
	// except
	if len(this.ExceptURLPatterns) > 0 {
		for _, pattern := range this.ExceptURLPatterns {
			if pattern.Match(url) {
				return false
			}
		}
	}

	// only
	if len(this.OnlyURLPatterns) > 0 {
		for _, pattern := range this.OnlyURLPatterns {
			if pattern.Match(url) {
				return true
			}
		}
		return false
	}

	return true
}
