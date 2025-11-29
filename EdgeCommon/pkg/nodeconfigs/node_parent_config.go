// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package nodeconfigs

// ParentNodeConfig 鐖剁骇鑺傜偣閰嶇疆
type ParentNodeConfig struct {
	Id         int64    `yaml:"id" json:"id"`
	Addrs      []string `yaml:"addrs" json:"addrs"`
	SecretHash string   `yaml:"secretHash" json:"secretHash"`
}
