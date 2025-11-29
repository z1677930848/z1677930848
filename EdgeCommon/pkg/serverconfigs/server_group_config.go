// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package serverconfigs

// ServerGroupConfig 鏈嶅姟鍒嗙粍閰嶇疆
type ServerGroupConfig struct {
	Id   int64  `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
	IsOn bool   `yaml:"isOn" json:"isOn"`

	// 鍙嶅悜浠ｇ悊閰嶇疆
	HTTPReverseProxyRef *ReverseProxyRef    `yaml:"httpReverseProxyRef" json:"httpReverseProxyRef"`
	HTTPReverseProxy    *ReverseProxyConfig `yaml:"httpReverseProxy" json:"httpReverseProxy"`
	TCPReverseProxyRef  *ReverseProxyRef    `yaml:"tcpReverseProxyRef" json:"tcpReverseProxyRef"`
	TCPReverseProxy     *ReverseProxyConfig `yaml:"tcpReverseProxy" json:"tcpReverseProxy"`
	UDPReverseProxyRef  *ReverseProxyRef    `yaml:"udpReverseProxyRef" json:"udpReverseProxyRef"`
	UDPReverseProxy     *ReverseProxyConfig `yaml:"udpReverseProxy" json:"udpReverseProxy"`

	Web *HTTPWebConfig `yaml:"web" json:"web"`
}
