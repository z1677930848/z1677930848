// Copyright 2021 Lingcdn CDN Lingcdn.cdn@gmail.com. All rights reserved.

package firewallconfigs

import "github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"

type FirewallMode = string

const (
	FirewallModeDefend  FirewallMode = "defend"
	FirewallModeObserve FirewallMode = "observe"
	FirewallModeBypass  FirewallMode = "bypass"
)

func FindAllFirewallModes() []*shared.Definition {
	return []*shared.Definition{
		{Name: "Defend", Description: "apply firewall rules and actions", Code: FirewallModeDefend},
		{Name: "Observe", Description: "apply rules but only log without actions", Code: FirewallModeObserve},
		{Name: "Bypass", Description: "skip firewall rules", Code: FirewallModeBypass},
	}
}

func FindFirewallMode(code FirewallMode) *shared.Definition {
	for _, def := range FindAllFirewallModes() {
		if def.Code == code {
			return def
		}
	}
	return nil
}
