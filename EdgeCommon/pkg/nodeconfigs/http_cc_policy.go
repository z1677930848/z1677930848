// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"

// HTTPCCPolicy CC绛栫暐
type HTTPCCPolicy struct {
	IsOn       bool                             `json:"isOn" yaml:"isOn"`
	Thresholds []*serverconfigs.HTTPCCThreshold `json:"thresholds" yaml:"thresholds"` // 闃堝€?
}

func NewHTTPCCPolicy() *HTTPCCPolicy {
	return &HTTPCCPolicy{
		IsOn: true,
	}
}

func (this *HTTPCCPolicy) Init() error {
	return nil
}
