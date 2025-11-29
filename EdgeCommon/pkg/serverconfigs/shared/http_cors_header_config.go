// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package shared

// HTTPCORSHeaderConfig 鍙傝€?https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
type HTTPCORSHeaderConfig struct {
	IsOn              bool     `yaml:"isOn" json:"isOn"`
	AllowMethods      []string `yaml:"allowMethods" json:"allowMethods"`
	AllowOrigin       string   `yaml:"allowOrigin" json:"allowOrigin"`           // TODO
	AllowCredentials  bool     `yaml:"allowCredentials" json:"allowCredentials"` // TODO锛屽疄鐜版椂闇€瑕佸崌绾т互寰€鐨勮€佹暟鎹?
	ExposeHeaders     []string `yaml:"exposeHeaders" json:"exposeHeaders"`
	MaxAge            int32    `yaml:"maxAge" json:"maxAge"`
	RequestHeaders    []string `yaml:"requestHeaders" json:"requestHeaders"` // TODO
	RequestMethod     string   `yaml:"requestMethod" json:"requestMethod"`
	OptionsMethodOnly bool     `yaml:"optionsMethodOnly" json:"optionsMethodOnly"` // 鏄惁浠呮敮鎸丱PTIONS鏂规硶
}

func NewHTTPCORSHeaderConfig() *HTTPCORSHeaderConfig {
	return &HTTPCORSHeaderConfig{
		AllowCredentials: true,
	}
}

func (this *HTTPCORSHeaderConfig) Init() error {
	return nil
}
