// Copyright 2022 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package nodeconfigs

type DNSResolverType = string

const (
	DNSResolverTypeDefault  = "default"
	DNSResolverTypeGoNative = "goNative"
	DNSResolverTypeCGO      = "cgo"
)

func DefaultDNSResolverConfig() *DNSResolverConfig {
	return &DNSResolverConfig{
		Type: DNSResolverTypeDefault,
	}
}

type DNSResolverConfig struct {
	Type string `yaml:"type" json:"type"` // 浣跨敤Go璇█鍐呯疆鐨凞NS瑙ｆ瀽鍣?
}

func (this *DNSResolverConfig) Init() error {
	return nil
}
