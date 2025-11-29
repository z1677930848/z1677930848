// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

// HTTPRequestLimitConfig HTTP璇锋眰闄愬埗鐩稿叧闄愬埗閰嶇疆
type HTTPRequestLimitConfig struct {
	IsPrior             bool                 `yaml:"isPrior" json:"isPrior"`                         // 鏄惁瑕嗙洊鐖剁骇
	IsOn                bool                 `yaml:"isOn" json:"isOn"`                               // 鏄惁鍚敤
	MaxConns            int                  `yaml:"maxConns" json:"maxConns"`                       // 骞跺彂杩炴帴鏁?
	MaxConnsPerIP       int                  `yaml:"maxConnsPerIP" json:"maxConnsPerIP"`             // 鍗曚釜IP骞跺彂杩炴帴鏁?
	OutBandwidthPerConn *shared.SizeCapacity `yaml:"outBandwidthPerConn" json:"outBandwidthPerConn"` // 涓嬭娴侀噺闄愬埗
	MaxBodySize         *shared.SizeCapacity `yaml:"maxBodySize" json:"maxBodySize"`                 // 鍗曚釜璇锋眰鏈€澶у昂瀵?

	outBandwidthPerConnBytes int64
	maxBodyBytes             int64
}

func (this *HTTPRequestLimitConfig) Init() error {
	if this.OutBandwidthPerConn != nil {
		this.outBandwidthPerConnBytes = this.OutBandwidthPerConn.Bytes()
	}
	if this.MaxBodySize != nil {
		this.maxBodyBytes = this.MaxBodySize.Bytes()
	}

	return nil
}

func (this *HTTPRequestLimitConfig) OutBandwidthPerConnBytes() int64 {
	return this.outBandwidthPerConnBytes
}

func (this *HTTPRequestLimitConfig) MaxBodyBytes() int64 {
	return this.maxBodyBytes
}
