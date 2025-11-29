// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"net/http"
)

// HTTPAuthPolicy HTTP璁よ瘉绛栫暐
type HTTPAuthPolicy struct {
	Id     int64                  `json:"id"`
	Name   string                 `json:"name"`
	IsOn   bool                   `json:"isOn"`
	Type   HTTPAuthType           `json:"type"`
	Params map[string]interface{} `json:"params"`

	method HTTPAuthMethodInterface
}

// MatchRequest 妫€鏌ユ槸鍚﹀尮閰嶈姹?
func (this *HTTPAuthPolicy) MatchRequest(req *http.Request) bool {
	if this.method == nil {
		return false
	}
	return this.method.MatchRequest(req)
}

// Filter 杩囨护
func (this *HTTPAuthPolicy) Filter(req *http.Request, subReqFunc func(subReq *http.Request) (status int, err error), formatter func(string) string) (ok bool, newURI string, uriChanged bool, err error) {
	if this.method == nil {
		// 濡傛灉璁剧疆姝ｇ‘鐨勬柟娉曪紝鎴戜滑鐩存帴鍏佽璇锋眰
		return true, "", false, nil
	}
	return this.method.Filter(req, subReqFunc, formatter)
}

// Method 鑾峰彇璁よ瘉瀹炰緥
func (this *HTTPAuthPolicy) Method() HTTPAuthMethodInterface {
	return this.method
}
