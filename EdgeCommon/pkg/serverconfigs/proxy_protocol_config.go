// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"errors"

	"github.com/iwind/TeaGo/types"
)

type ProxyProtocolVersion = int

const (
	ProxyProtocolVersion1 ProxyProtocolVersion = 1
	ProxyProtocolVersion2 ProxyProtocolVersion = 2
)

// ProxyProtocolConfig PROXY Protocol閰嶇疆
type ProxyProtocolConfig struct {
	IsOn    bool                 `yaml:"isOn" json:"isOn"`
	Version ProxyProtocolVersion `yaml:"version" json:"version"`
}

// Init 鍒濆鍖?
func (this *ProxyProtocolConfig) Init() error {
	if this.IsOn {
		if this.Version != ProxyProtocolVersion1 && this.Version != ProxyProtocolVersion2 {
			return errors.New("invalid ProxyProtocol version '" + types.String(this.Version) + "'")
		}
	}

	return nil
}
