// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// AccessLogFileStorageConfig 鏂囦欢瀛樺偍閰嶇疆
type AccessLogFileStorageConfig struct {
	Path       string `yaml:"path" json:"path"`             // 鏂囦欢璺緞锛屾敮鎸佸彉閲忥細${year|month|week|day|hour|minute|second}
	AutoCreate bool   `yaml:"autoCreate" json:"autoCreate"` // 鏄惁鑷姩鍒涘缓鐩綍
}
