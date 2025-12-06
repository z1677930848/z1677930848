// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

import (
	"strings"

	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

type HTTPBaseOptimizationConfig struct {
	OnlyURLPatterns   []*shared.URLPattern `yaml:"onlyURLPatterns" json:"onlyURLPatterns"`     // 浠呴檺鐨刄RL
	ExceptURLPatterns []*shared.URLPattern `yaml:"exceptURLPatterns" json:"exceptURLPatterns"` // 鎺掗櫎鐨刄RL
}

func (this *HTTPBaseOptimizationConfig) Init() error {
	// only url
	for _, pattern := range this.OnlyURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	// except url
	for _, pattern := range this.ExceptURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *HTTPBaseOptimizationConfig) MatchURL(url string) bool {
	// 鍘婚櫎闂彿
	var index = strings.Index(url, "?")
	if index >= 0 {
		url = url[:index]
	}

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
