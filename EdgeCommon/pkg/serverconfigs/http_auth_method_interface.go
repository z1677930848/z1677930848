// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import "net/http"

// HTTPAuthMethodInterface HTTP璁よ瘉鎺ュ彛瀹氫箟
type HTTPAuthMethodInterface interface {
	// Init 鍒濆鍖?
	Init(params map[string]any) error

	// MatchRequest 鏄惁鍖归厤璇锋眰
	MatchRequest(req *http.Request) bool

	// Filter 杩囨护
	Filter(req *http.Request, subReqFunc func(subReq *http.Request) (status int, err error), formatter func(string) string) (ok bool, newURI string, uriChanged bool, err error)

	// SetExts 璁剧疆鎵╁睍鍚?
	SetExts(exts []string)

	// SetDomains 璁剧疆鍩熷悕
	SetDomains(domains []string)
}
