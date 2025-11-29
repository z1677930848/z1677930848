//go:build !cgo

package injectionutils

import "github.com/TeaOSLab/EdgeNode/internal/waf/utils"

// 在禁用 cgo 的环境下不执行 libinjection 检测，直接返回未命中。
func DetectSQLInjectionCache(input string, isStrict bool, cacheLife utils.CacheLife) bool {
	return false
}

func DetectSQLInjection(input string, isStrict bool) bool {
	return false
}

func DetectXSSCache(input string, isStrict bool, cacheLife utils.CacheLife) bool {
	return false
}

func DetectXSS(input string, isStrict bool) bool {
	return false
}
