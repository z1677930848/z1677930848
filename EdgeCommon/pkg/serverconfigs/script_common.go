// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// CommonScript 鍏叡鑴氭湰
type CommonScript struct {
	Id       int64  `yaml:"id" json:"id"`
	IsOn     bool   `yaml:"isOn" json:"isOn"`
	Filename string `yaml:"filename" json:"filename"`
	Code     string `yaml:"code" json:"code"`
}
