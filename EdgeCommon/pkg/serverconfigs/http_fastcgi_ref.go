// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

type HTTPFastcgiRef struct {
	IsPrior    bool    `yaml:"isPrior" json:"isPrior"`       // 鏄惁瑕嗙洊
	IsOn       bool    `yaml:"isOn" json:"isOn"`             // 鏄惁寮€鍚?
	FastcgiIds []int64 `yaml:"fastcgiIds" json:"fastcgiIds"` // Fastcgi ID鍒楄〃
}
