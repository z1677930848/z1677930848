// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

// UAMConfig UAM閰嶇疆
type UAMConfig struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"`
	IsOn    bool `yaml:"isOn" json:"isOn"`

	AddToWhiteList    bool                           `yaml:"addToWhiteList" json:"addToWhiteList"`       // 鏄惁灏咺P鍔犲叆鍒扮櫧鍚嶅崟
	OnlyURLPatterns   []*shared.URLPattern           `yaml:"onlyURLPatterns" json:"onlyURLPatterns"`     // 浠呴檺鐨刄RL
	ExceptURLPatterns []*shared.URLPattern           `yaml:"exceptURLPatterns" json:"exceptURLPatterns"` // 鎺掗櫎鐨刄RL
	MinQPSPerIP       int                            `yaml:"minQPSPerIP" json:"minQPSPerIP"`             // 鍚敤瑕佹眰鐨勫崟IP鏈€浣庡钩鍧嘠PS
	Conds             *shared.HTTPRequestCondsConfig `yaml:"conds" json:"conds"`                         // 鍖归厤鏉′欢
	KeyLife           int                            `yaml:"keyLife" json:"keyLife"`                     // Key鏈夋晥鏈?
}

func NewUAMConfig() *UAMConfig {
	return &UAMConfig{
		AddToWhiteList: true,
	}
}

func (this *UAMConfig) Init() error {
	// only urls
	for _, pattern := range this.OnlyURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	// except urls
	for _, pattern := range this.ExceptURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	// conds
	if this.Conds != nil {
		err := this.Conds.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *UAMConfig) MatchURL(url string) bool {
	// except
	if len(this.ExceptURLPatterns) > 0 {
		for _, pattern := range this.ExceptURLPatterns {
			if pattern.Match(url) {
				return false
			}
		}
	}

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

func (this *UAMConfig) MatchRequest(formatter func(s string) string) bool {
	if this.Conds == nil {
		return true
	}
	return this.Conds.MatchRequest(formatter)
}
