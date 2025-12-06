// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

// NewReferersConfig 鑾峰彇鏂伴槻鐩楅摼閰嶇疆瀵硅薄
func NewReferersConfig() *ReferersConfig {
	return &ReferersConfig{
		CheckOrigin: true,
	}
}

// ReferersConfig 闃茬洍閾捐缃?
type ReferersConfig struct {
	IsPrior         bool     `yaml:"isPrior" json:"isPrior"`
	IsOn            bool     `yaml:"isOn" json:"isOn"`
	AllowEmpty      bool     `yaml:"allowEmpty" json:"allowEmpty"`           // 鏉ユ簮鍩熷悕鍏佽涓虹┖
	AllowSameDomain bool     `yaml:"allowSameDomain" json:"allowSameDomain"` // 鍏佽鏉ユ簮鍩熷悕鍜屽綋鍓嶈闂殑鍩熷悕涓€鑷达紝鐩稿綋浜庡湪绔欏唴璁块棶
	AllowDomains    []string `yaml:"allowDomains" json:"allowDomains"`       // 鍏佽鐨勬潵婧愬煙鍚嶅垪琛?
	DenyDomains     []string `yaml:"denyDomains" json:"denyDomains"`         // 绂佹鐨勬潵婧愬煙鍚嶅垪琛?
	CheckOrigin     bool     `yaml:"checkOrigin" json:"checkOrigin"`         // 鏄惁妫€鏌rigin

	OnlyURLPatterns   []*shared.URLPattern `yaml:"onlyURLPatterns" json:"onlyURLPatterns"`     // 浠呴檺鐨刄RL
	ExceptURLPatterns []*shared.URLPattern `yaml:"exceptURLPatterns" json:"exceptURLPatterns"` // 鎺掗櫎鐨刄RL
}

func (this *ReferersConfig) Init() error {
	// url patterns
	for _, pattern := range this.ExceptURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	for _, pattern := range this.OnlyURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ReferersConfig) MatchDomain(requestDomain string, refererDomain string) bool {
	if len(refererDomain) == 0 {
		if this.AllowEmpty {
			return true
		}
		return false
	}

	if this.AllowSameDomain && requestDomain == refererDomain {
		return true
	}

	if len(this.AllowDomains) == 0 {
		if len(this.DenyDomains) > 0 {
			return !configutils.MatchDomains(this.DenyDomains, refererDomain)
		}
		return false
	}

	if configutils.MatchDomains(this.AllowDomains, refererDomain) {
		if len(this.DenyDomains) > 0 && configutils.MatchDomains(this.DenyDomains, refererDomain) {
			return false
		}
		return true
	}

	return false
}

func (this *ReferersConfig) MatchURL(url string) bool {
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
