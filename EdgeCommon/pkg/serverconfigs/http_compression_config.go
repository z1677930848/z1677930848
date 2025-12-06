// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"strings"

	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/lists"
)

var DefaultHTTPCompressionTypes = []HTTPCompressionType{HTTPCompressionTypeBrotli, HTTPCompressionTypeGzip, HTTPCompressionTypeDeflate}

type HTTPCompressionRef struct {
	Id   int64 `yaml:"id" json:"id"`
	IsOn bool  `yaml:"isOn" json:"isOn"`
}

// HTTPCompressionConfig 鍐呭鍘嬬缉閰嶇疆
type HTTPCompressionConfig struct {
	IsPrior bool `yaml:"isPrior" json:"isPrior"`
	IsOn    bool `yaml:"isOn" json:"isOn"`

	UseDefaultTypes bool                  `yaml:"useDefaultTypes" json:"useDefaultTypes"` // 鏄惁浣跨敤榛樿鐨勭被鍨?
	Types           []HTTPCompressionType `yaml:"types" json:"types"`                     // 鏀寔鐨勭被鍨嬶紝濡傛灉涓虹┖琛ㄧず榛樿椤哄簭
	Level           int8                  `yaml:"level" json:"level"`                     // 绾у埆锛?-12
	DecompressData  bool                  `yaml:"decompressData" json:"decompressData"`   // 鏄惁瑙ｅ帇宸插帇缂╁唴瀹?

	GzipRef    *HTTPCompressionRef `yaml:"gzipRef" json:"gzipRef"`
	DeflateRef *HTTPCompressionRef `yaml:"deflateRef" json:"deflateRef"`
	BrotliRef  *HTTPCompressionRef `yaml:"brotliRef" json:"brotliRef"`

	Gzip    *HTTPGzipCompressionConfig    `yaml:"gzip" json:"gzip"`
	Deflate *HTTPDeflateCompressionConfig `yaml:"deflate" json:"deflate"`
	Brotli  *HTTPBrotliCompressionConfig  `yaml:"brotli" json:"brotli"`

	MinLength            *shared.SizeCapacity           `yaml:"minLength" json:"minLength"`                       // 鏈€灏忓帇缂╁璞℃瘮濡?m, 24k
	MaxLength            *shared.SizeCapacity           `yaml:"maxLength" json:"maxLength"`                       // 鏈€澶у帇缂╁璞?
	MimeTypes            []string                       `yaml:"mimeTypes" json:"mimeTypes"`                       // 鏀寔鐨凪imeType锛屾敮鎸乮mage/*杩欐牱鐨勯€氶厤绗︿娇鐢?
	Extensions           []string                       `yaml:"extensions" json:"extensions"`                     // 鏂囦欢鎵╁睍鍚嶏紝鍖呭惈鐐圭鍙凤紝涓嶅尯鍒嗗ぇ灏忓啓
	ExceptExtensions     []string                       `yaml:"exceptExtensions" json:"exceptExtensions"`         // 渚嬪鎵╁睍鍚?
	Conds                *shared.HTTPRequestCondsConfig `yaml:"conds" json:"conds"`                               // 鍖归厤鏉′欢
	EnablePartialContent bool                           `yaml:"enablePartialContent" json:"enablePartialContent"` // 鏀寔PartialContent鍘嬬缉

	OnlyURLPatterns   []*shared.URLPattern `yaml:"onlyURLPatterns" json:"onlyURLPatterns"`     // 浠呴檺鐨刄RL
	ExceptURLPatterns []*shared.URLPattern `yaml:"exceptURLPatterns" json:"exceptURLPatterns"` // 鎺掗櫎鐨刄RL

	minLength        int64
	maxLength        int64
	mimeTypeRules    []*shared.MimeTypeRule
	extensions       []string
	exceptExtensions []string

	types []HTTPCompressionType

	supportGzip    bool
	supportDeflate bool
	supportBrotli  bool
	supportZSTD    bool
}

// Init 鍒濆鍖?
func (this *HTTPCompressionConfig) Init() error {
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

	this.exceptExtensions = []string{}
	for _, ext := range this.ExceptExtensions {
		ext = strings.ToLower(ext)
		if len(ext) > 0 && ext[0] != '.' {
			ext = "." + ext
		}
		this.exceptExtensions = append(this.exceptExtensions, ext)
	}

	if this.Gzip != nil {
		err := this.Gzip.Init()
		if err != nil {
			return err
		}
	}

	if this.Deflate != nil {
		err := this.Deflate.Init()
		if err != nil {
			return err
		}
	}

	if this.Brotli != nil {
		err := this.Brotli.Init()
		if err != nil {
			return err
		}
	}

	var supportedTypes = []HTTPCompressionType{}
	if !this.UseDefaultTypes {
		supportedTypes = append(supportedTypes, this.Types...)
	} else {
		supportedTypes = append(supportedTypes, DefaultHTTPCompressionTypes...)
	}
	this.types = supportedTypes

	this.supportGzip = false
	this.supportDeflate = false
	this.supportDeflate = false
	for _, supportType := range supportedTypes {
		switch supportType {
		case HTTPCompressionTypeGzip:
			if this.GzipRef == nil || (this.GzipRef != nil && this.GzipRef.IsOn && this.Gzip != nil && this.Gzip.IsOn) {
				this.supportGzip = true
			}
		case HTTPCompressionTypeDeflate:
			if this.DeflateRef == nil || (this.DeflateRef != nil && this.DeflateRef.IsOn && this.Deflate != nil && this.Deflate.IsOn) {
				this.supportDeflate = true
			}
		case HTTPCompressionTypeBrotli:
			if this.BrotliRef == nil || (this.BrotliRef != nil && this.BrotliRef.IsOn && this.Brotli != nil && this.Brotli.IsOn) {
				this.supportBrotli = true
			}
		case HTTPCompressionTypeZSTD:
			this.supportZSTD = true
		}
	}

	// url patterns
	for _, pattern := range this.ExceptURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	for _, pattern := range this.OnlyURLPatterns {
		err := pattern.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

// MinBytes 鍙帇缂╂渶灏忓昂瀵?
func (this *HTTPCompressionConfig) MinBytes() int64 {
	return this.minLength
}

// MaxBytes 鍙帇缂╂渶澶у昂瀵?
func (this *HTTPCompressionConfig) MaxBytes() int64 {
	return this.maxLength
}

// MatchResponse 鏄惁鍖归厤鍝嶅簲
func (this *HTTPCompressionConfig) MatchResponse(mimeType string, contentLength int64, requestExt string, formatter shared.Formatter) bool {
	if this.Conds != nil && formatter != nil {
		if !this.Conds.MatchRequest(formatter) {
			return false
		}
		if !this.Conds.MatchResponse(formatter) {
			return false
		}
	}

	// min length
	if this.minLength > 0 && contentLength < this.minLength {
		return false
	}

	// max length
	if this.maxLength > 0 && contentLength > this.maxLength {
		return false
	}

	// except extensions
	if len(this.exceptExtensions) > 0 {
		if len(requestExt) > 0 {
			for _, ext := range this.exceptExtensions {
				if ext == requestExt {
					return false
				}
			}
		}
	}

	// extensions
	if len(this.extensions) > 0 {
		if len(requestExt) > 0 {
			for _, ext := range this.extensions {
				if ext == requestExt {
					return true
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

// MatchAcceptEncoding 鏍规嵁Accept-Encoding閫夋嫨浼樺厛鐨勫帇缂╂柟寮?
func (this *HTTPCompressionConfig) MatchAcceptEncoding(acceptEncodings string) (compressionType HTTPCompressionType, compressionEncoding string, ok bool) {
	if len(acceptEncodings) == 0 {
		return
	}

	if len(this.types) == 0 {
		return
	}

	var pieces = strings.Split(acceptEncodings, ",")
	var encodings = []string{}
	for _, piece := range pieces {
		var qualityIndex = strings.Index(piece, ";")
		if qualityIndex >= 0 {
			// TODO 瀹炵幇浼樺厛绾?
			piece = piece[:qualityIndex]
		}

		encodings = append(encodings, strings.TrimSpace(piece))
	}

	if len(encodings) == 0 {
		return
	}

	for _, supportType := range this.types {
		switch supportType {
		case HTTPCompressionTypeGzip:
			if this.supportGzip && lists.ContainsString(encodings, "gzip") {
				return HTTPCompressionTypeGzip, "gzip", true
			}
		case HTTPCompressionTypeDeflate:
			if this.supportDeflate && lists.ContainsString(encodings, "deflate") {
				return HTTPCompressionTypeDeflate, "deflate", true
			}
		case HTTPCompressionTypeBrotli:
			if this.supportBrotli && lists.ContainsString(encodings, "br") {
				return HTTPCompressionTypeBrotli, "br", true
			}
		case HTTPCompressionTypeZSTD:
			if this.supportZSTD && lists.ContainsString(encodings, "zstd") {
				return HTTPCompressionTypeZSTD, "zstd", true
			}
		}
	}

	return "", "", false
}

func (this *HTTPCompressionConfig) MatchURL(url string) bool {
	// except
	if len(this.ExceptURLPatterns) > 0 {
		for _, pattern := range this.ExceptURLPatterns {
			if pattern.Match(url) {
				return false
			}
		}
	}

	// only
	if len(this.OnlyURLPatterns) > 0 {
		for _, pattern := range this.OnlyURLPatterns {
			if pattern.Match(url) {
				return true
			}
		}
		return false
	}

	return true
}
