//go:build !cgo
// +build !cgo

// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .

package injectionutils

import (
	"github.com/TeaOSLab/EdgeNode/internal/waf/utils"
)

// DetectSQLInjectionCache detect sql injection in string with cache (stub for non-CGO builds)
func DetectSQLInjectionCache(input string, isStrict bool, cacheLife utils.CacheLife) bool {
	return false // WAF功能在非CGO编译中禁用
}

// DetectSQLInjection detect sql injection in string (stub for non-CGO builds)
func DetectSQLInjection(input string, isStrict bool) bool {
	return false // WAF功能在非CGO编译中禁用
}

// DetectXSSCache detect xss in string with cache (stub for non-CGO builds)
func DetectXSSCache(input string, isStrict bool, cacheLife utils.CacheLife) bool {
	return false // WAF功能在非CGO编译中禁用
}

// DetectXSS detect xss in string (stub for non-CGO builds)
func DetectXSS(input string, isStrict bool) bool {
	return false // WAF功能在非CGO编译中禁用
}
