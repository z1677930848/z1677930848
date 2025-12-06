// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"regexp"
	"strings"
)

type HTTPRemoteAddrType = string

const (
	HTTPRemoteAddrTypeDefault       HTTPRemoteAddrType = "default"       // 榛樿锛堢洿杩烇級
	HTTPRemoteAddrTypeProxy         HTTPRemoteAddrType = "proxy"         // 浠ｇ悊
	HTTPRemoteAddrTypeRequestHeader HTTPRemoteAddrType = "requestHeader" // 璇锋眰鎶ュご
	HTTPRemoteAddrTypeVariable      HTTPRemoteAddrType = "variable"      // 鍙橀噺
)

// HTTPRemoteAddrConfig HTTP鑾峰彇瀹㈡埛绔疘P鍦板潃鏂瑰紡
type HTTPRemoteAddrConfig struct {
	IsPrior bool               `yaml:"isPrior" json:"isPrior"`
	IsOn    bool               `yaml:"isOn" json:"isOn"`
	Value   string             `yaml:"value" json:"value"` // 鍊煎彉閲?
	Type    HTTPRemoteAddrType `yaml:"type" json:"type"`   // 绫诲瀷

	RequestHeaderName string `yaml:"requestHeaderName" json:"requestHeaderName"` // 璇锋眰鎶ュご鍚嶇О锛坱ype = requestHeader鏃剁敓鏁堬級

	isEmpty   bool
	values    []string
	hasValues bool
}

// Init 鍒濆鍖?
func (this *HTTPRemoteAddrConfig) Init() error {
	this.Value = strings.TrimSpace(this.Value)
	this.isEmpty = false
	if len(this.Value) == 0 {
		this.isEmpty = true
	} else if regexp.MustCompile(`\s+`).ReplaceAllString(this.Value, "") == "${remoteAddr}" {
		this.isEmpty = true
	}

	// values
	this.values = []string{}
	var headerVarReg = regexp.MustCompile(`(\$\{header\.)([\w-,]+)(})`)
	if headerVarReg.MatchString(this.Value) {
		var subMatches = headerVarReg.FindStringSubmatch(this.Value)
		if len(subMatches) > 3 {
			var prefix = subMatches[1]
			var headerNamesString = subMatches[2]
			var suffix = subMatches[3]
			for _, headerName := range strings.Split(headerNamesString, ",") {
				headerName = strings.TrimSpace(headerName)
				if len(headerName) > 0 {
					this.values = append(this.values, prefix+headerName+suffix)
				}
			}
		}
	}
	this.hasValues = len(this.values) > 1 // MUST be 1, not 0

	return nil
}

// IsEmpty 鏄惁涓虹┖
func (this *HTTPRemoteAddrConfig) IsEmpty() bool {
	return this.isEmpty
}

// Values 鍙兘鐨勫€煎彉閲?
func (this *HTTPRemoteAddrConfig) Values() []string {
	return this.values
}

// HasValues 妫€鏌ユ槸鍚︽湁涓€缁勫€?
func (this *HTTPRemoteAddrConfig) HasValues() bool {
	return this.hasValues
}
