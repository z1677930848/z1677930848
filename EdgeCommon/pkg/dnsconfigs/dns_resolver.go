// Copyright 2023 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved. Official site: https://lingcdn.cloud .

package dnsconfigs

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/iwind/TeaGo/types"
)

type DNSResolver struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func (this *DNSResolver) Addr() string {
	var port = this.Port
	if port <= 0 {
		// 鏆傛椂涓嶆敮鎸丏oH
		// 瀹為檯搴旂敤涓彧鏀寔udp
		switch this.Protocol {
		case "tls":
			port = 853
		default:
			port = 53
		}
	}
	return configutils.QuoteIP(this.Host) + ":" + types.String(port)
}
