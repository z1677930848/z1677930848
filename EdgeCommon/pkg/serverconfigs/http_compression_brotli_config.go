// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

type HTTPBrotliCompressionConfig struct {
	Id   int64 `yaml:"id" json:"id"` // ID
	IsOn bool  `yaml:"isOn" json:"isOn"`

	Level     int8                           `yaml:"level" json:"level"`         // 绾у埆
	MinLength *shared.SizeCapacity           `yaml:"minLength" json:"minLength"` // 鏈€灏忓帇缂╁璞℃瘮濡?m, 24k
	MaxLength *shared.SizeCapacity           `yaml:"maxLength" json:"maxLength"` // 鏈€澶у帇缂╁璞?
	Conds     *shared.HTTPRequestCondsConfig `yaml:"conds" json:"conds"`         // 鍖归厤鏉′欢

	minLength int64
	maxLength int64
}

func (this *HTTPBrotliCompressionConfig) Init() error {
	if this.MinLength != nil {
		this.minLength = this.MinLength.Bytes()
	}
	if this.MaxLength != nil {
		this.maxLength = this.MaxLength.Bytes()
	}

	if this.Conds != nil {
		err := this.Conds.Init()
		if err != nil {
			return err
		}
	}

	return nil
}
