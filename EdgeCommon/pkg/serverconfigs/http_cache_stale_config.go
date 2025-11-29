// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

// HTTPCacheStaleConfig Stale绛栫暐閰嶇疆
type HTTPCacheStaleConfig struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"`
	IsOn    bool `yaml:"isOn" json:"isOn"` // 鏄惁鍚敤

	Status                    []int                `yaml:"status" json:"status"`                                       // 鐘舵€佸垪琛?
	SupportStaleIfErrorHeader bool                 `yaml:"supportStaleIfErrorHeader" json:"supportStaleIfErrorHeader"` // 鏄惁鏀寔stale-if-error
	Life                      *shared.TimeDuration `yaml:"life" json:"life"`                                           // 闄堟棫鍐呭鐢熷懡鍛ㄦ湡
}
