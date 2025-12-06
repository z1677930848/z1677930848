// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package ddosconfigs

type TCPConfig struct {
	IsPrior             bool  `json:"isPrior"`
	IsOn                bool  `json:"isOn"`
	MaxConnections      int32 `json:"maxConnections"`
	MaxConnectionsPerIP int32 `json:"maxConnectionsPerIP"`

	// 鍒嗛挓绾ч€熺巼
	NewConnectionsMinutelyRate             int32 `json:"newConnectionsRate"`             // 鍒嗛挓
	NewConnectionsMinutelyRateBlockTimeout int32 `json:"newConnectionsRateBlockTimeout"` // 鎷︽埅鏃堕棿

	// 绉掔骇閫熺巼
	NewConnectionsSecondlyRate             int32 `json:"newConnectionsSecondlyRate"`
	NewConnectionsSecondlyRateBlockTimeout int32 `json:"newConnectionsSecondlyRateBlockTimeout"`

	AllowIPList []*IPConfig   `json:"allowIPList"`
	Ports       []*PortConfig `json:"ports"`
}

func (this *TCPConfig) Init() error {
	return nil
}
