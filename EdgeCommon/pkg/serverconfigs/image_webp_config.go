// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"strings"

	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
)

// WebPImageConfig WebP閰嶇疆
type WebPImageConfig struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"`
	IsOn    bool `yaml:"isOn" json:"isOn"`

	MinLength  *shared.SizeCapacity           `yaml:"minLength" json:"minLength"`   // 鏈€灏忓帇缂╁璞℃瘮濡?m, 24k
	MaxLength  *shared.SizeCapacity           `yaml:"maxLength" json:"maxLength"`   // 鏈€澶у帇缂╁璞?
	MimeTypes  []string                       `yaml:"mimeTypes" json:"mimeTypes"`   // 鏀寔鐨凪imeType锛屾敮鎸乮mage/*杩欐牱鐨勯€氶厤绗︿娇鐢?
	Extensions []string                       `yaml:"extensions" json:"extensions"` // 鏂囦欢鎵╁睍鍚嶏紝鍖呭惈鐐圭鍙凤紝涓嶅尯鍒嗗ぇ灏忓啓
	Conds      *shared.HTTPRequestCondsConfig `yaml:"conds" json:"conds"`           // 鍖归厤鏉′欢

	minLength     int64
	maxLength     int64
	mimeTypeRules []*shared.MimeTypeRule
	extensions    []string
}

func (this *WebPImageConfig) Init() error {
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

	// mime types
	this.mimeTypeRules = []*shared.MimeTypeRule{}
	for _, mimeType := range this.MimeTypes {
		rule, err := shared.NewMimeTypeRule(mimeType)
		if err != nil {
			return err
		}
		this.mimeTypeRules = append(this.mimeTypeRules, rule)
	}

	// extensions
	this.extensions = []string{}
	for _, ext := range this.Extensions {
		ext = strings.ToLower(ext)
		if len(ext) > 0 && ext[0] != '.' {
			ext = "." + ext
		}
		this.extensions = append(this.extensions, ext)
	}

	return nil
}

// MatchResponse 鏄惁鍖归厤鍝嶅簲
func (this *WebPImageConfig) MatchResponse(mimeType string, contentLength int64, requestExt string, formatter shared.Formatter) bool {
	if this.Conds != nil && formatter != nil {
		if !this.Conds.MatchRequest(formatter) {
			return false
		}
		if !this.Conds.MatchResponse(formatter) {
			return false
		}
	}

	// no content
	if contentLength == 0 {
		return false
	}

	// min length
	if this.minLength > 0 && contentLength >= 0 && contentLength < this.minLength {
		return false
	}

	// max length
	if this.maxLength > 0 && contentLength > this.maxLength {
		return false
	}

	// extensions
	if len(this.extensions) > 0 {
		if len(requestExt) > 0 {
			for _, ext := range this.extensions {
				if ext == requestExt {
					if strings.Contains(mimeType, "image/") {
						return true
					}
				}
			}
		}
	}

	// mime types
	if len(this.mimeTypeRules) > 0 {
		if len(mimeType) > 0 {
			var index = strings.Index(mimeType, ";")
			if index >= 0 {
				mimeType = mimeType[:index]
			}
			for _, rule := range this.mimeTypeRules {
				if rule.Match(mimeType) {
					return true
				}
			}
		}
	}

	// 濡傛灉娌℃湁鎸囧畾鏉′欢锛屽垯鎵€鏈夌殑閮借兘鍘嬬缉
	if len(this.extensions) == 0 && len(this.mimeTypeRules) == 0 {
		return true
	}

	return false
}

// MatchRequest 鏄惁鍖归厤璇锋眰
func (this *WebPImageConfig) MatchRequest(requestExt string, formatter shared.Formatter) bool {
	if this.Conds != nil && formatter != nil {
		if !this.Conds.MatchRequest(formatter) {
			return false
		}
	}

	// extensions
	if len(this.mimeTypeRules) == 0 && len(this.extensions) > 0 && len(requestExt) > 0 {
		for _, ext := range this.extensions {
			if ext == requestExt {
				return true
			}
		}
		return false
	}

	return true
}

// MatchAccept 妫€鏌ュ鎴风鏄惁鑳芥帴鍙梂ebP
func (this *WebPImageConfig) MatchAccept(acceptContentTypes string) bool {
	var t = "image/webp"
	var index = strings.Index(acceptContentTypes, t)
	if index < 0 {
		return false
	}
	var l = len(acceptContentTypes)
	if index > 0 && acceptContentTypes[index-1] != ',' {
		return false
	}

	if index+len(t) < l && acceptContentTypes[index+len(t)] != ',' {
		return false
	}
	return true
}
