// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package reporterconfigs

// GlobalSetting 鍏ㄥ眬璁剧疆
type GlobalSetting struct {
	MinNotifyConnectivity float64 `json:"minNotifyConnectivity"` // 闇€瑕侀€氱煡鐨勬渶灏忚繛閫氬€?
	NotifyWebHookURL      string  `json:"notifyWebHookURL"`      // WebHook閫氱煡鍦板潃
}

func DefaultGlobalSetting() *GlobalSetting {
	return &GlobalSetting{
		MinNotifyConnectivity: 100,
	}
}
