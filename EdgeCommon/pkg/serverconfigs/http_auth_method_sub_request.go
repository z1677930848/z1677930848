// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"regexp"
	"time"
)

var httpAuthSubRequestHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

// HTTPAuthSubRequestMethod 浣跨敤URL璁よ瘉
type HTTPAuthSubRequestMethod struct {
	HTTPAuthBaseMethod

	URL    string `json:"url"`
	Method string `json:"method"`

	// TODO 澧炲姞鑷畾涔塇eader銆佽秴鏃躲€佽瘉涔︾瓑

	isFullURL bool
}

func NewHTTPAuthSubRequestMethod() *HTTPAuthSubRequestMethod {
	return &HTTPAuthSubRequestMethod{}
}

// Init 鍒濆鍖?
func (this *HTTPAuthSubRequestMethod) Init(params map[string]any) error {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}
	err = json.Unmarshal(paramsJSON, this)
	if err != nil {
		return err
	}

	// 鏄惁鏄畬鏁寸殑URL
	this.isFullURL = false
	if regexp.MustCompile(`^(?i)(http|https)://`).MatchString(this.URL) {
		this.isFullURL = true
	} else {
		if len(this.URL) == 0 || this.URL[0] != '/' {
			this.URL = "/" + this.URL
		}
	}

	return nil
}

// Filter 杩囨护
func (this *HTTPAuthSubRequestMethod) Filter(req *http.Request, doSubReq func(subReq *http.Request) (status int, err error), formatter func(string) string) (ok bool, newURI string, uriChanged bool, err error) {
	var method = this.Method
	if len(method) == 0 {
		method = req.Method
	}

	var url = formatter(this.URL)
	var scheme = "http"
	if req.TLS != nil {
		scheme = "https"
	}
	var host = req.URL.Host
	if len(host) == 0 {
		host = req.Host
	}
	if !this.isFullURL {
		url = scheme + "://" + host + url
	}
	newReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false, "", false, err
	}
	for k, v := range req.Header {
		if k != "Connection" {
			newReq.Header[k] = v
		}
	}

	if !this.isFullURL {
		status, err := doSubReq(newReq)
		if err != nil {
			return false, "", false, err
		}
		return status >= 200 && status < 300, "", false, nil
	}

	// TODO 闇€瑕佸皢Header鍜孲tatusCode銆丷esponseBody杈撳嚭鍒板鎴风
	newReq.Header.Set("Referer", scheme+"://"+host+req.URL.RequestURI())
	resp, err := httpAuthSubRequestHTTPClient.Do(newReq)
	if err != nil {
		return false, "", false, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return resp.StatusCode >= 200 && resp.StatusCode < 300, "", false, nil
}
