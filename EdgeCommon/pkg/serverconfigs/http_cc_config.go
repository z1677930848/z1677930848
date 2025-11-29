// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package serverconfigs

type HTTPCCThreshold struct {
	// stub
}

func NewHTTPCCThreshold() *HTTPCCThreshold {
	return &HTTPCCThreshold{}
}

func (this *HTTPCCThreshold) Merge(threshold *HTTPCCThreshold) {
	// stub
}

func (this *HTTPCCThreshold) MergeIfEmpty(threshold *HTTPCCThreshold) {
	// stub
}

func (this *HTTPCCThreshold) Clone() *HTTPCCThreshold {
	return &HTTPCCThreshold{}
}

var DefaultHTTPCCThresholds = []*HTTPCCThreshold{} // stub

// DefaultHTTPCCConfig 榛樿鐨凜C閰嶇疆
func DefaultHTTPCCConfig() *HTTPCCConfig {
	return &HTTPCCConfig{}
}

// HTTPCCConfig HTTP CC闃叉姢閰嶇疆
type HTTPCCConfig struct {
	IsPrior    bool               `yaml:"isPrior" json:"isPrior"`       // 鏄惁瑕嗙洊鐖剁骇
	IsOn       bool               `yaml:"isOn" json:"isOn"`             // 鏄惁鍚敤
	Thresholds []*HTTPCCThreshold `yaml:"thresholds" json:"thresholds"` // 闃堝€艰缃?
}

func NewHTTPCCConfig() *HTTPCCConfig {
	return &HTTPCCConfig{}
}

func (this *HTTPCCConfig) Init() error {
	return nil
}

func (this *HTTPCCConfig) MatchURL(url string) bool {
	return false
}
