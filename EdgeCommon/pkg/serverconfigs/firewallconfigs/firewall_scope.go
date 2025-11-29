// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package firewallconfigs

type FirewallScope = string

const (
	FirewallScopeGlobal FirewallScope = "global"
	FirewallScopeServer FirewallScope = "service" // 鍘嗗彶鍘熷洜锛屼唬鍙蜂负 service 鑰岄潪 server
)
