// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// AccessLogTCPStorageConfig TCP瀛樺偍绛栫暐
type AccessLogTCPStorageConfig struct {
	Network string `yaml:"network" json:"network"` // tcp, unix
	Addr    string `yaml:"addr" json:"addr"`
}
