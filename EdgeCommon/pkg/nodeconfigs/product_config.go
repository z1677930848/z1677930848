// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package nodeconfigs

const DefaultProductName = "Lingcdn"

// ProductConfig 浜у搧鐩稿叧璁剧疆
type ProductConfig struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}
