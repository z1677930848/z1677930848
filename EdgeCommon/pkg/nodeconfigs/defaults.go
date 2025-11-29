// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package nodeconfigs

import "github.com/iwind/TeaGo/maps"

// 涓€缁勭郴缁熼粯璁ゅ€?
// 淇敼鍗曚釜IP鐩稿叧闄愬埗鍊兼椂瑕佽€冭檻鍒癗AT涓瘡涓狪P浼氫唬琛ㄥ緢澶氫釜涓绘満锛屽苟闈?瀵?鐨勫叧绯?

const (
	DefaultMaxThreads    = 20000   // 鍗曡妭鐐规渶澶х嚎绋嬫暟
	DefaultMaxThreadsMin = 1000    // 鍗曡妭鐐规渶澶х嚎绋嬫暟鏈€灏忓€?
	DefaultMaxThreadsMax = 100_000 // 鍗曡妭鐐规渶澶х嚎绋嬫暟鏈€澶у€?

	DefaultTCPMaxConnections      = 100_000 // 鍗曡妭鐐筎CP鏈€澶ц繛鎺ユ暟
	DefaultTCPMaxConnectionsPerIP = 1000    // 鍗旾P鏈€澶ц繛鎺ユ暟
	DefaultTCPMinConnectionsPerIP = 5       // 鍗旾P鏈€灏忚繛鎺ユ暟

	DefaultTCPNewConnectionsMinutelyRate    = 500 // 鍗旾P杩炴帴閫熺巼闄愬埗锛堟寜鍒嗛挓锛?
	DefaultTCPNewConnectionsMinMinutelyRate = 3   // 鍗旾P鏈€灏忚繛鎺ラ€熺巼

	DefaultTCPNewConnectionsSecondlyRate    = 300 // 鍗旾P杩炴帴閫熺巼闄愬埗锛堟寜绉掞級
	DefaultTCPNewConnectionsMinSecondlyRate = 3   // 鍗旾P鏈€灏忚繛鎺ラ€熺巼

	DefaultTCPLinger           = 5 // 鍗曡妭鐐筎CP Linger鍊?
	DefaultTLSHandshakeTimeout = 3 // TLS鎻℃墜瓒呮椂鏃堕棿
)

var DefaultConfigs = maps.Map{
	"tcpMaxConnections":                DefaultTCPMaxConnections,
	"tcpMaxConnectionsPerIP":           DefaultTCPMaxConnectionsPerIP,
	"tcpMinConnectionsPerIP":           DefaultTCPMinConnectionsPerIP,
	"tcpNewConnectionsMinutelyRate":    DefaultTCPNewConnectionsMinutelyRate,
	"tcpNewConnectionsMinMinutelyRate": DefaultTCPNewConnectionsMinMinutelyRate,
	"tcpNewConnectionsSecondlyRate":    DefaultTCPNewConnectionsSecondlyRate,
	"tcpNewConnectionsMinSecondlyRate": DefaultTCPNewConnectionsMinSecondlyRate,
}
