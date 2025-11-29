// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package nodeconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

func init() {
	_ = DefaultWebPImagePolicy.Init()
}

var DefaultWebPImagePolicy = NewWebPImagePolicy()

func NewWebPImagePolicy() *WebPImagePolicy {
	return &WebPImagePolicy{
		IsOn:         true,
		RequireCache: true,
		MinLength:    shared.NewSizeCapacity(0, shared.SizeCapacityUnitKB),
		MaxLength:    shared.NewSizeCapacity(128, shared.SizeCapacityUnitMB),
	}
}

// WebPImagePolicy WebP绛栫暐
type WebPImagePolicy struct {
	IsOn         bool                 `yaml:"isOn" json:"isOn"`                 // 鏄惁鍚敤
	RequireCache bool                 `yaml:"requireCache" json:"requireCache"` // 闇€瑕佸湪缂撳瓨鏉′欢涓嬭繘琛?
	MinLength    *shared.SizeCapacity `yaml:"minLength" json:"minLength"`       // 鏈€灏忓帇缂╁璞℃瘮濡?m, 24k
	MaxLength    *shared.SizeCapacity `yaml:"maxLength" json:"maxLength"`       // 鏈€澶у帇缂╁璞?
	Quality      int                  `yaml:"quality" json:"quality"`           // 鐢熸垚鐨勫浘鐗囪川閲忥細0-100

	minLength int64
	maxLength int64
}

func (this *WebPImagePolicy) Init() error {
	if this.MinLength != nil {
		this.minLength = this.MinLength.Bytes()
	}
	if this.MaxLength != nil {
		this.maxLength = this.MaxLength.Bytes()
	}

	return nil
}

func (this *WebPImagePolicy) MinLengthBytes() int64 {
	return this.minLength
}

func (this *WebPImagePolicy) MaxLengthBytes() int64 {
	return this.maxLength
}
